package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

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
		filterConditions, client.urlSpaceEncoding)
	if err != nil {
		return nil, err
	}
	urlString := requestUrl.String()
	log.DefaultLogger.Debug("Constructed request url", "url", urlString)
	return client.get(ctx, urlString, "application/json")
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
