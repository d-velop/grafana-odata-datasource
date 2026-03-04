package odata

import (
	"fmt"
	"time"
)

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
		result, err := mapNumber(value.(float64), propertyType)
		if err != nil {
			return nil
		}
		return result
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

func mapNumber(value float64, propertyType string) (interface{}, error) {
	switch propertyType {
	case EdmSingle:
		y := float32(value)
		return &y, nil
	case EdmDecimal, EdmDouble:
		return &value, nil
	case EdmSByte:
		y := int8(value)
		return &y, nil
	case EdmByte:
		y := uint8(value)
		return &y, nil
	case EdmInt16:
		y := int16(value)
		return &y, nil
	case EdmInt32:
		y := int32(value)
		return &y, nil
	case EdmInt64:
		y := int64(value)
		return &y, nil
	default:
		return nil, fmt.Errorf("unexpected property type: %s", propertyType)
	}
}
