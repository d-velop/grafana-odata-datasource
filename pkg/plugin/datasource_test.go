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

func TestNewODataSourceInstanceSuccess(t *testing.T) {
	// Arrange
	url := "http://localhost:8080"

	// Act
	dsi, err := newDatasourceInstance(context.TODO(), backend.DataSourceInstanceSettings{
		URL: url,
	})

	// Assert
	require.NoError(t, err)

	odsi := dsi.(*ODataSourceInstance)
	odsic := odsi.client.(*ODataClientImpl)

	require.Equal(t, url, odsic.baseUrl)
}

func TestNewODataSourceInstanceInvalidJSON(t *testing.T) {
	// Act
	dsi, err := newDatasourceInstance(context.TODO(), backend.DataSourceInstanceSettings{
		URL:      "http://localhost:8080",
		JSONData: []byte(`this is no json`),
	})

	// Assert
	require.Error(t, err)
	require.Nil(t, dsi)
}

func TestCallResource(t *testing.T) {
	tables := []struct {
		name        string
		req         *backend.CallResourceRequest
		expRespCode int
	}{
		{
			name: "Call not found absolute",
			req: &backend.CallResourceRequest{
				Path: "http://localhost:8080/path/does/not/exist",
			},
			expRespCode: 404,
		},
		{
			name: "Call not found",
			req: &backend.CallResourceRequest{
				Path: "/path/does/not/exist",
			},
			expRespCode: 404,
		},
		{
			name: "Call metadata, success",
			req: &backend.CallResourceRequest{
				Path: "metadata",
			},
			expRespCode: 200,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			body := "<?xml version=\"1.0\" encoding=\"utf-8\"?><edmx:Edmx Version=\"4.0\" xmlns:edmx=\"http://docs.oasis-open.org/odata/ns/edmx\"></edmx:Edmx>"
			client := &clientMock{
				statusCode: table.expRespCode,
				body:       []byte(body),
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
		})
	}
}

func TestCallResourceMetadata(t *testing.T) {
	tables := []struct {
		name        string
		respXml     odata.Edmx
		respErr     error
		expRespCode int
		expResponse metadataResource
	}{
		{
			name:        "Minimal metadata response",
			respXml:     anOdataEdmx("4.0"),
			respErr:     nil,
			expRespCode: 200,
			expResponse: aMetadataResource(),
		},
		{
			name: "Full metadata response",
			respXml: anOdataEdmx("4.0",
				withDataService(
					withSchema("some-namespace",
						withEntityType("entity-type-name",
							withKey("key-name",
								withPropertyRef("property-name")),
							withProperty("property-name", "property-type")),
						withEntityContainer("entity-container-name",
							withEntitySet("entity-set-name", "some-namespace.entity-set-name"))))),
			respErr:     nil,
			expRespCode: 200,
			expResponse: aMetadataResource(
				withEntityTypeResource("entity-type-name", "some-namespace",
					withPropertyResource("property-name", "property-type")),
				withEntitySetResource("entity-set-name", "some-namespace.entity-set-name")),
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			body, _ := xml.Marshal(table.respXml)
			client := &clientMock{
				body:       body,
				err:        table.respErr,
				statusCode: table.expRespCode,
			}
			im := managerMock{}
			ds := ODataSource{&im}

			is := ODataSourceInstance{client}
			im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)
			crs := callResourceResponseSenderMock{}

			// Act
			err := ds.getMetadata(context.TODO(), &backend.CallResourceRequest{Path: "metadata"}, &crs)

			// Assert
			require.NoError(t, err)
			require.Equal(t, table.expRespCode, crs.csr.Status)

			// Parse crs.csr.Body into a metadataResponse struct
			var resp metadataResource
			err = json.Unmarshal(crs.csr.Body, &resp)

			require.Equal(t, table.expResponse, resp)
		})
	}
}
