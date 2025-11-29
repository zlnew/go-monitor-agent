// Package pkg
package pkg

import "strings"

func ContainsAny(s string, xs []string) bool {
	for _, x := range xs {
		if strings.Contains(s, strings.ToLower(x)) {
			return true
		}
	}

	return false
}
