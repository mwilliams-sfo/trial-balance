package util

import (
	"strings"
)

func IsXmlSpace(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("\t\n\r ", r) {
			return false
		}
	}
	return true
}
