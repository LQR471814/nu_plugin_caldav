package main

import (
	"reflect"
	"strings"

	"github.com/zeebo/xxh3"
)

func typeId(t reflect.Type) uint64 {
	return xxh3.Hash([]byte(t.String()))
}

func pascalToSnakeCase(pascalcase string) string {
	var result strings.Builder
	for i, c := range pascalcase {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result.WriteString("_")
			}
			result.WriteRune(c + ('a' - 'A'))
			continue
		}
		result.WriteRune(c)
	}
	return result.String()
}
