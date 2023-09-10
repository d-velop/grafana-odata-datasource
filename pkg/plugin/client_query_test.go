package plugin

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapFilterWithTimeRange(t *testing.T) {
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
	}

	for _, table := range tables {
		var filterString = mapFilter(table.timeProperty, table.timeRange, table.filterConditions)
		assert.Equal(t, table.expected, filterString)
	}
}

func TestMapFilterWithoutTimeRange(t *testing.T) {
	tables := []struct {
		filterConditions []filterCondition
		expected         string
	}{
		{
			filterConditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:         " and string eq ''",
		},
	}

	for _, table := range tables {
		var filterString = mapFilter("", backend.TimeRange{}, table.filterConditions)
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
			expected:         "http://localhost:5000/Temperatures?%24filter=time+ge+2022-04-21T12%3A30%3A50Z+and+time+le+2022-04-21T12%3A30%3A50Z+and+string+eq+%27%27&%24select=Value1%2Ctime",
		},
	}

	for _, table := range tables {
		var builtUrl, err = buildQueryUrl(table.baseUrl, table.entitySet, table.properties, table.timeProperty,
			table.timeRange, table.filterConditions, "+")
		assert.NoError(t, err)
		assert.Equal(t, table.expected, builtUrl.String())
	}
}
