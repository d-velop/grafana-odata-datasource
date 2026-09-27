package plugin

import (
	"time"

	"github.com/d-velop/grafana-odata-datasource/pkg/plugin/odata"
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

func CompleteProperties(properties []property) []property {
	var result []property
	for _, p := range properties {
		if p.Name != "" {
			result = append(result, p)
		}
	}
	return result
}

func CompleteFilterConditions(conditions []filterCondition) []filterCondition {
	var result []filterCondition
	for _, c := range conditions {
		if c.Property.Name == "" || c.Operator == "" {
			continue
		}
		if c.Value == "" && c.Property.Type != odata.EdmString {
			continue
		}
		result = append(result, c)
	}
	return result
}
