package ssg

import (
	"regexp"
	"strings"
)

// NormalizeSlug sanitizes a slug to ensure it's URL-safe:
// - Lowercase
// - Replace spaces with hyphens
// - Remove non-alphanumeric except hyphens
// - Collapse consecutive hyphens
// - Trim leading/trailing hyphens
func NormalizeSlug(input string) string {
	s := strings.ToLower(input)
	s = strings.ReplaceAll(s, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}
