package components

import (
	"fmt"
	"strings"
)

func ComponentID(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func ToString[T any](v *T, fallback string) string {
	if v == nil {
		return fallback
	}

	return fmt.Sprintf("%v", *v)
}
