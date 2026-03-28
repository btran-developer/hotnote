package slugify

import (
	"strings"
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	s = result.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}

func SlugifyPath(path string) string {
	if path == "" {
		return ""
	}
	parts := strings.Split(path, "/")
	var result []string
	for _, part := range parts {
		slugged := Slugify(part)
		if slugged != "" {
			result = append(result, slugged)
		}
	}
	return strings.Join(result, "/")
}
