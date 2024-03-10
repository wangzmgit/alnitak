package utils

import (
	"fmt"
	"strings"
)

func IntJoin(elems []int, sep string) string {
	str := make([]string, len(elems))
	for i, v := range elems {
		str[i] = fmt.Sprint(v)
	}

	return strings.Join(str, ",")
}

func UintJoin(elems []uint, sep string) string {
	str := make([]string, len(elems))
	for i, v := range elems {
		str[i] = fmt.Sprint(v)
	}

	return strings.Join(str, ",")
}
