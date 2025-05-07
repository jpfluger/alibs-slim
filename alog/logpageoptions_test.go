package alog

import (
	"encoding/csv"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadLogFilePaged(t *testing.T) {
	logPath := filepath.Join("test_data", "tests.log")

	tests := []struct {
		name     string
		opts     LogPageOptions
		expected []string
		total    int
	}{
		{
			name: "basic no filter",
			opts: LogPageOptions{
				FilePath: logPath,
				Count:    3,
			},
			expected: []string{
				"[INFO] Server started",
				"[DEBUG] Loading configuration",
				"[INFO] Connected to database",
			},
			total: 10,
		},
		{
			name: "case-insensitive substring filter (error)",
			opts: LogPageOptions{
				FilePath:        logPath,
				Count:           2,
				Filter:          "error",
				IsCaseSensitive: false,
				NoRegex:         true,
			},
			expected: []string{
				"[ERROR] Failed to connect to service A",
				"[ERROR] Timeout on service B",
			},
			total: 2,
		},
		{
			name: "case-sensitive substring filter (ERROR)",
			opts: LogPageOptions{
				FilePath:        logPath,
				Count:           1,
				Filter:          "ERROR",
				IsCaseSensitive: true,
				NoRegex:         true,
			},
			expected: []string{
				"[ERROR] Failed to connect to service A",
			},
			total: 2,
		},
		{
			name: "regex match lines with INFO at start",
			opts: LogPageOptions{
				FilePath: logPath,
				Count:    10,
				Filter:   `^\[INFO\]`,
			},
			expected: []string{
				"[INFO] Server started",
				"[INFO] Connected to database",
				"[INFO] Request received from 127.0.0.1",
				"[INFO] Job completed successfully",
			},
			total: 4,
		},
		{
			name: "tail last 2 warning lines",
			opts: LogPageOptions{
				FilePath: logPath,
				Count:    2,
				Filter:   "WARN",
				Tail:     true,
				NoRegex:  true,
			},
			expected: []string{
				"[WARN] Disk space low",
				"[WARN] Memory usage high",
			},
			total: 2,
		},
		{
			name: "offset and count",
			opts: LogPageOptions{
				FilePath: logPath,
				Filter:   "INFO",
				Count:    2,
				Offset:   1,
			},
			expected: []string{
				"[INFO] Connected to database",
				"[INFO] Request received from 127.0.0.1",
			},
			total: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadLogFilePaged(tt.opts)
			if err := result.Err(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lines := result.Lines()

			if result.TotalCount() != tt.total {
				t.Errorf("expected total %d, got %d", tt.total, result.TotalCount())
			}

			if len(lines) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(lines))
			}

			for i := range lines {
				if lines[i] != tt.expected[i] {
					t.Errorf("line %d: expected %q, got %q", i, tt.expected[i], lines[i])
				}
			}
		})
	}
}

