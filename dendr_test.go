package main

import "testing"

func TestComparePath(t *testing.T) {
	before := "/a/before"
	test := "/b/sample"
	after := "/c/after"

	a := func(fe *fileEntry, path string, expected int) {
		actual := fe.comparePath(path)
		if actual != expected {
			t.Errorf("expected %v, got %v, for %v", expected, actual, fe)
		}
	}

	a(&fileEntry{path: before}, test, -1)
	a(&fileEntry{path: test}, test, 0)
	a(&fileEntry{path: after}, test, 1)
	a(nil, test, -1)
}
