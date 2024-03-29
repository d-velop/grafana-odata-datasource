package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type ODataClient interface {
	GetServiceRoot() (*http.Response, error)
	GetMetadata() (*http.Response, error)
	Get(entitySet string, properties []property,
		filterConditions []filterCondition) (*http.Response, error)
}

type ODataClientImpl struct {
	httpClient       *http.Client
	baseUrl          string
	urlSpaceEncoding string
}

func (client *ODataClientImpl) GetServiceRoot() (*http.Response, error) {
	return client.httpClient.Get(client.baseUrl)
}

func (client *ODataClientImpl) GetMetadata() (*http.Response, error) {
	requestUrl, err := url.Parse(client.baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, odata.Metadata)
	return client.httpClient.Get(requestUrl.String())
}

func (client *ODataClientImpl) Get(entitySet string, properties []property, filterConditions []filterCondition) (*http.Response, error) {
	requestUrl, err := buildQueryUrl(client.baseUrl, entitySet, properties,
		filterConditions, client.urlSpaceEncoding)
	if err != nil {
		return nil, err
	}
	urlString := requestUrl.String()
	log.DefaultLogger.Debug("Constructed request url: ", urlString)
	return client.httpClient.Get(urlString)
}

func buildQueryUrl(baseUrl string, entitySet string, properties []property, filterConditions []filterCondition, urlSpaceEncoding string) (*url.URL, error) {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.Path = path.Join(requestUrl.Path, entitySet)
	params, _ := url.ParseQuery(requestUrl.RawQuery)
	filterParam := mapFilter(filterConditions)
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

func mapFilter(filterConditions []filterCondition) string {
	var filter = ""
	for index, element := range filterConditions {
		if element.Property.Type == odata.EdmString {
			filter += fmt.Sprintf("%s %s '%s'", element.Property.Name, element.Operator, element.Value)
		} else {
			filter += fmt.Sprintf("%s %s %s", element.Property.Name, element.Operator, element.Value)
		}
		if index < (len(filterConditions) - 1) {
			filter += " and "
		}
	}

	return filter
}
