package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ODataClient interface {
	ODataVersion() string
	GetServiceRoot() (*http.Response, error)
	GetMetadata() (*http.Response, error)
	Get(entitySet string, properties []property, timeProperty property, timeRange backend.TimeRange,
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

func (client *ODataClientImpl) GetServiceRoot() (*http.Response, error) {
	return get(client.httpClient, client.baseUrl, "application/json")
}

func (client *ODataClientImpl) GetMetadata() (*http.Response, error) {
	requestUrl, err := url.Parse(client.baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, odata.Metadata)
	return get(client.httpClient, requestUrl.String(), "application/xml")
}

func (client *ODataClientImpl) Get(entitySet string, properties []property, timeProperty property,
	timeRange backend.TimeRange, filterConditions []filterCondition) (*http.Response, error) {
	requestUrl, err := buildQueryUrl(client.baseUrl, entitySet, properties, timeProperty, timeRange,
		filterConditions, client.urlSpaceEncoding)
	if err != nil {
		return nil, err
	}
	return get(client.httpClient, requestUrl.String(), "application/json")
}

func get(httpClient *http.Client, url string, accept string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	return httpClient.Do(req)
}

func buildQueryUrl(baseUrl string, entitySet string, properties []property, timeProperty property,
	timeRange backend.TimeRange, filterConditions []filterCondition, urlSpaceEncoding string) (*url.URL, error) {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, entitySet)
	params, _ := url.ParseQuery(requestUrl.RawQuery)
	filterParam := mapFilter(timeProperty, timeRange, filterConditions)
	if len(filterParam) > 0 {
		params.Add(odata.Filter, filterParam)
	}
	selectParam := mapSelect(properties, timeProperty)
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

func mapSelect(properties []property, timeProperty property) string {
	var result []string
	if len(properties) > 0 {
		for _, selectProp := range properties {
			result = append(result, selectProp.Name)
		}
	}
	if len(timeProperty.Name) > 0 {
		result = append(result, timeProperty.Name)
	}
	return strings.Join(result[:], ",")
}

func mapFilter(timeProperty property, timeRange backend.TimeRange, filterConditions []filterCondition) string {
	var filter string
	if len(timeProperty.Name) > 0 {
		filter = fmt.Sprintf("%s ge %s and %s le %s", timeProperty.Name, timeRange.From.UTC().Format(time.RFC3339),
			timeProperty.Name, timeRange.To.UTC().Format(time.RFC3339))
	}
	var customFilter = ""
	for index, element := range filterConditions {
		if element.Property.Type == odata.EdmString {
			customFilter += fmt.Sprintf("%s %s '%s'", element.Property.Name, element.Operator, element.Value)
		} else {
			customFilter += fmt.Sprintf("%s %s %s", element.Property.Name, element.Operator, element.Value)
		}
		if index < (len(filterConditions) - 1) {
			customFilter += " and "
		}
	}
	if len(customFilter) > 0 {
		filter += " and " + customFilter
	}
	return filter
}
