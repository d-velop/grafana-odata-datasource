package odata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateTimeWithoutTZ = "2006-01-02T15:04:05"

// ToArray maps OData property types to Grafana Field type
func ToArray(propertyType string) interface{} {
	switch propertyType {
	case EdmBoolean:
		return []*bool{}
	case EdmSingle:
		return []*float32{}
	case EdmDouble:
		return []*float64{}
	case EdmDecimal:
		return []*float64{}
	case EdmSByte:
		return []*int8{}
	case EdmByte:
		return []*uint8{}
	case EdmInt16:
		return []*int16{}
	case EdmInt32:
		return []*int32{}
	case EdmInt64:
		return []*int64{}
	case EdmDateTimeOffset:
		return []*time.Time{}
	case EdmDate:
		return []*time.Time{}
	default:
		return []*string{}
	}
}

func parseOffset(s string, sign int) (int, error) {
	offset, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return sign * offset * 60, nil
}

func localOffset() int {
	currentTime := time.Now()
	_, offset := currentTime.Zone()
	return offset
}

func parseV2Time(timeString string) (time.Time, error) {
	trimmed := strings.TrimSuffix(strings.TrimPrefix(timeString, "/Date("), ")/")
	var err error
	var parts []string
	var offset int
	if strings.Contains(trimmed, "+") {
		parts = strings.Split(trimmed, "+")
		offset, err = parseOffset(parts[1], 1)
	} else if strings.Contains(trimmed, "-") {
		parts = strings.Split(trimmed, "-")
		offset, err = parseOffset(parts[1], -1)
	} else {
		parts = []string{trimmed}
	}
	if err != nil {
		return time.Time{}, err
	}
	ms, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	result := time.Unix(seconds, nanoseconds).Add(time.Duration(offset) * time.Second)
	var loc *time.Location
	if localOffset() == offset {
		loc = time.Local
	} else {
		loc = time.FixedZone("", offset)
	}
	return result.In(loc), nil
}

func ParseTime(timeString string) (time.Time, error) {
	if strings.HasPrefix(timeString, "/") {
		ts, err := parseV2Time(timeString)
		if err == nil && !ts.IsZero() {
			return ts, nil
		}
	}
	formats := []string{
		time.RFC3339Nano,
		time.DateOnly,
	}
	var ts time.Time
	var err error
	for _, format := range formats {
		ts, err = time.Parse(format, timeString)
		if err == nil && !ts.IsZero() {
			return ts, nil
		}
	}
	return time.ParseInLocation(DateTimeWithoutTZ, timeString, time.UTC)
}

// MapValue maps OData values to Grafana (Go) values
func MapValue(value interface{}, propertyType string) interface{} {
	if value == nil {
		return nil
	}
	switch propertyType {
	case EdmBoolean:
		boolValue := value.(bool)
		return &boolValue
	case EdmSingle, EdmDecimal, EdmDouble, EdmSByte, EdmByte, EdmInt16, EdmInt32, EdmInt64:
		if stringValue, ok := value.(string); ok {
			// TODO: work in progress
			floatValue, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				panic("could not parse number value in string")
			}
			return mapNumber(floatValue, propertyType)
		} else if floatValue, ok := value.(float64); ok {
			return mapNumber(floatValue, propertyType)
		} else {
			// TODO: fall back to string?
			x := fmt.Sprint(value)
			return &x
		}
	case EdmDateTimeOffset, EdmDate:
		if timeValue, err := time.Parse(time.RFC3339Nano, fmt.Sprint(value)); err == nil {
			return &timeValue
		} else {
			return nil
		}
	default:
		x := fmt.Sprint(value)
		return &x
	}
}

func mapNumber(value float64, propertyType string) interface{} {
	switch propertyType {
	case EdmSingle:
		y := float32(value)
		return &y
	case EdmDecimal, EdmDouble:
		return &value
	case EdmSByte:
		y := int8(value)
		return &y
	case EdmByte:
		y := uint8(value)
		return &y
	case EdmInt16:
		y := int16(value)
		return &y
	case EdmInt32:
		y := int32(value)
		return &y
	case EdmInt64:
		y := int64(value)
		return &y
	default:
		panic("unexpected property type")
	}
}

func MapToResponse(bodyBytes []byte) ([]interface{}, error) {
	var response Response
	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	if response.Value != nil {
		return response.Value, nil
	} else if response.D != nil {
		switch d := response.D.(type) {
		case map[string]interface{}:
			if results, ok := d["results"].([]interface{}); ok {
				return results, nil
			}
		case []interface{}:
			return d, nil
		}
	} else if response.Results != nil {
		return response.Results, nil
	}
	return nil, nil
}
