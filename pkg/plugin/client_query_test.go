package plugin

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
	"github.com/stretchr/testify/assert"
)

func TestMapFilter(t *testing.T) {
	tables := []struct {
		name             string
		version          string
		filterConditions []filterCondition
		expected         string
	}{
		// --- V4 / default behavior ---
		{
			name: "V4: DateTimeOffset filter",
			filterConditions: someFilterConditions(
				withFilterCondition(timeProp, "ge", aOneDayTimeRange().From.Format(time.RFC3339)),
				withFilterCondition(timeProp, "le", aOneDayTimeRange().To.Format(time.RFC3339)),
				int32Eq5),
			expected: "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5",
		},
		{
			name: "V4: DateTimeOffset and string filter",
			filterConditions: someFilterConditions(
				withFilterCondition(timeProp, "ge", aOneDayTimeRange().From.Format(time.RFC3339)),
				withFilterCondition(timeProp, "le", aOneDayTimeRange().To.Format(time.RFC3339)),
				int32Eq5,
				withFilterCondition(stringProp, "eq", "Hello")),
			expected: "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and int32 eq 5 and string eq 'Hello'",
		},
		{
			name: "V4: DateTimeOffset and empty string filter",
			filterConditions: someFilterConditions(
				withFilterCondition(timeProp, "ge", aOneDayTimeRange().From.Format(time.RFC3339)),
				withFilterCondition(timeProp, "le", aOneDayTimeRange().To.Format(time.RFC3339)),
				withFilterCondition(stringProp, "eq", "")),
			expected: "time ge 2022-04-21T12:30:50Z and time le 2022-04-21T12:30:50Z and string eq ''",
		},
		{
			name:             "V4: String filter only",
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         "string eq ''",
		},
		// --- V2 behavior ---
		{
			name:    "V2: DateTimeOffset wraps with datetimeoffset prefix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(timeProp, "ge", "2022-04-21T12:30:50Z"),
				withFilterCondition(timeProp, "le", "2022-04-21T12:30:50Z")),
			expected: "time ge datetimeoffset'2022-04-21T12:30:50Z' and time le datetimeoffset'2022-04-21T12:30:50Z'",
		},
		{
			name:    "V2: DateTime wraps with datetime prefix and strips timezone",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "ts"; p.Type = odata.EdmDateTime }, "ge", "2022-04-21T12:30:50Z")),
			expected: "ts ge datetime'2022-04-21T12:30:50'",
		},
		{
			name:    "V2: Int64 gets L suffix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "count"; p.Type = odata.EdmInt64 }, "eq", "42")),
			expected: "count eq 42L",
		},
		{
			name:    "V2: Decimal gets M suffix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "amount"; p.Type = odata.EdmDecimal }, "eq", "12.34")),
			expected: "amount eq 12.34M",
		},
		{
			name:    "V2: Single gets f suffix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "discount"; p.Type = odata.EdmSingle }, "gt", "0")),
			expected: "discount gt 0f",
		},
		{
			name:    "V2: Double gets d suffix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "ratio"; p.Type = odata.EdmDouble }, "gt", "1.5")),
			expected: "ratio gt 1.5d",
		},
		{
			name:    "V2: Guid gets guid prefix",
			version: "V2",
			filterConditions: someFilterConditions(
				withFilterCondition(func(p *property) { p.Name = "id"; p.Type = odata.EdmGuid }, "eq", "12345678-1234-1234-1234-123456789abc")),
			expected: "id eq guid'12345678-1234-1234-1234-123456789abc'",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			var filterString = mapFilter(table.filterConditions, table.version)

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
		timeRange        []filterCondition
		filterConditions []filterCondition
		expected         string
	}{
		{
			name:       "Success",
			baseUrl:    "http://localhost:5000",
			entitySet:  "Temperatures",
			properties: []property{aProperty(int32Prop), aProperty(timeProp)},
			filterConditions: someFilterConditions(
				withFilterCondition(timeProp, "ge", aOneDayTimeRange().From.Format(time.RFC3339)),
				withFilterCondition(timeProp, "le", aOneDayTimeRange().To.Format(time.RFC3339)),
				withFilterCondition(stringProp, "eq", "")),
			expected: "http://localhost:5000/Temperatures?%24filter=time+ge+2022-04-21T12%3A30%3A50Z+and+time+le+2022-04-21T12%3A30%3A50Z+and+string+eq+%27%27&%24select=int32%2Ctime",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			var builtUrl, err = buildQueryUrl(table.baseUrl, table.entitySet, table.properties, table.filterConditions, "+", "")

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
				_, _ = w.Write([]byte("{\"value\":[{\"hello\":\"world\"}]}"))
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
			resp, err := client.Get(context.TODO(), "Temperatures", []property{aProperty(int32Prop)}, someFilterConditions(int32Eq5))

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
				_, _ = w.Write([]byte("<?xml version=\"1.0\" encoding=\"utf-8\"?><edmx:Edmx Version=\"4.0\" xmlns:edmx=\"https://docs.oasis-open.org/odata/ns/edmx\"></edmx:Edmx>"))
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
			resp, err := client.GetMetadata(context.TODO())

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
