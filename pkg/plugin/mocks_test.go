package plugin

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/stretchr/testify/mock"
)

type clientMock struct {
	baseUrl             string
	statusCode          int
	body                string
	metadataResponse    *odata.Edmx
	getEntitiesResponse *odata.Response
	err                 error
	mock.Mock
}

type managerMock struct {
	mock.Mock
}

type callResourceResponseSenderMock struct {
	csr *backend.CallResourceResponse
}

func (client *clientMock) GetServiceRoot() (*http.Response, error) {
	return &http.Response{StatusCode: client.statusCode,
		Body: io.NopCloser(strings.NewReader(client.body))}, client.err
}

func (client *clientMock) GetMetadata() (*odata.Edmx, error) {
	return client.metadataResponse, client.err
}

func (client *clientMock) GetEntities(_ string, _ []property, _ string, _ backend.TimeRange, _ []filterCondition) (*odata.Response, error) {
	return client.getEntitiesResponse, client.err
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
