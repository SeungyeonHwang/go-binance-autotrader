package utilities

import (
	"strconv"
	"strings"
)

func FormatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 5, 64)

	for strings.HasSuffix(s, "0") && strings.Contains(s, ".") {
		s = s[:len(s)-1]
	}

	s = strings.TrimSuffix(s, ".")

	return s
}
