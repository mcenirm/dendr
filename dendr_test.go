package main

import (
	"os"
	"testing"
)

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

func TestOpenPastInventoryReaderStdin(t *testing.T) {
	r := openPastInventoryReader("-")
	defer r.Close()

	if r.e != nil {
		t.Error(r.e)
	}
	if r.f != os.Stdin {
		t.Errorf("openPastInventoryReader(\"-\") should open stdin")
	}
}

func TestOpenNextInventoryWriterStdout(t *testing.T) {
	r := openNextInventoryWriter("-")
	defer r.Close()

	if r.e != nil {
		t.Error(r.e)
	}
	if r.f != os.Stdout {
		t.Errorf("openNextInventoryWriter(\"-\") should open stdout")
	}
}
