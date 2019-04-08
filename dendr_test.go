package main

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

type FakeFileInfo struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    os.FileMode // file mode bits
	modTime time.Time   // modification time
}

func (info FakeFileInfo) Name() string       { return info.name }
func (info FakeFileInfo) Size() int64        { return info.size }
func (info FakeFileInfo) Mode() os.FileMode  { return info.mode }
func (info FakeFileInfo) ModTime() time.Time { return info.modTime }
func (info FakeFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info FakeFileInfo) Sys() interface{}   { return nil }

type FakeFile struct {
	path string
	info os.FileInfo
}

type FakeDB struct {
	sample FakeFile
}

func (db FakeDB) AddSampleFile() FakeFile {
	info := FakeFileInfo{"sample", 123, 0644, time.Now()}
	path_ := path.Join("test", info.Name())
	db.sample = FakeFile{path_, info}
	return db.sample
}

type FakeCollector struct {
	change Change
	path_  string
	info   os.FileInfo
}

func (c FakeCollector) Collect(change Change, path string, info os.FileInfo) error {
	c.change = change
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

func callInspector(t *testing.T, inspector filepath.WalkFunc, f FakeFile) {
	err := inspector(f.path, f.info, nil)
	if err != nil {
		t.Errorf("unexpected error from inspector: %v", err)
	}
}

func checkCollector(t *testing.T, collector FakeCollector, expected Change) {
	if collector.change != expected {
		t.Errorf("expected %v but got %v", expected, collector.change)
	}
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
	db, collector, inspector := makeFakes(t)
	f := db.AddSampleFile()
	callInspector(t, inspector, f)
	checkCollector(t, collector, Unchanged)
}

func TestInspectorWithAddedFile(t *testing.T) {
	_, collector, inspector := makeFakes(t)
	db2 := FakeDB{}
	f := db2.AddSampleFile()
	callInspector(t, inspector, f)
	checkCollector(t, collector, Added)
}

func TestInspectorWithFileWithDifferentModificationTime(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestInspectorWithFileWithDifferentSize(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestInspectorWithFileWithDifferentModificationTimeAndSize(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestInspectorWithDirectory(t *testing.T) {
	t.Errorf("test not implemented")
}
