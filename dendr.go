package main

import (
	"errors"
	"fmt"
	"path/filepath"
)

type DB interface{}
type Collector interface{}

func CreateInspector(olddb DB, newdb DB, collector Collector) (inspector filepath.WalkFunc, err error) {
	return nil, errors.New("not implemented")
}

func main() {
	fmt.Println("TODO:")
	fmt.Println(" * Open old database")
	fmt.Println(" * Open new database")
	fmt.Println(" * Scan tree")
}
