package odata

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
	default:
		return []*string{}
	}
}

func parseV2Time(timeString string) (time.Time, error) {
	trimmed := strings.TrimPrefix(timeString, "/Date(")
	trimmed = strings.TrimSuffix(trimmed, ")/")
	var offsetString string
	var parts []string
	if strings.Contains(trimmed, "+") {
		parts = strings.Split(trimmed, "+")
		offsetString = "+" + parts[1]
	} else if strings.Contains(trimmed, "-") {
		parts = strings.Split(trimmed, "-")
		offsetString = "-" + parts[1]
	} else {
		parts = []string{trimmed}
	}
	ms, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	result := time.Unix(seconds, nanoseconds)
	if offsetString != "" {
		offset, err := time.ParseDuration(offsetString + "m")
		if err != nil {
			return time.Time{}, err
		}
		result = result.Add(offset)
	}
	return result, nil
}

func ParseTime(timeString string) (time.Time, error) {
	if strings.HasPrefix(timeString, "/") {
		ts, err := parseV2Time(timeString)
		if err == nil && !ts.IsZero() {
			return ts, nil
		}
	}
	return time.Parse(time.RFC3339Nano, timeString)
}

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
