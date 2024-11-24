package utils

import "strconv"

func ParseInt(s string, defaultValue int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultValue
}
