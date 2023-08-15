package giu

import "strings"

func CheckParam(method string, pattern string) {

	var IsSearched bool

	for _, Method := range httpMethod {
		if Method == method {
			IsSearched = true
			break
		}
	}

	if !IsSearched {
		panic("invalid http request method")
	}

	if pattern[0] != '/' && pattern[1] != '*' {
		panic("path must be started with '/' and '*'")
	}

}

func parsePattern(pattern string) []string {
	split := strings.Split(pattern, "/")

	parts := make([]string, 0)

	for i := 0; i < len(split); i++ {
		if split[i] != "" {
			parts = append(parts, split[i])
			if split[i][0] == '*' {
				break
			}
		}
	}
	return parts

}
