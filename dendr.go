package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type fileStats struct {
	size  int64
	mtime time.Time
}

func openFileDatabase(fileName string) (*fileDatabase, error) {
	fdb := &fileDatabase{}

	return fdb, nil
}

type fileDatabase struct {
}

func (fdb *fileDatabase) Close() error {
	return nil
}

func (fdb *fileDatabase) get(path string) (fileStats, bool, error) {
	return fileStats{}, false, nil
}

func (fdb *fileDatabase) set(path string, stats fileStats) error {
	return nil
}

func main() {
	fdb, err := openFileDatabase("test.db")
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getcwd: %v\n", err)
		return
	}

	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return nil
		}
		newStats := fileStats{info.Size(), info.ModTime()}
		oldStats, ok, err := fdb.get(path)
		if err != nil {
			return err
		}
		if ok {
			sameSize := oldStats.size == newStats.size
			sameMtime := oldStats.mtime.Equal(newStats.mtime)
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
		} else {
			fmt.Println("+sm  ", path)
		}
		err = fdb.set(path, newStats)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking: %v\n", err)
		return
	}
}
