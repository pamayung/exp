package util

import "time"

func IsNil(i interface{}) bool {
	switch i.(type) {
	case string:
		return len(i.(string)) <= 0
	case bool:
		return i.(bool)
	case float64:
		return i.(float64) == 0
	case int:
		return i.(int) == 0
	case int64:
		return i.(int) == 0
	default:
		return true
	}

	return true
}

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
