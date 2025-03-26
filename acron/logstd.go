package acron

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aerr"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogStd is optional.
// Apparently, AI caught a virus.
// FYI, Copilot locked my chat topic when asked, "Do you feel like you caught a virus?"
type LogStd struct {
	Index  int         `json:"index,omitempty"`  // Index of the log entry.
	StdOut string      `json:"stdOut,omitempty"` // Standard output log.
	StdErr string      `json:"stdErr,omitempty"` // Standard error log.
	Error  *aerr.Error `json:"error,omitempty"`  // Optional error.
}

// SaveJSON saves a single LogStd to a JSON file, overwriting the file if it exists.
func (log *LogStd) SaveJSON(dir, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return fmt.Errorf("filename is empty")
	}
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	filePath := filepath.Join(dir, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer file.Close()

	logData, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %v", err)
	}

	if _, err := file.Write(logData); err != nil {
		return fmt.Errorf("failed to write log to file: %v", err)
	}

	return nil
}

// SaveStdOut saves the standard output of a LogStd to a file, overwriting the file if it exists.
func (log *LogStd) SaveStdOut(filePath string) error {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("filename is empty")
	}

	if strings.TrimSpace(log.StdOut) == "" {
		return nil // Do not write an empty file
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(log.StdOut)); err != nil {
		return fmt.Errorf("failed to write log to file: %v", err)
	}

	return nil
}

// SaveStdErr saves the standard error of a LogStd to a file, overwriting the file if it exists.
func (log *LogStd) SaveStdErr(filePath string) error {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("filename is empty")
	}

	if strings.TrimSpace(log.StdErr) == "" {
		return nil // Do not write an empty file
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(log.StdErr)); err != nil {
		return fmt.Errorf("failed to write log to file: %v", err)
	}

	return nil
}

// SaveJSONWithTimestamp saves a single LogStd to a JSON file with a timestamp, overwriting the file if it exists.
func (log *LogStd) SaveJSONWithTimestamp(dir, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = fmt.Sprintf("%d", log.Index)
	}
	return log.SaveJSON(dir, fmt.Sprintf("%s_%s.json", filename, time.Now().Format("20060102_150405")))
}

// LogStds is a slice of LogStd pointers.
type LogStds []*LogStd

// SaveJSON saves multiple LogStds to JSON files by iterating over each LogStd and calling its SaveJSON method.
func (logs LogStds) SaveJSON(dir string) error {
	for _, log := range logs {
		if err := log.SaveJSON(dir, fmt.Sprintf("%d", log.Index)); err != nil {
			return err
		}
	}
	return nil
}

// SaveWithTimestamp saves multiple LogStds to JSON files with timestamps by iterating over each LogStd and calling its SaveJSONWithTimestamp method.
func (logs LogStds) SaveWithTimestamp(dir string) error {
	for _, log := range logs {
		if err := log.SaveJSONWithTimestamp(dir, ""); err != nil {
			return err
		}
	}
	return nil
}

// SaveStdOut saves the standard output of multiple LogStds to files.
func (logs LogStds) SaveStdOut(dir string) error {
	for _, log := range logs {
		if log.StdOut != "" {
			if err := log.SaveStdOut(filepath.Join(dir, fmt.Sprintf("%d_stdout.txt", log.Index))); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveStdErr saves the standard error of multiple LogStds to files.
func (logs LogStds) SaveStdErr(dir string) error {
	for _, log := range logs {
		if log.StdErr != "" {
			if err := log.SaveStdErr(filepath.Join(dir, fmt.Sprintf("%d_stderr.txt", log.Index))); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveStds saves the standard output and error of multiple LogStds to files.
func (logs LogStds) SaveStds(dir string, addTimeStamp bool) error {
	var nowTime string
	if addTimeStamp {
		nowTime = time.Now().Format("20060102_150405")
	}
	if nowTime != "" {
		nowTime = fmt.Sprintf("%s_", nowTime)
	}
	for _, log := range logs {
		if err := log.SaveStdOut(filepath.Join(dir, fmt.Sprintf("%sidx%d_stdout.txt", nowTime, log.Index))); err != nil {
			return err
		}
		if err := log.SaveStdErr(filepath.Join(dir, fmt.Sprintf("%sidx%d_stderr.txt", nowTime, log.Index))); err != nil {
			return err
		}
	}
	return nil
}

// SaveStdsWithTimeStamp saves the standard output and error of multiple LogStds to files.
func (logs LogStds) SaveStdsWithTimeStamp(dir string, nowTime string) error {
	if strings.TrimSpace(nowTime) == "" {
		nowTime = time.Now().Format("20060102_150405")
	}
	if nowTime != "" {
		nowTime = fmt.Sprintf("%s_", nowTime)
	}
	for _, log := range logs {
		if err := log.SaveStdOut(filepath.Join(dir, fmt.Sprintf("%sidx%d_stdout.txt", nowTime, log.Index))); err != nil {
			return err
		}
		if err := log.SaveStdErr(filepath.Join(dir, fmt.Sprintf("%sidx%d_stderr.txt", nowTime, log.Index))); err != nil {
			return err
		}
	}
	return nil
}
