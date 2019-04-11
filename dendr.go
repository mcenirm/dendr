package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type fileDatabase interface{}

//go:generate stringer -type=change,changedStats

type change int

const (
	_ change = iota
	unchanged
	added
	statsChanged
)

type changedStats int

const (
	changedModTime changedStats = 1 << iota
	changedSize
)

type collector interface {
	collect(theChange change, theChangedStats changedStats, path string, info os.FileInfo) error
}

func createInspector(theDB fileDatabase, theCollector collector) (inspector filepath.WalkFunc, err error) {
	if theDB == nil {
		return nil, errors.New("theDB must not be nil")
	}
	if theCollector == nil {
		return nil, errors.New("theCollector must not be nil")
	}
	return func(path string, info os.FileInfo, err error) error {
		fmt.Println("theCollector pre:  ", theCollector)
		theCollector.collect(unchanged, 0, path, info)
		fmt.Println("theCollector post: ", theCollector)
		return nil
	}, nil
}

func main() {
	fmt.Println("TODO:")
	fmt.Println(" * Open fileDatabase")
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
