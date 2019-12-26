package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type AnyMap = map[string]interface{}

func AnyToAnyMap(value interface{}) AnyMap {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case AnyMap:
		return val
	case map[interface{}]interface{}:
		count := len(val)
		if count == 0 {
			return nil
		}
		m := make(AnyMap, count)
		for k, v := range val {
			m[AnyToString(k)] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch val := value.(type) {
	case *string:
		return *val
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	default:
		return fmt.Sprint(value)
	}
}

func AnyToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case *string:
		if i, err := StringToInt64(*val); err == nil {
			return i
		}
	case string:
		if i, err := StringToInt64(val); err == nil {
			return i
		}
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case bool:
		if val {
			return 1
		} else {
			return 0
		}
	case json.Number:
		v, _ := val.Int64()
		return v
	}
	return 0
}

func AnyToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}
	switch val := value.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	case *string:
		if v, err := strconv.ParseFloat(*val, 64); err == nil {
			return v
		}
	case string:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}
	case bool:
		if val {
			return 1
		} else {
			return 0
		}
	case json.Number:
		v, _ := val.Float64()
		return v
	}
	return 0
}

func AnyToBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		if len(v) == 0 {
			return false
		}
		c := strings.ToLower(v[0:1])
		return c == "y" || c == "t" || c == "1"
	case *string:
		return AnyToBool(*v)
	default:
		return false
	}
}

func AnyArrayToMap(mapInterface []interface{}) AnyMap {
	if len(mapInterface)/2 < 1 {
		return nil
	}
	elementMap := make(AnyMap)
	for i := 0; i < len(mapInterface)/2; i += 1 {
		key := AnyToString(mapInterface[i*2])
		elementMap[key] = mapInterface[i*2+1]
	}
	return elementMap
}

func AnyArrayToStringArray(arrInterface []interface{}) []string {
	elementArray := make([]string, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToString(v)
	}
	return elementArray
}

func StringToInt64(value string) (int64, error) {
	if index := strings.Index(value, "."); index > 0 {
		value = value[:index]
	}
	return strconv.ParseInt(value, 10, 64)
}

func FindInAnyMap(m AnyMap, keys ...string) interface{} {
	return FindInAnyMapWithKeys(m, keys)
}

func FindInAnyMapWithKeys(m AnyMap, keys []string) interface{} {
	if m == nil {
		return nil
	}
	l := len(keys)
	if l == 0 {
		return nil
	}
	value := m[keys[0]]
	if l == 1 {
		return value
	}
	m = AnyToAnyMap(value)
	return FindInAnyMapWithKeys(m, keys[1:])
}
