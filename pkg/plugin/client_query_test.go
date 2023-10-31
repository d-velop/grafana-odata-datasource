package plugin

import (
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
		name             string
		timeProperty     string
		timeRange        backend.TimeRange
		filterConditions []filterCondition
		expected         string
	}{
		{
			name:             "Time filter only",
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(int32Eq5),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5",
		},
		{
			name:             "Time filter and int and string filter",
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5 and string eq 'Hello'",
		},
		{
			name:             "Time filter and string filter",
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and string eq ''",
		},
		{
			name:             "String filter only",
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         " and string eq ''",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			var filterString = mapFilter(table.timeProperty, table.timeRange, table.filterConditions)

			// Assert
			assert.Equal(t, table.expected, filterString)
		})
	}
}

func TestBuildQueryUrl(t *testing.T) {
	tables := []struct {
		name             string
		baseUrl          string
		entitySet        string
		properties       []property
		timeProperty     string
		timeRange        backend.TimeRange
		filterConditions []filterCondition
		expected         string
	}{
		{
			name:             "Success",
			baseUrl:          "http://localhost:5000",
			entitySet:        "Temperatures",
			properties:       []property{aProperty(int32Prop)},
			timeProperty:     "time",
			timeRange:        aOneDayTimeRange(),
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         "http://localhost:5000/Temperatures?%24filter=time+ge+2022-04-21T12%3A30%3A50Z+and+time+le+2022-04-21T12%3A30%3A50Z+and+string+eq+%27%27&%24select=int32%2C+time",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			var builtUrl, err = buildQueryUrl(table.baseUrl, table.entitySet, table.properties, table.timeProperty,
				table.timeRange, table.filterConditions)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, table.expected, builtUrl.String())
		})
	}
}

func TestGetEntities(t *testing.T) {
	tables := []struct {
		name             string
		expectedError    error
		expectedRespCode int
		handlerCallback  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:             "Success",
			expectedRespCode: 200,
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"value\":[{\"hello\":\"world\"}]}"))
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
			name:             "Server 500 error",
			expectedRespCode: 500,
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
			resp, err := client.Get("Temperatures", []property{aProperty(int32Prop)}, "time", aOneDayTimeRange(), someFilterConditions(int32Eq5))

			// Assert
			if table.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, table.expectedRespCode, resp.StatusCode)
			} else {
				assert.Error(t, err)
				assert.IsType(t, table.expectedError, err)
			}
		})
	}
}

func TestGetMetadata(t *testing.T) {
	tables := []struct {
		name             string
		expectedResult   odata.Edmx
		expectedRespCode int
		expectedError    error
		handlerCallback  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:             "Success",
			expectedResult:   anOdataEdmx("4.0"),
			expectedError:    nil,
			expectedRespCode: 200,
			handlerCallback: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><edmx:Edmx Version=\"4.0\" xmlns:edmx=\"http://docs.oasis-open.org/odata/ns/edmx\"></edmx:Edmx>"))
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
			name:             "Server 500 error",
			expectedRespCode: 500,
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
				assert.Equal(t, table.expectedRespCode, resp.StatusCode)
			} else {
				assert.IsType(t, table.expectedError, err)
			}
		})
	}
}
