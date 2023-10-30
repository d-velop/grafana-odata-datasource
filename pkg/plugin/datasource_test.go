package plugin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetClientInstance(t *testing.T) {
	// Arrange
	client := &clientMock{}
	im := managerMock{}
	ds := ODataSource{&im}

	is := ODataSourceInstance{client}
	im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

	// Act
	odataClient := ds.getClientInstance(context.TODO(), backend.PluginContext{})

	// Assert
	require.Equal(t, client, odataClient)
}

func TestNewODataSourceInstance(t *testing.T) {
	tables := []backend.DataSourceInstanceSettings{
		{
			URL: "http://localhost:8080",
		},
	}

	for _, table := range tables {
		// Arrange
		// Act
		dsi, err := newDatasourceInstance(context.TODO(), table)

		// Assert
		require.NoError(t, err)

		odsi := dsi.(*ODataSourceInstance)
		odsic := odsi.client.(*ODataClientImpl)

		require.Equal(t, table.URL, odsic.baseUrl)
	}
}

func TestCallResource(t *testing.T) {
	tables := []struct {
		req         *backend.CallResourceRequest
		expRespCode int
	}{
		{
			req: &backend.CallResourceRequest{
				Path: "http://localhost:8080/path/does/not/exist",
			},
			expRespCode: 404,
		},
		{
			req: &backend.CallResourceRequest{
				Path: "/path/does/not/exist",
			},
			expRespCode: 404,
		},
		{
			req: &backend.CallResourceRequest{
				Path: "metadata",
			},
			expRespCode: 200,
		},
	}

	for _, table := range tables {
		// Arrange
		client := &clientMock{
			metadataResponse: &odata.Edmx{},
		}
		im := managerMock{}
		ds := ODataSource{&im}

		is := ODataSourceInstance{client}
		im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)
		crs := callResourceResponseSenderMock{}

		// Act
		err := ds.CallResource(context.TODO(), table.req, &crs)

		// Assert
		require.NoError(t, err)
		require.Equal(t, table.expRespCode, crs.csr.Status)
	}
}

