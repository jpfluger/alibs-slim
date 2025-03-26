package autils

import (
	"encoding/csv"
	"strings"
	"testing"
)

// mockRead simulates reading from a CSV file by providing predefined data.
func mockRead(data string) *csv.Reader {
	return csv.NewReader(strings.NewReader(data))
}

// TestReadCSVSuccess tests the ReadCSV function with a successful CSV read.
func TestReadCSVSuccess(t *testing.T) {
	// Simulate CSV data with two header lines and two data lines.
	csvData := "header1,header2\nheaderA,headerB\ndata1,data2\ndataA,dataB\n"
	csvReader := mockRead(csvData)

	// Define a callback function that simply counts the number of data lines processed.
	lineCount := 0
	fn := func(headers map[int][]string, lineOn int, line []string) (doStop bool, err error) {
		lineCount++
		return false, nil
	}

	// Call ReadCSV with the simulated CSV data and the callback function.
	err := ReadCSV(csvReader, 2, fn)
	if err != nil {
		t.Errorf("ReadCSV failed: %v", err)
	}
	if lineCount != 2 {
		t.Errorf("ReadCSV processed %d lines; want 2", lineCount)
	}
}

// TestReadCSVHeaderCount tests the ReadCSV function with varying header line counts.
func TestReadCSVHeaderCount(t *testing.T) {
	csvData := "header1,header2\nheaderA,headerB\ndata1,data2\n"

	// Test with different header line counts.
	for _, headerCount := range []int{1, 2} {
		csvReader := mockRead(csvData)
		headersProcessed := make(map[int][]string)
		fn := func(headers map[int][]string, lineOn int, line []string) (doStop bool, err error) {
			for k, v := range headers {
				headersProcessed[k] = v
			}
			return false, nil
		}

		err := ReadCSV(csvReader, headerCount, fn)
		if err != nil {
			t.Errorf("ReadCSV with headerCount %d failed: %v", headerCount, err)
		}
		if len(headersProcessed) != headerCount {
			t.Errorf("ReadCSV with headerCount %d processed %d headers; want %d", headerCount, len(headersProcessed), headerCount)
		}
	}
}

// TestReadCSVStopEarly tests the ReadCSV function with an early stop.
func TestReadCSVStopEarly(t *testing.T) {
	csvData := "header1,header2\ndata1,data2\ndataA,dataB\n"
	csvReader := mockRead(csvData)

	fn := func(headers map[int][]string, lineOn int, line []string) (doStop bool, err error) {
		return true, nil // Stop after the first data line.
	}

	lineCount := 0
	err := ReadCSV(csvReader, 1, func(headers map[int][]string, lineOn int, line []string) (doStop bool, err error) {
		lineCount++
		return fn(headers, lineOn, line)
	})
	if err != nil {
		t.Errorf("ReadCSV failed: %v", err)
	}
	if lineCount != 1 {
		t.Errorf("ReadCSV should have stopped early but processed %d lines", lineCount)
	}
}
