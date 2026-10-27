package util

import (
	"encoding/xml"
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

func GetAttr(start xml.StartElement, name xml.Name) (string, bool) {
	for _, attr := range start.Attr {
		if attr.Name == name {
			return attr.Value, true
		}
	}
	return "", false
}
