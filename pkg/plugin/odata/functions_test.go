package odata

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		name          string
		input         string
		expectedError error
		expectedTime  time.Time
	}{
		{
			name:         "OData V2",
			input:        "/Date(1672531200000)/",
			expectedTime: defaultTime,
		},
		{
			name:         "OData V2 with positive offset minutes",
			input:        "/Date(1672531200000+0060)/",
			expectedTime: defaultTime.Add(1 * time.Hour),
		},
		{
			name:         "OData V2 with negative offset minutes",
			input:        "/Date(1672531200000-0060)/",
			expectedTime: defaultTime.Add(-1 * time.Hour),
		},
		{
			name:         "RFC3389 base",
			input:        "2023-01-01T00:00:00Z",
			expectedTime: defaultTime,
		},
		{
			name:         "RFC3389 with fractional second",
			input:        "2023-01-01T00:00:00.000Z",
			expectedTime: defaultTime,
		},
		{
			name:         "RFC3389 with offset +01:00",
			input:        "2023-01-01T00:00:00+01:00",
			expectedTime: defaultTime.Add(-1 * time.Hour),
		},
		{
			name:         "RFC3389 with offset -01:00",
			input:        "2023-01-01T00:00:00-01:00",
			expectedTime: defaultTime.Add(1 * time.Hour),
		},
		{
			name:         "Datetime without time zone",
			input:        "2023-01-01T00:00:00",
			expectedTime: defaultTime,
		},
		{
			name:         "Date only",
			input:        "2023-01-01",
			expectedTime: defaultTime,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name+": "+tc.input, func(t *testing.T) {
			ts, err := ParseTime(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTime.UTC(), ts.UTC())
		})
	}
}

func TestResponseMapping(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "",
			input: `{
				"value": [
					{"test": "test"}
				]
			}`,
		},
		{
			name: "",
			input: `{
				"d": {
					"results": [
						{"test": "test"}
					]
				}
			}`,
		},
		{
			name: "",
			input: `{
				"d": [
					{"test": "test"}
				]
			}`,
		},
		{
			name: "",
			input: `{
				"results": [
					{"test": "test"}
				]
			}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := MapToResponse([]byte(tc.input))
			require.NoError(t, err)
			for _, entry := range entries {
				object, ok := entry.(map[string]interface{})
				require.Equal(t, true, ok)
				require.Equal(t, "test", object["test"])
			}
		})
	}
}
