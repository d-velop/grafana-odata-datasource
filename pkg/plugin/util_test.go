package plugin

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
)

func TestBackendTimeRangeToODataFilter(t *testing.T) {
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
			result := BackendTimeRangeToODataFilter(table.timeRange, table.timeProperty)

			// Assert
			assert.Equal(t, table.expected, result)
		})
	}
}
