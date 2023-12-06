package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Just test if multiple queries lead to multiple responses
func TestQueryData(t *testing.T) {
	tables := []struct {
		name     string
		query    backend.QueryDataRequest
		expected backend.QueryDataResponse
	}{
		{
			name:     "Zero queries",
			query:    aQueryDataRequest(),
			expected: aQueryDataResponse(),
		},
		{
			name:     "One query",
			query:    aQueryDataRequest(withDataQuery("one", withQueryModel())),
			expected: aQueryDataResponse(withDataResponse("one", withDefaultTestFrame())),
		},
		{
			name:  "Two queries",
			query: aQueryDataRequest(withDataQuery("one", withQueryModel()), withDataQuery("two", withQueryModel())),
			expected: aQueryDataResponse(
				withDataResponse("one", withDefaultTestFrame()),
				withDataResponse("two", withDefaultTestFrame())),
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			im := managerMock{}
			ds := ODataSource{&im}

			body, _ := json.Marshal(odata.Response{})
			client := clientMock{body: body}
			is := ODataSourceInstance{&client}
			im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

			// Act
			result, err := ds.QueryData(context.TODO(), &table.query)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, len(table.expected.Responses), len(result.Responses))
		})
	}
}

// TODO: Tests with invalid query models
func TestQuery(t *testing.T) {
	tables := []struct {
		name              string
		mockODataResponse odata.Response
		query             backend.DataQuery
		expected          backend.DataResponse
	}{
		{
			name: "success simple",
			query: aDataQuery("defaultTestFrame", withQueryModel(withTimeProperty("time"),
				withFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
				withProperties(int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(withDefaultEntity()),
			expected:          aDataResponse(withDefaultTestFrame()),
		},
		{
			name: "success ordered",
			query: aDataQuery("defaultTestFrame", withQueryModel(withTimeProperty("time"),
				withProperties(int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(
				withEntity(
					withProp("string", "Hello"),
					withProp("int32", 10.0),
					withProp("boolean", false),
					withProp("time", "2022-01-02T00:00:00Z")),
				withEntity(
					withProp("time", "2000-01-02T00:00:00Z"),
				),
				withEntity(
					withProp("time", "2010-01-02T00:00:00Z"),
					withProp("string", "World"),
				),
			),
			expected: aDataResponse(withBaseFrame("defaultTestFrame",
				withTimeField("time"),
				withField("int32", []*int32{}),
				withField("boolean", []*bool{}),
				withField("string", []*string{}),
				withRow(
					withRowValue(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)),
					withRowValue(int32(10)),
					withRowValue(false),
					withRowValue("Hello"),
				),
				withRow(
					withRowValue(time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)),
					nil, nil, nil,
				),
				withRow(
					withRowValue(time.Date(2010, 1, 2, 0, 0, 0, 0, time.UTC)),
					nil, nil,
					withRowValue("World"),
				),
			)),
		},
		{
			name:  "success select time without time property",
			query: aDataQuery("defaultTestFrame", withQueryModel(withProperties(timeProp, int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(
				withEntity(
					withProp("string", "Hello"),
					withProp("int32", 10.0),
					withProp("boolean", false),
					withProp("time", "2022-01-02T00:00:00Z")),
			),
			expected: aDataResponse(withBaseFrame("defaultTestFrame",
				withTimeField("time"),
				withField("int32", []*int32{}),
				withField("boolean", []*bool{}),
				withField("string", []*string{}),
				withRow(
					withRowValue(time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)),
					withRowValue(int32(10)),
					withRowValue(false),
					withRowValue("Hello"),
				),
			)),
		},
		{
			name:  "success select no time property",
			query: aDataQuery("defaultTestFrame", withQueryModel(withProperties(int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(
				withEntity(
					withProp("string", "Hello"),
					withProp("int32", 10.0),
					withProp("boolean", false)),
			),
			expected: aDataResponse(withBaseFrame("defaultTestFrame",
				withField("int32", []*int32{}),
				withField("boolean", []*bool{}),
				withField("string", []*string{}),
				withRow(
					withRowValue(int32(10)),
					withRowValue(false),
					withRowValue("Hello"),
				),
			)),
		},
		{
			name:              "success minimal",
			query:             aDataQuery("baseFrame", withQueryModel()),
			mockODataResponse: anOdataResponse(),
			expected:          aDataResponse(),
		},
		{
			name:  "success select time property that does not exist",
			query: aDataQuery("defaultTestFrame", withQueryModel(withProperties(timeProp, int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(
				withEntity(
					withProp("string", "Hello"),
					withProp("int32", 10.0),
					withProp("boolean", false),
					withProp("otherTimePropName", "2022-01-02T00:00:00Z")),
			),
			expected: aDataResponse(withBaseFrame("defaultTestFrame",
				withTimeField("time"),
				withField("int32", []*int32{}),
				withField("boolean", []*bool{}),
				withField("string", []*string{}),
				withRow(
					nil,
					withRowValue(int32(10)),
					withRowValue(false),
					withRowValue("Hello"),
				),
			)),
		},
		{
			name: "failure",
			query: aDataQuery("one", withQueryModel(withTimeProperty("time"),
				withFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
				withProperties(int32Prop, booleanProp, stringProp))),
			mockODataResponse: anOdataResponse(),
			expected:          aDataResponse(withErrorResponse(errors.New("Something went wrong"))),
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			im := managerMock{}
			ds := ODataSource{&im}

			body, _ := json.Marshal(table.mockODataResponse)
			client := clientMock{
				body:       body,
				err:        table.expected.Error,
				statusCode: 200,
			}
			is := ODataSourceInstance{&client}
			im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

			// Act
			resp := ds.query(&client, table.query)

			// Assert
			assert.Equal(t, table.expected, resp)
		})
	}
}

func TestInvalidQueryModels(t *testing.T) {
	tables := []struct {
		name             string
		query            backend.DataQuery
		expectedErrorMsg string
	}{
		{
			name: "Invalid json",
			query: backend.DataQuery{
				JSON: []byte(`{`),
			},
			expectedErrorMsg: "error unmarshalling query json",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			im := managerMock{}
			ds := ODataSource{&im}

			client := clientMock{}
			is := ODataSourceInstance{&client}
			im.On("Get", context.TODO(), mock.Anything).Return(&is, nil)

			// Act
			resp := ds.query(&client, table.query)

			// Assert
			assert.NotNil(t, resp.Error)
			assert.Contains(t, resp.Error.Error(), table.expectedErrorMsg)
		})
	}
}
