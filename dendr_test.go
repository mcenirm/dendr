package main

import (
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

type FakeDB struct {
	before fakeFileInfo
}

type FakeCollector struct {
	change       Change
	changedStats ChangedStats
	path_        string
	info         os.FileInfo
}

func (c FakeCollector) Collect(change Change, changedStats ChangedStats, path string, info os.FileInfo) error {
	c.change = change
	c.changedStats = changedStats
	c.path_ = path
	c.info = info
	return nil
}

func makeFakes(t *testing.T) (FakeDB, FakeCollector, filepath.WalkFunc) {
	db := FakeDB{}
	collector := FakeCollector{}
	inspector, err := CreateInspector(db, collector)
	if err != nil {
		t.Errorf("unexpected error from CreateInspector: %v", err)
	}
	return db, collector, inspector
}

func callInspector(t *testing.T, inspector filepath.WalkFunc, f fakeFileInfo) {
	err := inspector(testpath, f, nil)
	if err != nil {
		t.Errorf("unexpected error from inspector: %v", err)
	}
}

func checkCollector(t *testing.T, collector FakeCollector, expected Change) {
	if collector.change != expected {
		t.Errorf("expected %v but got %v", expected, collector.change)
	}
}

func testInspector(t *testing.T, change Change, changedStats ChangedStats) {
	db, collector, inspector := makeFakes(t)

	after := fakeFileInfo{123, time.Now()}
	if change != Added {
		db.before = after
		if change != Unchanged {
			if changedStats&ChangedModTime != 0 {
				earlier := db.before.modTime.Add(time.Duration(-7000000000))
				db.before.modTime = earlier
			}
			if changedStats&ChangedSize != 0 {
				larger := db.before.size + 100
				db.before.size = larger
			}
		}
	}

	callInspector(t, inspector, after)
	checkCollector(t, collector, change)
}

func TestCreateInspectorWhenDBIsNil(t *testing.T) {
	collector := FakeCollector{}
	_, err := CreateInspector(nil, collector)
	if err == nil {
		t.Errorf("expected error from CreateInpector")
	}
}

func TestCreateInspectorWhenCollectorIsNil(t *testing.T) {
	db := FakeDB{}
	_, err := CreateInspector(db, nil)
	if err == nil {
		t.Errorf("expected error from CreateInpector: %v", err)
	}
}

func TestInspectorWithUnchangedFile(t *testing.T) {
	testInspector(t, Unchanged, 0)
}

func TestInspectorWithAddedFile(t *testing.T) {
	testInspector(t, Added, 0)
}

func TestInspectorWithFileWithDifferentModificationTime(t *testing.T) {
	testInspector(t, StatsChanged, ChangedModTime)
}

func TestInspectorWithFileWithDifferentSize(t *testing.T) {
	testInspector(t, StatsChanged, ChangedSize)
}

func TestInspectorWithFileWithDifferentModificationTimeAndSize(t *testing.T) {
	testInspector(t, StatsChanged, ChangedModTime|ChangedSize)
}

func TestInspectorWithDirectory(t *testing.T) {
	t.Errorf("test not implemented")
}
