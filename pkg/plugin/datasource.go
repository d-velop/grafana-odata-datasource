package plugin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

func newDatasourceInstance(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	clientOptions, err := settings.HTTPClientOptions()
	if err != nil {
		return nil, err
	}
	client, err := httpclient.New(clientOptions)
	if err != nil {
		return nil, err
	}
	return &ODataSourceInstance{
		&ODataClientImpl{client, settings.URL},
	}, nil
}

type ODataSourceInstance struct {
	client ODataClient
}

func NewODataSource(_ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	im := datasource.NewInstanceManager(newDatasourceInstance)
	ds := &ODataSource{
		im: im,
	}
	return ds, nil
}

func (ds *ODataSource) getClientInstance(pluginContext backend.PluginContext) ODataClient {
	instance, _ := ds.im.Get(pluginContext)
	clientInstance := instance.(*ODataSourceInstance).client
	return clientInstance
}

func (ds *ODataSource) QueryData(_ context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse,
	error) {
	clientInstance := ds.getClientInstance(req.PluginContext)
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		res := ds.query(clientInstance, q)
		response.Responses[q.RefID] = res
	}
	return response, nil
}

func (ds *ODataSource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult,
	error) {
	var status backend.HealthStatus
	var message string
	clientInstance := ds.getClientInstance(req.PluginContext)
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

func (ds *ODataSource) CallResource(_ context.Context, req *backend.CallResourceRequest,
	sender backend.CallResourceResponseSender) error {
	switch req.Path {
	case "metadata":
		return ds.getMetadata(req, sender)
	default:
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
		})
	}
}

func (ds *ODataSource) query(clientInstance ODataClient, query backend.DataQuery) backend.DataResponse {
	log.DefaultLogger.Debug("query", "query.JSON", string(query.JSON))
	response := backend.DataResponse{}
	var qm queryModel
	err := json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		response.Error = fmt.Errorf("error unmarshalling query json: %w", err)
		return response
	}
	timeProperty := qm.TimeProperty.Name
	frame := data.NewFrame("response")
	frame.Name = query.RefID
	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}
	frame.Meta.PreferredVisualization = data.VisTypeTable
	labels, err := data.LabelsFromString("time=" + timeProperty)
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
	var entities []map[string]interface{}
	entities, err = ds.getEntities(clientInstance, qm.EntitySet.Name, qm.Properties, timeProperty,
		query.TimeRange, qm.FilterConditions)
	log.DefaultLogger.Debug("query complete", "noOfEntities", len(entities))
	if err != nil {
		response.Error = err
		return response
	}
	for _, entry := range entities {
		values := make([]interface{}, len(qm.Properties)+1)
		if timeValue, err := time.Parse(time.RFC3339Nano, fmt.Sprint(entry[timeProperty])); err == nil {
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

func (ds *ODataSource) getEntities(client ODataClient, entitySet string, properties []property, timeProperty string,
	timeRange backend.TimeRange, filterConditions []filterCondition) ([]map[string]interface{}, error) {
	resp, err := client.Get(entitySet, properties, timeProperty, timeRange, filterConditions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.DefaultLogger.Debug("request response status", "status", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get failed with status code %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result odata.Response
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return result.Value, nil
}

func (ds *ODataSource) getMetadata(req *backend.CallResourceRequest,
	sender backend.CallResourceResponseSender) error {
	clientInstance := ds.getClientInstance(req.PluginContext)
	resp, err := clientInstance.GetMetadata()
	if err != nil {
		log.DefaultLogger.Error("error in http request")
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
	metadata := make(map[string]interface{})
	entityTypes := make(map[string]interface{})
	entitySets := make(map[string]interface{})
	for _, ds := range edmx.DataServices {
		for _, s := range ds.Schemas {
			for _, et := range s.EntityTypes {
				qualifiedName := s.Namespace + "." + et.Name
				var properties []interface{}
				for _, p := range et.Properties {
					prop := map[string]interface{}{
						"name": p.Name,
						"type": p.Type,
					}
					properties = append(properties, prop)
				}
				entityTypes[qualifiedName] = map[string]interface{}{
					"name":          et.Name,
					"qualifiedName": qualifiedName,
					"properties":    properties,
				}
			}
			for _, ec := range s.EntityContainers {
				for _, es := range ec.EntitySet {
					entitySets[es.Name] = map[string]interface{}{
						"name":       es.Name,
						"entityType": es.EntityType,
					}
				}
			}
		}
	}
	metadata["entityTypes"] = entityTypes
	metadata["entitySets"] = entitySets
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
