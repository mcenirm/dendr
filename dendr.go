package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileEntry struct {
	path  string
	size  int64
	mtime time.Time
}

type inventoryReader struct {
	f *os.File
	e error
}

func newFileListReader(fileName string) *inventoryReader {
	f, err := os.Open(fileName)
	return &inventoryReader{f, err}
}

func (r *inventoryReader) Close() error {
	return r.f.Close()
}

func inventoryFileNameFor(name string) string {
	return name + ".inventory"
}

func (fe *fileEntry) comparePath(path string) int {
	if fe == nil {
		return -1
	}
	return strings.Compare(fe.path, path)
}

func (r *inventoryReader) readEntry() *fileEntry {
	return nil
}

func main() {
	var err error

	start := "testpath"

	pastName := "past"
	pastInventoryFileName := inventoryFileNameFor(pastName)
	pastInventoryReader := newFileListReader(pastInventoryFileName)

	past := pastInventoryReader.readEntry()
	err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// ignore errors (TODO unless verbose)
			return nil
		}
		size := info.Size()
		mtime := info.ModTime()
		cmp := past.comparePath(path)
		switch {
		case cmp == 0:
			sameSize := past.size == size
			sameMtime := past.mtime.Equal(mtime)
			if sameSize && sameMtime {
				// do nothing?
			} else {
				fmt.Print("=")
				if sameSize {
					fmt.Print(".")
				} else {
					fmt.Print("s")
				}
				if sameMtime {
					fmt.Print(".")
				} else {
					fmt.Print("m")
				}
				fmt.Println("  ", path)
			}
		case cmp < 0:
			fmt.Println("compare", cmp, past, path)
		case cmp > 0:
			fmt.Println("compare", cmp, past, path)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking: %v\n", err)
		return
	}
}
