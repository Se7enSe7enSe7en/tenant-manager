package utils

import "strings"

func ClassNameJoin(val ...string) string {
	return strings.Join(val, " ")
}
