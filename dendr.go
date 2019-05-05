package main

import (
	"bufio"
	"flag"
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

func openPastInventoryReader(pastName string) *inventoryReader {
	pastInventoryFileName := inventoryFileNameFor(pastName)
	return newInventoryReader(pastInventoryFileName)
}

func openNextInventoryWriter(nextName string) *inventoryWriter {
	nextInventoryFileName := inventoryFileNameFor(nextName)
	return newInventoryWriter(nextInventoryFileName)
}

func realmain(start string, pastName string, nextName string, quiet bool, verbose bool) {
	var err error

	pastInventoryReader := openPastInventoryReader(pastName)
	defer pastInventoryReader.Close()

	if e := pastInventoryReader.e; e != nil {
		if !quiet {
			fmt.Println(e)
		}
	}

	nextInventoryWriter := openNextInventoryWriter(nextName)
	defer nextInventoryWriter.Close()

	if e := nextInventoryWriter.e; e != nil {
		if !quiet {
			fmt.Println(e)
		}
	}

	past := pastInventoryReader.readEntry()
	err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// ignore errors (TODO unless verbose)
			return nil
		}

		if !info.Mode().IsRegular() {
			// ignore directories, symlinks, etc (TODO unless extra verbose?)
			return nil
		}

		next := &fileEntry{path, info.Size(), info.ModTime().UTC()}
		if past == nil {
			if !quiet {
				fmt.Println("+++  ", path)
			}
		} else {
		pastloop:
			for keepgoing := true; keepgoing && past != nil; past = pastInventoryReader.readEntry() {
				cmp := past.comparePath(path)
				switch {
				case cmp < 0:
					if !quiet {
						fmt.Println("---  ", past.path)
					}
				case cmp == 0:
					sameSize := past.size == next.size
					sameMtime := past.mtime.Equal(next.mtime)
					if sameSize && sameMtime {
						// do nothing?
					} else {
						if !quiet {
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
					}
					keepgoing = false
				default:
					if !quiet {
						fmt.Println("+++  ", path)
					}
					break pastloop
				}
			}
		}

		nextInventoryWriter.writeEntry(next)

		return nil
	})
	if err != nil {
		if !quiet {
			fmt.Printf("error walking: %v\n", err)
		}
		return
	}
	for ; past != nil; past = pastInventoryReader.readEntry() {
		if !quiet {
			fmt.Println("---  ", past.path)
		}
	}
}

func main() {
	var (
		flagpath     string
		flagpastname string
		flagnextname string
		flagquiet    bool
		flagverbose  bool
	)
	flag.StringVar(&flagpath, "path", "testpath", "path to scan")
	flag.StringVar(&flagpastname, "pastname", "past", "name of past inventory")
	flag.StringVar(&flagnextname, "nextname", "next", "name of next inventory")
	flag.BoolVar(&flagquiet, "quiet", false, "suppress output")
	flag.BoolVar(&flagverbose, "verbose", false, "show more details")
	flag.Parse()
	realmain(flagpath, flagpastname, flagnextname, flagquiet, flagverbose)
}
