package plugin

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func BackendTimeRangeToODataFilter(timeRange backend.TimeRange, timeProperty string) []filterCondition {
	if timeProperty == "" {
		return []filterCondition{}
	}

	return []filterCondition{
		{
			Property: property{
				Name: timeProperty,
				Type: "Edm.DateTimeOffset",
			},
			Operator: "ge",
			Value:    timeRange.From.UTC().Format(time.RFC3339),
		},
		{
			Property: property{
				Name: timeProperty,
				Type: "Edm.DateTimeOffset",
			},
			Operator: "le",
			Value:    timeRange.To.UTC().Format(time.RFC3339),
		},
	}
}
