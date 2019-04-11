package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type DB interface{}

//go:generate stringer -type=Change,ChangedStats

type Change int

const (
	_ Change = iota
	Unchanged
	Added
	StatsChanged
)

type ChangedStats int

const (
	ChangedModTime ChangedStats = 1 << iota
	ChangedSize
)

type Collector interface {
	Collect(change Change, changedStats ChangedStats, path string, info os.FileInfo) error
}

func CreateInspector(db DB, collector Collector) (inspector filepath.WalkFunc, err error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	if collector == nil {
		return nil, errors.New("collector must not be nil")
	}
	return func(path string, info os.FileInfo, err error) error {
		fmt.Println("collector pre:  ", collector)
		collector.Collect(Unchanged, 0, path, info)
		fmt.Println("collector post: ", collector)
		return nil
	}, nil
}

func main() {
	fmt.Println("TODO:")
	fmt.Println(" * Open database")
	fmt.Println(" * Scan tree and report changes")
	fmt.Println()

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
		if info.IsDir() && info.Name() == ".git" {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n  %8v %v %v %q\n", path, info.Size(), info.Mode(), info.ModTime(), info.Name())
		return nil
	})
	if err != nil {
		fmt.Printf("error walking: %v\n", err)
		return
	}
}
