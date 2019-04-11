package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	placeholder = "placeholder"
	testpath    = "/test/placeholder"
)

type fakeFileInfo struct {
	size    int64     // length in bytes for regular files; system-dependent for others
	modTime time.Time // modification time
}

func (info fakeFileInfo) Name() string       { return placeholder }
func (info fakeFileInfo) Size() int64        { return info.size }
func (info fakeFileInfo) Mode() os.FileMode  { return 0644 }
func (info fakeFileInfo) ModTime() time.Time { return info.modTime }
func (info fakeFileInfo) IsDir() bool        { return false }
func (info fakeFileInfo) Sys() interface{}   { return nil }

type fakeDB struct {
	before fakeFileInfo
}

type fakeCollector struct {
	change       change
	changedStats changedStats
	path_        string
	info         os.FileInfo
}

func (c fakeCollector) collect(theChange change, theChangedStats changedStats, path string, info os.FileInfo) error {
	fmt.Println("        c pre:  ", c)
	c.change = theChange
	c.changedStats = theChangedStats
	c.path_ = path
	c.info = info
	fmt.Println("        c post: ", c)
	return nil
}

func makeFakes(t *testing.T) (fakeDB, fakeCollector, filepath.WalkFunc) {
	theDB := fakeDB{}
	theCollector := fakeCollector{}
	theInspector, err := createInspector(theDB, theCollector)
	if err != nil {
		t.Errorf("unexpected error from createInspector: %v", err)
	}
	return theDB, theCollector, theInspector
}

func callInspector(t *testing.T, theInspector filepath.WalkFunc, f fakeFileInfo) {
	err := theInspector(testpath, f, nil)
	if err != nil {
		t.Errorf("unexpected error from inspector: %v", err)
	}
}

func checkCollector(t *testing.T, theCollector fakeCollector, theExpectedChange change) {
	if theCollector.change != theExpectedChange {
		t.Errorf("expected %v but got %v", theExpectedChange, theCollector.change)
	}
}

func testInspector(t *testing.T, theChange change, theChangedStats changedStats) {
	theDB, theCollector, theInspector := makeFakes(t)

	after := fakeFileInfo{123, time.Now()}
	if theChange != added {
		theDB.before = after
		if theChange != unchanged {
			if theChangedStats&changedModTime != 0 {
				earlier := theDB.before.modTime.Add(time.Duration(-7000000000))
				theDB.before.modTime = earlier
			}
			if theChangedStats&changedSize != 0 {
				larger := theDB.before.size + 100
				theDB.before.size = larger
			}
		}
	}

	callInspector(t, theInspector, after)
	checkCollector(t, theCollector, theChange)
}

func TestCreateInspectorWhenDBIsNil(t *testing.T) {
	theCollector := fakeCollector{}
	_, err := createInspector(nil, theCollector)
	if err == nil {
		t.Errorf("expected error from CreateInpector")
	}
}

func TestCreateInspectorWhenCollectorIsNil(t *testing.T) {
	theDB := fakeDB{}
	_, err := createInspector(theDB, nil)
	if err == nil {
		t.Errorf("expected error from CreateInpector: %v", err)
	}
}

func TestInspectorWithUnchangedFile(t *testing.T) {
	testInspector(t, unchanged, 0)
}

func TestInspectorWithAddedFile(t *testing.T) {
	testInspector(t, added, 0)
}

func TestInspectorWithFileWithDifferentModificationTime(t *testing.T) {
	testInspector(t, statsChanged, changedModTime)
}

func TestInspectorWithFileWithDifferentSize(t *testing.T) {
	testInspector(t, statsChanged, changedSize)
}

func TestInspectorWithFileWithDifferentModificationTimeAndSize(t *testing.T) {
	testInspector(t, statsChanged, changedModTime|changedSize)
}

func TestInspectorWithDirectory(t *testing.T) {
	t.Errorf("test not implemented")
}
