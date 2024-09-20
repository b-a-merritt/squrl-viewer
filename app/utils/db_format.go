package utils

import (
	"time"
)

func FormatFromDB(entry interface{}) interface{} {
	switch v := entry.(type) {
	default:
		return v
	case time.Time:
		return v.Format(time.RFC3339)
	case nil:
		return "null"
	case []uint8:
		str := ""
		for _, val := range entry.([]uint8) {
			str += string(byte(val))
		}
		return str
	}
}