func TestReadLogFilePaged_JSONLog(t *testing.T) {
	logPath := filepath.Join("test_data", "tests.json")

	tests := []struct {
		name         string
		opts         LogPageOptions
		expected     []string
		total        int
		expectIsJSON bool
	}{
		{
			name: "basic read json file",
			opts: LogPageOptions{
				FilePath: logPath,
				Count:    2,
			},
			expected: []string{
				`{"level":"info","msg":"Service started","ts":"2024-04-01T10:00:00Z"}`,
				`{"level":"debug","msg":"Connecting to database","ts":"2024-04-01T10:00:01Z"}`,
			},
			total:        7,
			expectIsJSON: true,
		},
		{
			name: "filter error logs with substring",
			opts: LogPageOptions{
				FilePath:        logPath,
				Count:           5,
				Filter:          "error",
				NoRegex:         true,
				IsCaseSensitive: false,
			},
			expected: []string{
				`{"level":"error","msg":"Database connection failed","ts":"2024-04-01T10:00:04Z"}`,
				`{"level":"error","msg":"Service timeout","ts":"2024-04-01T10:00:06Z"}`,
			},
			total:        2,
			expectIsJSON: true,
		},
		{
			name: "regex filter messages ending in 'timeout'",
			opts: LogPageOptions{
				FilePath: logPath,
				Count:    5,
				Filter:   `"msg":"[^"]*timeout"`,
			},
			expected: []string{
				`{"level":"error","msg":"Service timeout","ts":"2024-04-01T10:00:06Z"}`,
			},
			total:        1,
			expectIsJSON: true,
		},
		{
			name: "substring filter for 'timeout' using NoRegex",
			opts: LogPageOptions{
				FilePath:        logPath,
				Count:           5,
				Filter:          "timeout",
				NoRegex:         true,
				IsCaseSensitive: false,
			},
			expected: []string{
				`{"level":"error","msg":"Service timeout","ts":"2024-04-01T10:00:06Z"}`,
			},
			total:        1,
			expectIsJSON: true,
		},
		{
			name: "tail last 2 info logs",
			opts: LogPageOptions{
				FilePath: logPath,
				Filter:   `"level":"info"`,
				Count:    2,
				Tail:     true,
			},
			expected: []string{
				`{"level":"info","msg":"Request received","ts":"2024-04-01T10:00:02Z"}`,
				`{"level":"info","msg":"Job completed","ts":"2024-04-01T10:00:05Z"}`,
			},
			total:        3,
			expectIsJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadLogFilePaged(tt.opts)
			if err := result.Err(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lines := result.Lines()

			if result.TotalCount() != tt.total {
				t.Errorf("expected total %d, got %d", tt.total, result.TotalCount())
			}

			if len(lines) != len(tt.expected) {
				t.Fatalf("expected %d lines, got %d", len(tt.expected), len(lines))
			}

			for i := range lines {
				if lines[i] != tt.expected[i] {
					t.Errorf("line %d: expected %q, got %q", i, tt.expected[i], lines[i])
				}
			}

			if result.IsFirstLineJSON() != tt.expectIsJSON {
				t.Errorf("expected isJSON = %v, got %v", tt.expectIsJSON, result.IsFirstLineJSON())
			}
		})
	}
}

func TestLogPageResult_AnalyzeJSON(t *testing.T) {
	logPath := filepath.Join("test_data", "tests.json")

	result := ReadLogFilePaged(LogPageOptions{
		FilePath: logPath,
		Count:    10,
	})
	if err := result.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys, rows := result.AnalyzeJSON()

	if len(keys) == 0 {
		t.Errorf("expected at least one key in JSON logs")
	}
	if len(rows) != result.TotalCount() {
		t.Errorf("expected %d rows, got %d", result.TotalCount(), len(rows))
	}

	expectedKeys := []string{"level", "msg", "ts"}
	for _, k := range expectedKeys {
		found := false
		for _, actual := range keys {
			if k == actual {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing expected key %q in analyzed keys", k)
		}
	}
}

func TestLogPageResult_ToJSON(t *testing.T) {
	logPath := filepath.Join("test_data", "tests.json")

	result := ReadLogFilePaged(LogPageOptions{
		FilePath: logPath,
		Count:    2,
	})
	if err := result.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonBytes, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed []string
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Errorf("ToJSON output is not valid JSON array: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(parsed))
	}
}

func TestLogPageResult_ToCSV(t *testing.T) {
	logPath := filepath.Join("test_data", "tests.json")

	result := ReadLogFilePaged(LogPageOptions{
		FilePath: logPath,
		Count:    3,
	})
	if err := result.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	csvOutput, err := result.ToCSV()
	if err != nil {
		t.Fatalf("ToCSV failed: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(csvOutput))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("CSV parsing failed: %v", err)
	}

	if len(records) < 2 {
		t.Errorf("expected at least 2 rows in CSV (1 header + 1 data), got %d", len(records))
	}

	headers := records[0]
	if len(headers) < 3 {
		t.Errorf("expected at least 3 columns in CSV header, got %d", len(headers))
	}
}
