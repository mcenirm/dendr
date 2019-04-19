package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
	s *bufio.Scanner
	e error
}

type inventoryWriter struct {
	f *os.File
	e error
}

func newInventoryReader(fileName string) *inventoryReader {
	f, err := os.Open(fileName)
	s := bufio.NewScanner(f)
	return &inventoryReader{f, s, err}
}

func newInventoryWriter(fileName string) *inventoryWriter {
	f, err := os.Create(fileName)
	return &inventoryWriter{f, err}
}

func (r *inventoryReader) Close() error {
	return r.f.Close()
}

func (w *inventoryWriter) Close() error {
	return w.f.Close()
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

const (
	inventoryFieldSep    string = "\t"
	inventoryMarkerSize         = 's'
	inventoryMarkerMtime        = 't'
	inventoryTimeLayout         = time.RFC3339Nano
	inventoryFormat             = "%v" + inventoryFieldSep + string(inventoryMarkerSize) + "%v" + inventoryFieldSep + string(inventoryMarkerMtime) + "%v\n"
)

func (r *inventoryReader) readEntry() *fileEntry {
	if r.s.Scan() {
		t := r.s.Text()
		fields := strings.Split(t, inventoryFieldSep)
		path, _ := url.PathUnescape(fields[0])
		entry := fileEntry{path: path}
		for _, field := range fields[1:] {
			marker := field[0]
			value := field[1:]
			switch marker {
			case inventoryMarkerSize:
				entry.size, _ = strconv.ParseInt(value, 10, 64)
			case inventoryMarkerMtime:
				mtime, _ := time.Parse(inventoryTimeLayout, value)
				entry.mtime = mtime.UTC()
			}
		}
		return &entry
	}
	return nil
}

func (w *inventoryWriter) writeEntry(fe *fileEntry) {
	path := url.PathEscape(fe.path)
	mtime := fe.mtime.UTC().Format(inventoryTimeLayout)
	fmt.Fprintf(w.f, inventoryFormat, path, fe.size, mtime)
}

func main() {
	var err error

	start := "testpath"

	pastName := "past"
	pastInventoryFileName := inventoryFileNameFor(pastName)
	pastInventoryReader := newInventoryReader(pastInventoryFileName)
	defer pastInventoryReader.Close()

	if e := pastInventoryReader.e; e != nil {
		fmt.Println(e)
	}

	nextName := "next"
	nextInventoryFileName := inventoryFileNameFor(nextName)
	nextInventoryWriter := newInventoryWriter(nextInventoryFileName)

	if e := nextInventoryWriter.e; e != nil {
		fmt.Println(e)
	}

	past := pastInventoryReader.readEntry()
	err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// ignore errors (TODO unless verbose)
			return nil
		}

		if info.IsDir() {
			// ignore directories (TODO unless extra verbose?)
			return nil
		}

		next := &fileEntry{path, info.Size(), info.ModTime().UTC()}
		cmp := past.comparePath(path)
		switch {
		case cmp == 0:
			sameSize := past.size == next.size
			sameMtime := past.mtime.Equal(next.mtime)
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

		nextInventoryWriter.writeEntry(next)

		return nil
	})
	if err != nil {
		fmt.Printf("error walking: %v\n", err)
		return
	}
}
