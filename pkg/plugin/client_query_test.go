package plugin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

func TestMapFilter(t *testing.T) {
	tables := []struct {
		timeProperty     string
		timeRange        backend.TimeRange
		filterConditions []filterCondition
		expected         string
	}{
		{
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(int32Eq5),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5",
		},
		{
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5 and string eq 'Hello'",
		},
		{
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and string eq ''",
		},
		{
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         " and string eq ''",
		},
	}

	for _, table := range tables {
		var filterString = mapFilter(table.timeProperty, table.timeRange, table.filterConditions)
		assert.Equal(t, table.expected, filterString)
	}
}

func TestBuildQueryUrl(t *testing.T) {
	tables := []struct {
		baseUrl          string
		entitySet        string
		properties       []property
		timeProperty     string
		timeRange        backend.TimeRange
		filterConditions []filterCondition
		expected         string
	}{
		{
			baseUrl:          "http://localhost:5000",
			entitySet:        "Temperatures",
			properties:       []property{{Name: "Value1", Type: "Edm.Double"}},
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         "http://localhost:5000/Temperatures?%24filter=time+ge+2022-04-21T12%3A30%3A50Z+and+time+le+2022-04-21T12%3A30%3A50Z+and+string+eq+%27%27&%24select=Value1%2C+time",
		},
	}

	for _, table := range tables {
		var builtUrl, err = buildQueryUrl(table.baseUrl, table.entitySet, table.properties, table.timeProperty,
			table.timeRange, table.filterConditions)
		assert.NoError(t, err)
		assert.Equal(t, table.expected, builtUrl.String())
	}
}

func TestGetEntities(t *testing.T) {
	tables := []struct {
		name            string
		expectedResult  odata.Response
		expectedError   error
		handlerCallback func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:           "Success",
			expectedResult: anOdataResponse(withEntity(withProp("hello", "world"))),
			expectedError:  nil,
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"value\":[{\"hello\":\"world\"}]}"))
			},
		},
		{
			name:          "Invalid json",
			expectedError: &json.SyntaxError{},
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Invalid json"))
			},
		},
		{
			name:          "Server Timeout",
			expectedError: &url.Error{},
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Second)
			},
		},
		{
			name:          "Server 500 error",
			expectedError: errors.New("get failed with status code 500"), // Is of type "errors.errorString"
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			client := GetOC("*", table.handlerCallback)

			// Act
			resp, err := client.GetEntities("Temperatures", []property{{Name: "Value1", Type: "Edm.Double"}}, "time", aOneDayTimeRange(), someFilterConditions(int32Eq5))

			// Assert
			if table.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, &table.expectedResult, resp)
			} else {
				assert.Error(t, err)
				assert.IsType(t, table.expectedError, err)
			}
		})
	}
}

func TestGetMetadata(t *testing.T) {
	tables := []struct {
		name            string
		expectedResult  odata.Edmx
		expectedError   error
		handlerCallback func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:           "Success",
			expectedResult: anOdataEdmx("4.0"),
			expectedError:  nil,
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><edmx:Edmx Version=\"4.0\" xmlns:edmx=\"http://docs.oasis-open.org/odata/ns/edmx\"></edmx:Edmx>"))
			},
		},
		{
			name:          "Invalid xml",
			expectedError: errors.New(""), // Is of type "errors.errorString",
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Invalid xml")) // Produces no syntax error but an eof error as no elements are present
			},
		},
		{
			name:          "Invalid xml, syntax error",
			expectedError: &xml.SyntaxError{},
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><edmx:Edmx Version=\"4.0\" xmlns:edmx=\"http://docs.oasis-open.org/odata/ns/edmx\">")) // Missing closing element
			},
		},
		{
			name:          "Server Timeout",
			expectedError: &url.Error{},
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Second)
			},
		},
		{
			name:          "Server 500 error",
			expectedError: errors.New("get failed with status code 500"), // Is of type "errors.errorString"
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Arrange
			client := GetOC("*", table.handlerCallback)

			// Act
			resp, err := client.GetMetadata()

			// Assert
			if table.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, &table.expectedResult, resp)
			} else {
				assert.Error(t, err)
				assert.IsType(t, table.expectedError, err)
			}
		})
	}
}
