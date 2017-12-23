package utils

import (
	"strconv"
	"strings"
)

func StrToInt(s string) int {
	val, err := strconv.Atoi(strings.TrimSpace(strings.Replace(s, ",", "", len(s))))
	if err != nil {
		return 0
	}
	return val
}
