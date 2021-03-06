package main

import (
	"testing"
)

func TestFormat(t *testing.T) {
	cases := [][]string{
		{"testCase", "test-case"},
		{"TestCase", "test-case"},
		{"Test Case", "test-case"},
		{" Test Case", "test-case"},
		{"Test Case ", "test-case"},
		{" Test Case ", "test-case"},
		{"test", "test"},
		{"test-case", "test-case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many-many-words"},
		{"manyManyWords", "many-many-words"},
		{"AnyKind of-string", "any-kind-of-string"},
		{"numbers2and55with000", "numbers2and55with000"},
		{"JSONData", "jsondata"},
		{"userID", "user-id"},
		{"AAAbbb", "aaabbb"},
		{"(test) case", "test-case"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := Format(in)
		if result != out {
			t.Errorf("'%s' ('%s' != '%s')", in, result, out)
		}
	}
}
