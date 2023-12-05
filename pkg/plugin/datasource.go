package plugin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler    = (*ODataSource)(nil)
	_ backend.CheckHealthHandler  = (*ODataSource)(nil)
	_ backend.CallResourceHandler = (*ODataSource)(nil)
)

type ODataSource struct {
	im instancemgmt.InstanceManager
}

type DatasourceSettings struct {
	URLSpaceEncoding string `json:"urlSpaceEncoding"`
	ODataVersion     string `json:"odataVersion"`
}

func newDatasourceInstance(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	clientOptions, err := settings.HTTPClientOptions(ctx)
	if err != nil {
		return nil, err
	}
	client, err := httpclient.New(clientOptions)
	if err != nil {
		return nil, err
	}

	var dsSettings DatasourceSettings
	if settings.JSONData != nil && len(settings.JSONData) > 1 {
		if err := json.Unmarshal(settings.JSONData, &dsSettings); err != nil {
			return nil, err
		}
	}

	return &ODataSourceInstance{
		&ODataClientImpl{client, strings.TrimSuffix(settings.URL, "/"), dsSettings.URLSpaceEncoding, dsSettings.ODataVersion},
	}, nil
}

type ODataSourceInstance struct {
	client ODataClient
}

func NewODataSource(ctx context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	im := datasource.NewInstanceManager(newDatasourceInstance)
	ds := &ODataSource{
		im: im,
	}
	return ds, nil
}

func (ds *ODataSource) getClientInstance(ctx context.Context, pluginContext backend.PluginContext) ODataClient {
	instance, _ := ds.im.Get(ctx, pluginContext)
	clientInstance := instance.(*ODataSourceInstance).client
	return clientInstance
}

func (ds *ODataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse,
	error) {
	clientInstance := ds.getClientInstance(ctx, req.PluginContext)
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		res := ds.query(clientInstance, q)
		response.Responses[q.RefID] = res
	}
	return response, nil
}

func (ds *ODataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult,
	error) {
	var status backend.HealthStatus
	var message string
	clientInstance := ds.getClientInstance(ctx, req.PluginContext)
	var res, err = clientInstance.GetServiceRoot()
	if err != nil {
		status = backend.HealthStatusError
		message = fmt.Sprintf("Health check failed: %s", err.Error())
	} else {
		if res.StatusCode == 200 {
			status = backend.HealthStatusOk
			message = "Data Source is working as expected."
		} else {
			status = backend.HealthStatusError
			message = fmt.Sprintf("Health check failed, datasource exists but given path does not. "+
				"Statuscode: %d", res.StatusCode)
		}
	}
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (ds *ODataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest,
	sender backend.CallResourceResponseSender) error {
	switch req.Path {
	case "metadata":
		return ds.getMetadata(ctx, req, sender)
	default:
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
		})
	}
}

func mapToResponse(bodyBytes []byte) ([]map[string]interface{}, error) {
	var response odata.Response
	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}

func mapToV2Response(bodyBytes []byte) ([]map[string]interface{}, error) {
	var response odata.ResponseV2
	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return response.D.Results, nil
}

