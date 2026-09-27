package plugin

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

func TestTimeRangeToFilter(t *testing.T) {
	tables := []struct {
		name         string
		timeProperty *property
		timeRange    backend.TimeRange
		expected     []filterCondition
	}{
		{
			name: "Time property set",
			timeProperty: &property{
				Name: "time",
				Type: "Edm.DateTimeOffset",
			},
			timeRange: aOneDayTimeRange(),
			expected: someFilterConditions(
				withFilterCondition(timeProp, "ge", aOneDayTimeRange().From.Format(time.RFC3339)),
				withFilterCondition(timeProp, "le", aOneDayTimeRange().To.Format(time.RFC3339)),
			),
		},
		{
			name:         "No time property set",
			timeProperty: nil,
			timeRange:    aOneDayTimeRange(),
			expected:     []filterCondition{},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			result := TimeRangeToFilter(table.timeRange, table.timeProperty)

			// Assert
			assert.Equal(t, table.expected, result)
		})
	}
}

func TestCompleteProperties(t *testing.T) {
	tables := []struct {
		name       string
		properties []property
		expected   []property
	}{
		{
			name:       "Keeps named properties",
			properties: []property{aProperty(int32Prop), aProperty(stringProp)},
			expected:   []property{aProperty(int32Prop), aProperty(stringProp)},
		},
		{
			name:       "Drops properties without name",
			properties: []property{aProperty(), aProperty(int32Prop), aProperty()},
			expected:   []property{aProperty(int32Prop)},
		},
		{
			name:       "Only properties without name",
			properties: []property{aProperty()},
			expected:   nil,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			result := CompleteProperties(table.properties)

			// Assert
			assert.Equal(t, table.expected, result)
		})
	}
}

func TestCompleteFilterConditions(t *testing.T) {
	tables := []struct {
		name       string
		conditions []filterCondition
		expected   []filterCondition
	}{
		{
			name:       "Keeps complete conditions",
			conditions: someFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
			expected:   someFilterConditions(int32Eq5, withFilterCondition(stringProp, "eq", "Hello")),
		},
		{
			name:       "Drops condition without property",
			conditions: someFilterConditions(withFilterCondition(func(p *property) {}, "eq", "5"), int32Eq5),
			expected:   someFilterConditions(int32Eq5),
		},
		{
			name:       "Drops condition without operator",
			conditions: someFilterConditions(withFilterCondition(int32Prop, "", "5")),
			expected:   nil,
		},
		{
			name:       "Drops non-string condition without value",
			conditions: someFilterConditions(withFilterCondition(int32Prop, "eq", "")),
			expected:   nil,
		},
		{
			name:       "Keeps string condition without value",
			conditions: someFilterConditions(withFilterCondition(stringProp, "eq", "")),
			expected:   someFilterConditions(withFilterCondition(stringProp, "eq", "")),
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			// Act
			result := CompleteFilterConditions(table.conditions)

			// Assert
			assert.Equal(t, table.expected, result)
		})
	}
}
