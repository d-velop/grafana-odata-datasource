package plugin

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/stretchr/testify/mock"
)

type clientMock struct {
	baseUrl    string
	statusCode int
	body       []byte
	err        error
	mock.Mock
}

type managerMock struct {
	mock.Mock
}

type callResourceResponseSenderMock struct {
	csr *backend.CallResourceResponse
}

func (client *clientMock) ODataVersion() string {
	return "V4"
}

func (client *clientMock) GetServiceRoot() (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(bytes.NewReader(client.body))}, client.err
}

func (client *clientMock) GetMetadata() (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(bytes.NewReader(client.body))}, client.err
}

func (client *clientMock) Get(_ string, _ []property, _ string, _ backend.TimeRange, _ []filterCondition) (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(bytes.NewReader(client.body))}, client.err
}

func (im *managerMock) Get(ctx context.Context, pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	args := im.Called(ctx, pluginContext)
	return args.Get(0), args.Error(1)
}

func (im *managerMock) Do(ctx context.Context, pluginContext backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	args := im.Called(ctx, pluginContext, fn)
	return args.Error(0)
}

func (crrsm *callResourceResponseSenderMock) Send(csr *backend.CallResourceResponse) error {
	crrsm.csr = csr
	return nil
}
