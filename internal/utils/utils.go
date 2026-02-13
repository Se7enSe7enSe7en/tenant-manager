package utils

import "strings"

func Ptr[T any](v T) *T {
	return &v
}

func ClassNameJoin(val ...string) string {
	return strings.Join(val, " ")
}
