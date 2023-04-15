package plugin

import (
	"context"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryData(t *testing.T) {
	tables := []struct {
		query    backend.QueryDataRequest
		mockBody string
		expected backend.QueryDataResponse
	}{
		{
			query: aQueryDataRequest(withDataQuery(withQueryModel(withFilterConditions(int32Eq5),
				withProperties(int32Prop, booleanProp, stringProp)))),
			mockBody: "{\"@odata.context\":\"http://localhost:6000/odata/$metadata#Temperatures\",\"value\":" +
				"[{\"int32\":5,\"time\":\"2022-01-02T00:00:00Z\",\"boolean\":true,\"string\":\"Hello World!\"}]}",
			expected: aQueryDataResponse(withDataResponse(withDefaultTestFrame())),
		},
		{
			query: aQueryDataRequest(
				withDataQuery(withQueryModel(
					withFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
					withProperties(int32Prop, booleanProp, stringProp))),
			),
			mockBody: "",
			expected: aQueryDataResponse(
				withDataResponse(withDefaultTestFrame()),
			),
		},
	}

	im := managerMock{}
	ds := ODataSource{&im}

	for _, table := range tables {
		client := clientMock{statusCode: 200, body: table.mockBody}
		is := ODataSourceInstance{&client}
		im.On("Get", mock.Anything).Return(&is, nil)

		// Result
		result, err := ds.QueryData(context.TODO(), &table.query)
		assert.NoError(t, err)
		assert.Equal(t, result.Responses, table.expected.Responses)
	}
}

func TestQuery(t *testing.T) {
	tables := []struct {
		query          backend.DataQuery
		mockStatusCode int
		mockBody       string
		expected       backend.DataResponse
	}{
		{
			query:          aDataQuery(withQueryModel()),
			mockStatusCode: 200,
			mockBody:       "{\"@odata.context\":\"http://localhost:6000/odata/$metadata#Temperatures\",\"value\":[]}",
			expected:       aDataResponse(withBaseFrame()),
		},
	}

	im := managerMock{}
	ds := ODataSource{&im}

	for _, table := range tables {
		client := clientMock{statusCode: table.mockStatusCode, body: table.mockBody}
		resp := ds.query(&client, table.query)
		require.Equal(t, table.expected, resp)
	}
}
