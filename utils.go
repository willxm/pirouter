package pirouter

import "strings"

func TrimPathPrefix(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

func SplitPattern(pattern string) []string {
	return strings.Split(pattern, "/")
}
