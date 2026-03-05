package utils

import "strings"

/*
This function is for turning any value into a pointer,
our main usecase for this as of the moment is to handle optional
or nullable behavior for the front end.
*/
func Ptr[T any](v T) *T {
	return &v
}

func ClassNameJoin(val ...string) string {
	return strings.Join(val, " ")
}
