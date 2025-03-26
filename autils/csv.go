package autils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// FNReadCSVLine is a function type that processes a line from a CSV file.
// It takes a map of headers, the current line number, and the line content.
// It returns a boolean indicating whether to stop reading the CSV and any error encountered.
type FNReadCSVLine func(headers map[int][]string, lineOn int, line []string) (doStop bool, err error)

// ReadCSVFile opens a CSV file and uses ReadCSV to process the file line by line.
func ReadCSVFile(fileCSV string, headerLineCount int, fn FNReadCSVLine) error {
	// Open the CSV file for reading.
	fcsv, err := os.Open(fileCSV)
	if err != nil {
		return err
	}
	defer fcsv.Close()

	// Create a new CSV reader with a buffered reader for efficiency.
	csvReader := csv.NewReader(bufio.NewReader(fcsv))

	// Use ReadCSV to process the CSV data.
	return ReadCSV(csvReader, headerLineCount, fn)
}

// ReadCSV processes CSV data from a *csv.Reader, handling headers and invoking a callback function for each line.
func ReadCSV(csvReader *csv.Reader, headerLineCount int, fn FNReadCSVLine) error {
	// Initialize line counter and a map to store headers.
	lineOn := 0
	headers := make(map[int][]string)

	// Read the CSV data line by line.
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached, stop reading.
		} else if err != nil {
			return fmt.Errorf("read error: %v", err)
		}

		lineOn++ // Increment line counter.

		// Check if the current line is a header line.
		if lineOn <= headerLineCount {
			headers[lineOn-1] = line
			continue // Skip processing for header lines.
		}

		// Process the current line using the provided function.
		doStop, err := fn(headers, lineOn, line)
		if doStop {
			return nil // Stop reading as indicated by the processing function.
		} else if err != nil {
			return fmt.Errorf("error processing CSV line: %v", err)
		}
	}

	return nil // Successfully processed the CSV data without errors.
}
