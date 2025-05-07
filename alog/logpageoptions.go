package alog

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

// LogPageOptions defines the filtering and paging options for reading log files.
type LogPageOptions struct {
	FilePath        string `json:"filePath"`        // Full path to the log file
	Count           int    `json:"count"`           // Number of lines to return
	Offset          int    `json:"offset"`          // Number of matching lines to skip (ignored if Tail is true)
	Filter          string `json:"filter"`          // Include only lines containing this substring
	Tail            bool   `json:"tail"`            // If true, return the last N matching lines
	IsCaseSensitive bool   `json:"isCaseSensitive"` // If false, filter will match case-insensitively
	NoRegex         bool   `json:"noRegex"`         // If true, use plain substring matching instead of regex
}

// sanitizeRegexPattern escapes regex characters if NoRegex is true.
func sanitizeRegexPattern(input string, noRegex bool) string {
	if noRegex {
		// Literal match: escape all regex metacharacters
		return regexp.QuoteMeta(input)
	}
	return input
}

func isJSONLine(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}

	var js interface{}
	return json.Unmarshal([]byte(line), &js) == nil
}

// ReadLogFilePaged reads a log file with filtering, offset, and tailing.
func ReadLogFilePaged(opts LogPageOptions) *LogPageResult {
	result := &LogPageResult{opts: opts}

	if opts.Count <= 0 {
		opts.Count = 100
	}

	file, err := os.Open(opts.FilePath)
	if err != nil {
		result.parseErr = err
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var filterPattern *regexp.Regexp
	filter := opts.Filter

	if filter != "" && !opts.NoRegex {
		pattern := filter
		if !opts.IsCaseSensitive {
			pattern = `(?i)` + pattern
		}
		filterPattern, err = regexp.Compile(pattern)
		if err != nil {
			result.parseErr = fmt.Errorf("invalid regex pattern: %w", err)
			return result
		}
	}

	matchesFilter := func(line string) bool {
		if filter == "" {
			return true
		}
		if opts.NoRegex {
			if !opts.IsCaseSensitive {
				return strings.Contains(strings.ToLower(line), strings.ToLower(filter))
			}
			return strings.Contains(line, filter)
		}
		return filterPattern.MatchString(line)
	}

	firstMatchChecked := false
	if opts.Tail {
		buf := make([]string, opts.Count)
		index := 0
		size := 0

		for scanner.Scan() {
			line := scanner.Text()
			if !matchesFilter(line) {
				continue
			}

			if !firstMatchChecked {
				result.isFirstJSON = isJSONLine(line)
				firstMatchChecked = true
			}

			buf[index] = line
			index = (index + 1) % opts.Count
			if size < opts.Count {
				size++
			}
			result.totalCount++
		}

		result.lines = make([]string, size)
		for i := 0; i < size; i++ {
			result.lines[i] = buf[(index+i)%opts.Count]
		}

	} else {
		skipped := 0
		result.lines = make([]string, 0, opts.Count)

		for scanner.Scan() {
			line := scanner.Text()
			if !matchesFilter(line) {
				continue
			}

			if !firstMatchChecked {
				result.isFirstJSON = isJSONLine(line)
				firstMatchChecked = true
			}

			result.totalCount++
			if skipped < opts.Offset {
				skipped++
				continue
			}
			if len(result.lines) < opts.Count {
				result.lines = append(result.lines, line)
			}
		}
	}

	result.parseErr = scanner.Err()
	return result
}

type LogPageResult struct {
	lines        []string
	totalCount   int
	isFirstJSON  bool
	opts         LogPageOptions
	parseErr     error
	analyzed     bool
	analyzedKeys []string
	analyzedRows []map[string]interface{}
}

func (r *LogPageResult) Lines() []string {
	return r.lines
}

func (r *LogPageResult) TotalCount() int {
	return r.totalCount
}

func (r *LogPageResult) IsFirstLineJSON() bool {
	return r.isFirstJSON
}

func (r *LogPageResult) Err() error {
	return r.parseErr
}

func (r *LogPageResult) Options() LogPageOptions {
	return r.opts
}

func (r *LogPageResult) JSONKeys() []string {
	keys, _ := r.AnalyzeJSON()
	return keys
}

func (r *LogPageResult) JSONRows() []map[string]interface{} {
	_, rows := r.AnalyzeJSON()
	return rows
}

func (r *LogPageResult) AnalyzeJSON() ([]string, []map[string]interface{}) {
	if r.analyzed {
		return r.analyzedKeys, r.analyzedRows
	}
	r.analyzedKeys, r.analyzedRows = AnalyzeJSONLines(r.lines)
	r.analyzed = true
	return r.analyzedKeys, r.analyzedRows
}

func (r *LogPageResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r.lines, "", "  ")
}

func (r *LogPageResult) ToCSV() (string, error) {
	keys, rows := r.AnalyzeJSON()
	if len(keys) == 0 || len(rows) == 0 {
		return "", fmt.Errorf("no valid JSON rows for CSV")
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	_ = writer.Write(keys)

	for _, row := range rows {
		record := make([]string, len(keys))
		for i, key := range keys {
			if val, ok := row[key]; ok {
				record[i] = fmt.Sprint(val)
			}
		}
		_ = writer.Write(record)
	}
	writer.Flush()
	return sb.String(), nil
}

func AnalyzeJSONLines(logLines []string) ([]string, []map[string]interface{}) {
	uniqueKeys := map[string]bool{}
	parsedRows := make([]map[string]interface{}, 0, len(logLines))

	for _, line := range logLines {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			continue // skip invalid JSON lines
		}
		parsedRows = append(parsedRows, parsed)
		for k := range parsed {
			uniqueKeys[k] = true
		}
	}

	// Sort keys alphabetically (optional, but stable)
	keys := make([]string, 0, len(uniqueKeys))
	for k := range uniqueKeys {
		keys = append(keys, k)
	}
	slices.Sort(keys) // Go 1.21+ or use sort.Strings(keys)

	return keys, parsedRows
}
