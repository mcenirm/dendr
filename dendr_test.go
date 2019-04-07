package main

import (
	"testing"
)

func TestCreateInspectorWhenOlddbIsNil(t *testing.T) {
	_, err := CreateInspector(nil, true, true)
	if err != nil {
		t.Errorf("CreateInpector had an error: %s", err)
	}
}

func TestCreateInspectorWhenNewdbIsNil(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestCreateInspectorWhenCollectorIsNil(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestInspectorWithNewFile(t *testing.T) {
	t.Errorf("test not implemented")
}

func TestInspectorWithUnchangedFile(t *testing.T) {
	t.Errorf("test not implemented")
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
