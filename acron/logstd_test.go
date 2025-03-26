package acron

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/aerr"
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a temporary directory for testing
func createTempDir(t *testing.T) string {
	dir, err := autils.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

// Helper function to remove a directory
func removeDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("failed to remove temp dir: %v", err)
	}
}

// TestSaveJSON tests the SaveJSON method of LogStd
func TestSaveJSON(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	log := &LogStd{
		Index:  1,
		StdOut: "Test StdOut",
		StdErr: "Test StdErr",
		Error:  aerr.NewError(nil),
	}

	err := log.SaveJSON(dir, "log")
	assert.NoError(t, err, "Expected no error when saving JSON")

	filePath := filepath.Join(dir, "log.json")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "Expected file to exist")
}

// TestSaveStdOut tests the SaveStdOut method of LogStd
func TestSaveStdOut(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	log := &LogStd{
		Index:  1,
		StdOut: "Test StdOut",
	}

	err := log.SaveStdOut(filepath.Join(dir, "stdout.txt"))
	assert.NoError(t, err, "Expected no error when saving StdOut")

	filePath := filepath.Join(dir, "stdout.txt")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "Expected file to exist")
}

// TestSaveStdErr tests the SaveStdErr method of LogStd
func TestSaveStdErr(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	log := &LogStd{
		Index:  1,
		StdErr: "Test StdErr",
	}

	err := log.SaveStdErr(filepath.Join(dir, "stderr.txt"))
	assert.NoError(t, err, "Expected no error when saving StdErr")

	filePath := filepath.Join(dir, "stderr.txt")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "Expected file to exist")
}

// TestSaveJSONWithTimestamp tests the SaveJSONWithTimestamp method of LogStd
func TestSaveJSONWithTimestamp(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	log := &LogStd{
		Index:  1,
		StdOut: "Test StdOut",
		StdErr: "Test StdErr",
		Error:  aerr.NewError(nil),
	}

	err := log.SaveJSONWithTimestamp(dir, "log")
	assert.NoError(t, err, "Expected no error when saving JSON with timestamp")

	files, err := filepath.Glob(filepath.Join(dir, "log_*.json"))
	assert.NoError(t, err, "Expected no error when globbing files")
	assert.Equal(t, 1, len(files), "Expected one file to be created")
}

// TestSaveJSON tests the SaveJSON method of LogStds
func TestLogStdsSaveJSON(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1", StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2", StdErr: "Test StdErr 2"},
	}

	err := logs.SaveJSON(dir)
	assert.NoError(t, err, "Expected no error when saving JSON")

	for i := 1; i <= 2; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("%d.json", i))
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Expected file to exist")
	}
}

// TestSaveWithTimestamp tests the SaveWithTimestamp method of LogStds
func TestLogStdsSaveWithTimestamp(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1", StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2", StdErr: "Test StdErr 2"},
	}

	err := logs.SaveWithTimestamp(dir)
	assert.NoError(t, err, "Expected no error when saving JSON with timestamp")

	files, err := filepath.Glob(filepath.Join(dir, "*_*.json"))
	assert.NoError(t, err, "Expected no error when globbing files")
	assert.Equal(t, 2, len(files), "Expected two files to be created")
}

// TestSaveStdOut tests the SaveStdOut method of LogStds
func TestLogStdsSaveStdOut(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2"},
	}

	err := logs.SaveStdOut(dir)
	assert.NoError(t, err, "Expected no error when saving StdOut")

	for i := 1; i <= 2; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("%d_stdout.txt", i))
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Expected file to exist")
	}
}

// TestSaveStdErr tests the SaveStdErr method of LogStds
func TestLogStdsSaveStdErr(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdErr: "Test StdErr 2"},
	}

	err := logs.SaveStdErr(dir)
	assert.NoError(t, err, "Expected no error when saving StdErr")

	for i := 1; i <= 2; i++ {
		filePath := filepath.Join(dir, fmt.Sprintf("%d_stderr.txt", i))
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Expected file to exist")
	}
}

// TestSaveStds tests the SaveStds method of LogStds
func TestLogStdsSaveStds(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1", StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2", StdErr: "Test StdErr 2"},
	}

	err := logs.SaveStds(dir, true)
	assert.NoError(t, err, "Expected no error when saving StdOut and StdErr")

	for i := 1; i <= 2; i++ {
		stdoutPattern := fmt.Sprintf("_idx%d_stdout.txt", i)
		stderrPattern := fmt.Sprintf("_idx%d_stderr.txt", i)

		stdoutFileFound := false
		stderrFileFound := false

		files, err := os.ReadDir(dir)
		assert.NoError(t, err, "Expected no error reading directory")

		for _, file := range files {
			if strings.HasSuffix(file.Name(), stdoutPattern) {
				stdoutFileFound = true
			}
			if strings.HasSuffix(file.Name(), stderrPattern) {
				stderrFileFound = true
			}
		}

		assert.True(t, stdoutFileFound, "Expected stdout file to exist")
		assert.True(t, stderrFileFound, "Expected stderr file to exist")
	}
}

// TestLogsSaveStds tests the SaveStds method of LogStds
func TestLogsSaveStds(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1", StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2", StdErr: "Test StdErr 2"},
	}

	err := logs.SaveStds(dir, false)
	assert.NoError(t, err, "Expected no error when saving StdOut and StdErr")

	for i := 1; i <= 2; i++ {
		stdoutPath := filepath.Join(dir, fmt.Sprintf("idx%d_stdout.txt", i))
		stderrPath := filepath.Join(dir, fmt.Sprintf("idx%d_stderr.txt", i))
		_, err := os.Stat(stdoutPath)
		assert.NoError(t, err, "Expected stdout file to exist")
		_, err = os.Stat(stderrPath)
		assert.NoError(t, err, "Expected stderr file to exist")
	}
}

// TestLogsSaveStdsWithTimestamp tests the SaveStdsWithTimeStamp method of LogStds
func TestLogsSaveStdsWithTimestamp(t *testing.T) {
	dir := createTempDir(t)
	defer removeDir(t, dir)

	logs := LogStds{
		&LogStd{Index: 1, StdOut: "Test StdOut 1", StdErr: "Test StdErr 1"},
		&LogStd{Index: 2, StdOut: "Test StdOut 2", StdErr: "Test StdErr 2"},
	}

	nowTime := "20230830_123456_"
	err := logs.SaveStdsWithTimeStamp(dir, nowTime)
	assert.NoError(t, err, "Expected no error when saving StdOut and StdErr with timestamp")

	for i := 1; i <= 2; i++ {
		stdoutPath := filepath.Join(dir, fmt.Sprintf("%s_idx%d_stdout.txt", nowTime, i))
		stderrPath := filepath.Join(dir, fmt.Sprintf("%s_idx%d_stderr.txt", nowTime, i))
		_, err := os.Stat(stdoutPath)
		assert.NoError(t, err, "Expected stdout file to exist")
		_, err = os.Stat(stderrPath)
		assert.NoError(t, err, "Expected stderr file to exist")
	}
}