func (ds *ODataSource) query(clientInstance ODataClient, query backend.DataQuery) backend.DataResponse {
	log.DefaultLogger.Debug("query", "query.JSON", string(query.JSON))
	response := backend.DataResponse{}
	var qm queryModel
	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		response.Error = fmt.Errorf("error unmarshalling query json: %w", err)
		return response
	}

	timeProperty := qm.TimeProperty
	frame := data.NewFrame("response")
	frame.Name = query.RefID
	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}
	frame.Meta.PreferredVisualization = data.VisTypeTable
	labels, err := data.LabelsFromString("time=" + timeProperty.Name)
	if err != nil {
		response.Error = err
		return response
	}
	frame.Fields = append(frame.Fields,
		data.NewField("time", labels, []*time.Time{}),
	)
	for _, prop := range qm.Properties {
		frame.Fields = append(frame.Fields, data.NewField(prop.Name, nil, odata.ToArray(prop.Type)))
	}

	resp, err := clientInstance.Get(qm.EntitySet.Name, qm.Properties, timeProperty,
		query.TimeRange, qm.FilterConditions)
	if err != nil {
		response.Error = err
		return response
	}

	defer resp.Body.Close()

	log.DefaultLogger.Debug("request response status", "status", resp.Status)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.DefaultLogger.Error("error reading response body", "err", err)
			response.Error = fmt.Errorf("get failed - code %d", resp.StatusCode)
		} else {
			response.Error = fmt.Errorf("get failed - code %d: %s", resp.StatusCode, string(bodyBytes))
		}
		return response
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return response
	}
	version := clientInstance.ODataVersion()
	if version == "Auto" || version == "" {
		odataVersion := resp.Header.Get("DataServiceVersion")
		if strings.HasPrefix(odataVersion, "2") {
			version = "V2"
		} else if strings.HasPrefix(odataVersion, "3") {
			version = "V3"
		} else {
			odataVersion = resp.Header.Get("OData-Version")
			if strings.HasPrefix(odataVersion, "4") {
				version = "V4"
			}
		}
	}
	log.DefaultLogger.Debug("using odata version", "version", version)
	var entries []map[string]interface{}
	if version == "V2" {
		entries, err = mapToV2Response(bodyBytes)
		if err != nil {
			response.Error = err
			return response
		}
	} else {
		entries, err = mapToResponse(bodyBytes)
		if err != nil {
			response.Error = err
			return response
		}
	}
	log.DefaultLogger.Debug("query complete", "noOfEntities", len(entries))
	for _, entry := range entries {
		values := make([]interface{}, len(qm.Properties)+1)
		if timeValue, err := odata.ParseTime(fmt.Sprint(entry[timeProperty.Name])); err == nil {
			values[0] = &timeValue
		} else {
			values[0] = nil
		}
		for i, prop := range qm.Properties {
			if value, ok := entry[prop.Name]; ok {
				values[i+1] = odata.MapValue(value, prop.Type)
			} else {
				values[i+1] = nil
			}
		}
		frame.AppendRow(values...)
	}
	response.Frames = append(response.Frames, frame)
	return response
}

func (ds *ODataSource) getMetadata(ctx context.Context, req *backend.CallResourceRequest,
	sender backend.CallResourceResponseSender) error {
	clientInstance := ds.getClientInstance(ctx, req.PluginContext)
	resp, err := clientInstance.GetMetadata()

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get metadata failed with status code %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Error("error reading response body")
		return err
	}
	var edmx odata.Edmx
	err = xml.Unmarshal(bodyBytes, &edmx)
	if err != nil {
		log.DefaultLogger.Error("error unmarshalling response body")
		return err
	}

	metadata := schema{
		EntityTypes: make(map[string]entityType),
		EntitySets:  make(map[string]entitySet),
	}
	for _, ds := range edmx.DataServices {
		for _, s := range ds.Schemas {
			for _, et := range s.EntityTypes {
				qualifiedName := s.Namespace + "." + et.Name
				var properties []property
				for _, p := range et.Properties {
					prop := property{
						Name: p.Name,
						Type: p.Type,
					}
					properties = append(properties, prop)
				}
				metadata.EntityTypes[qualifiedName] = entityType{
					Name:          et.Name,
					QualifiedName: qualifiedName,
					Properties:    properties,
				}
			}
			for _, ec := range s.EntityContainers {
				for _, es := range ec.EntitySet {
					metadata.EntitySets[es.Name] = entitySet{
						Name:       es.Name,
						EntityType: es.EntityType,
					}
				}
			}
		}
	}

	responseBody, err := json.Marshal(metadata)
	if err != nil {
		log.DefaultLogger.Error("error marshalling response body")
		return err
	}
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   responseBody,
	})
}
