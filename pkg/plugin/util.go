package plugin

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TimeRangeToFilter(timeRange backend.TimeRange, timeProperty *property) []filterCondition {
	if timeProperty == nil {
		return []filterCondition{}
	}

	return []filterCondition{
		{
			Property: *timeProperty,
			Operator: "ge",
			Value:    timeRange.From.UTC().Format(time.RFC3339),
		},
		{
			Property: *timeProperty,
			Operator: "le",
			Value:    timeRange.To.UTC().Format(time.RFC3339),
		},
	}
}
