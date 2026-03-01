package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ODataClient interface {
	ODataVersion() string
	GetServiceRoot(ctx context.Context) (*http.Response, error)
	GetMetadata(ctx context.Context) (*http.Response, error)
	Get(ctx context.Context, entitySet string, properties []property,
		filterConditions []filterCondition) (*http.Response, error)
}

type ODataClientImpl struct {
	httpClient       *http.Client
	baseUrl          string
	urlSpaceEncoding string
	odataVersion     string
}

func (client *ODataClientImpl) ODataVersion() string {
	return client.odataVersion
}

func (client *ODataClientImpl) get(ctx context.Context, url string, mimeType string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new request with context: %w", err)
	}
	req.Header.Set("Accept", mimeType)
	return client.httpClient.Do(req)
}

func (client *ODataClientImpl) GetServiceRoot(ctx context.Context) (*http.Response, error) {
	return client.get(ctx, client.baseUrl, "application/json")
}

func (client *ODataClientImpl) GetMetadata(ctx context.Context) (*http.Response, error) {
	requestUrl, err := url.Parse(client.baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, odata.Metadata)
	return client.get(ctx, requestUrl.String(), "application/xml")
}

func (client *ODataClientImpl) Get(ctx context.Context, entitySet string, properties []property, filterConditions []filterCondition) (*http.Response, error) {
	requestUrl, err := buildQueryUrl(client.baseUrl, entitySet, properties,
		filterConditions, client.urlSpaceEncoding, client.odataVersion)
	if err != nil {
		return nil, err
	}
	urlString := requestUrl.String()
	log.DefaultLogger.Debug("Constructed request url", "url", urlString)
	return client.get(ctx, urlString, "application/json")
}

func buildQueryUrl(baseUrl string, entitySet string, properties []property, filterConditions []filterCondition, urlSpaceEncoding string, version string) (*url.URL, error) {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, entitySet)
	params, _ := url.ParseQuery(requestUrl.RawQuery)
	filterParam := mapFilter(filterConditions, version)
	if len(filterParam) > 0 {
		params.Add(odata.Filter, filterParam)
	}
	selectParam := mapSelect(properties)
	if len(selectParam) > 0 {
		params.Add(odata.Select, selectParam)
	}
	encodedUrl := params.Encode()
	if urlSpaceEncoding == "%20" {
		encodedUrl = strings.ReplaceAll(encodedUrl, "+", "%20")
	}
	requestUrl.RawQuery = encodedUrl
	return requestUrl, nil
}

func mapSelect(properties []property) string {
	var result []string
	if len(properties) > 0 {
		for _, selectProp := range properties {
			result = append(result, selectProp.Name)
		}
	}
	return strings.Join(result[:], ",")
}

func mapFilter(filterConditions []filterCondition, version string) string {
	isV2V3 := version == "V2" || version == "V3"
	var filter = ""
	for index, element := range filterConditions {
		name, op, val := element.Property.Name, element.Operator, element.Value
		switch element.Property.Type {
		case odata.EdmString:
			filter += fmt.Sprintf("%s %s '%s'", name, op, val)
		case odata.EdmDateTime:
			if isV2V3 {
				// V2/V3: datetime'yyyy-mm-ddThh:mm:ss' without timezone suffix
				filter += fmt.Sprintf("%s %s datetime'%s'", name, op, stripTimezoneForV2(val))
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmDateTimeOffset:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s datetimeoffset'%s'", name, op, val)
			} else {
				// V4: plain ISO 8601 value, no prefix
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmInt64:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s %sL", name, op, val)
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmDecimal:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s %sM", name, op, val)
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmSingle:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s %sf", name, op, val)
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmDouble:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s %sd", name, op, val)
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		case odata.EdmGuid:
			if isV2V3 {
				filter += fmt.Sprintf("%s %s guid'%s'", name, op, val)
			} else {
				filter += fmt.Sprintf("%s %s %s", name, op, val)
			}
		default:
			filter += fmt.Sprintf("%s %s %s", name, op, val)
		}
		if index < (len(filterConditions) - 1) {
			filter += " and "
		}
	}
	return filter
}

func stripTimezoneForV2(value string) string {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.UTC().Format("2006-01-02T15:04:05")
	}
	return value
}
