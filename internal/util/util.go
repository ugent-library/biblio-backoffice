package util

func ParseBoolean(v interface{}) bool {
	switch b := v.(type) {
	case int32:
		return b == 1
	case int64:
		return b == 1
	case string:
		return b == "true" || b == "1"
	}
	return false
}
