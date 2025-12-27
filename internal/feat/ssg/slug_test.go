package ssg_test

import (
	"testing"

	"github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNormalizeSlug(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"EstoNoEsValido", "estonoesvalido"},
		{"ESTO---tampoco", "esto-tampoco"},
		{"esto-si", "esto-si"},
		{"esto", "esto"},
		{"Mi Sitio Personal", "mi-sitio-personal"},
		{"Blog!!!2024", "blog2024"},
		{"---leading-trailing---", "leading-trailing"},
		{"Hello World", "hello-world"},
		{"", ""},
		{"123-test-456", "123-test-456"},
		{"UPPERCASE", "uppercase"},
		{"special@#$chars", "specialchars"},
		{"multiple   spaces", "multiple-spaces"},
	}

	for _, tc := range cases {
		result := ssg.NormalizeSlug(tc.input)
		if result != tc.expected {
			t.Errorf("NormalizeSlug(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
