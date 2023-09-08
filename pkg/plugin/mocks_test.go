package plugin

import (
	"io"
	"net/http"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/stretchr/testify/mock"
)

type clientMock struct {
	statusCode int
	body       string
	baseUrl    string
	mock.Mock
}

type managerMock struct {
	mock.Mock
}

func (client *clientMock) GetServiceRoot() (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(strings.NewReader(client.body))}, nil
}

func (client *clientMock) GetMetadata() (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(strings.NewReader(client.body))}, nil
}

func (client *clientMock) Get(_ string, _ []property, _ string, _ backend.TimeRange, _ []filterCondition) (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(strings.NewReader(client.body))}, nil
}

func (im *managerMock) Get(pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	args := im.Called(pluginContext)
	return args.Get(0), args.Error(1)
}

func (im *managerMock) Do(pluginContext backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	args := im.Called(pluginContext, fn)
	return args.Error(0)
}
