package atemplates

import (
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTempDir creates a temporary directory for testing and returns its path.
// If it fails to create the directory, the test is aborted with a fatal error.
func createTempDir(t *testing.T) string {
	dir, err := autils.CreateTempDir() // Attempt to create a temporary directory.
	if err != nil {
		t.Fatalf("cannot create temp directory for logging; %v", err) // Fatal error if creation fails.
	}
	return dir // Return the path to the newly created temporary directory.
}

// deleteDir removes a list of directories provided in the dir slice.
// If it fails to remove a directory, the test is aborted with a fatal error.
func deleteDir(t *testing.T, dir []string) {
	for _, d := range dir {
		if err := os.RemoveAll(d); err != nil {
			t.Fatalf("failed to remove test directory at %s; %v", d, err) // Fatal error if removal fails.
		}
	}
}

// TestParseDirective tests the ParseDirective function with various scenarios.
func TestParseDirective(t *testing.T) {
	dir1 := createTempDir(t)           // Create a temporary directory for the test.
	defer deleteDir(t, []string{dir1}) // Schedule the deletion of the temporary directory.

	data := "hello world"                                           // Test data.
	file1 := path.Join(dir1, "test.txt")                            // Create a file path in the temporary directory.
	if err := os.WriteFile(file1, []byte(data), 0777); err != nil { // Write test data to the file.
		t.Error(err) // Non-fatal error if writing fails.
		return
	}

	// Assert various conditions to test the ParseDirective and HasDirective functions.
	assert.Equal(t, false, HasDirective(data), "Data should not have a directive")
	assert.Equal(t, data, ParseDirective(data), "ParseDirective should return original data when no directive is present")
	assert.Equal(t, false, HasDirective(data), "Data should not have a directive after parsing")
	assert.Equal(t, file1, ParseDirective(file1), "ParseDirective should return file path when no directive is present")

	dataDirective := "{{file:" + file1 + "}}" // Create a directive with the file path.
	assert.Equal(t, true, HasDirective(dataDirective), "Data directive should be detected")
	assert.Equal(t, data, ParseDirective(dataDirective), "ParseDirective should return file content")
}

// TestMatchAndReplace tests the replacement of placeholders within a string using regex.
func TestMatchAndReplace(t *testing.T) {
	test := `
adsfasdf
asdf
asd
fa {{1}}
sdf {{2 }}a dfasdf
{{3 }}asdfd {{ UIFE={This is where it starts} }}asdfasdf`

	re := regexp.MustCompile(`\{\{(.*?)\}\}`) // Compile a regex to match placeholders.
	newT := re.ReplaceAllFunc([]byte(test), func(s []byte) []byte {
		// This function is called for each regex match.
		r := strings.TrimSpace(strings.TrimLeft(strings.TrimRight(string(s), "}}"), "{{"))
		if strings.HasPrefix(r, "UIFE=") {
			r = r[5:] // Remove the "UIFE=" prefix if present.
		}

		return []byte("{{0}}") // Replace the placeholder with "{{0}}".
	})

	expected := `
adsfasdf
asdf
asd
fa {{0}}
sdf {{0}}a dfasdf
{{0}}asdfd {{0}}asdfasdf`
	assert.Equal(t, expected, string(newT), "The placeholders should be replaced with {{0}}")
}

// TestIsEmpty checks the IsEmpty method of ParseDirectiveType.
func TestIsEmpty(t *testing.T) {
	var tests = []struct {
		pdType ParseDirectiveType
		want   bool
	}{
		{"", true},
		{" ", true},
		{"file", false},
	}

	for _, tt := range tests {
		testname := string(tt.pdType)
		t.Run(testname, func(t *testing.T) {
			ans := tt.pdType.IsEmpty()
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

// TestTrimSpace checks the TrimSpace method of ParseDirectiveType.
func TestTrimSpace(t *testing.T) {
	var tests = []struct {
		pdType ParseDirectiveType
		want   ParseDirectiveType
	}{
		{" file ", "file"},
		{"  file  ", "file"},
		{"\t\nfile\n\t", "file"},
	}

	for _, tt := range tests {
		testname := string(tt.pdType)
		t.Run(testname, func(t *testing.T) {
			ans := tt.pdType.TrimSpace()
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

// TestToStringTrimLower checks the ToStringTrimLower method of ParseDirectiveType.
func TestToStringTrimLower(t *testing.T) {
	var tests = []struct {
		pdType ParseDirectiveType
		want   string
	}{
		{" FILE ", "file"},
		{" File ", "file"},
		{" fIlE ", "file"},
	}

	for _, tt := range tests {
		testname := string(tt.pdType)
		t.Run(testname, func(t *testing.T) {
			ans := tt.pdType.ToStringTrimLower()
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

//// TestParseDirectiveKeyValueWithAutoload checks the ParseDirectiveKeyValueWithAutoload function.
//func TestParseDirectiveKeyValueWithAutoload(t *testing.T) {
//	//// Mock StringParserDirectives to return a test instance
//	//originalStringParserDirectives := StringParserDirectives
//	//defer func() { StringParserDirectives = originalStringParserDirectives }()
//	//StringParserDirectives = func() *pds {
//	//	return &pds{
//	//		instances: StringParserDirectiveMap{
//	//			"test": func(pkey ParseDirectiveType, pvalue string) (key ParseDirectiveType, value string) {
//	//				return pkey, strings.ToUpper(pvalue)
//	//			},
//	//		},
//	//	}
//	//}
//
//	var tests = []struct {
//		target     string
//		doAutoLoad bool
//		wantKey    ParseDirectiveType
//		wantValue  string
//	}{
//		{"{{test:value}}", true, "test", "VALUE"},
//		{"{{test:value}}", false, "test", "value"},
//		{"{{test:}}", true, "test", ""},
//		{"{{:value}}", true, "", ""},
//		{"normal text", true, "", "normal text"},
//	}
//
//	for _, tt := range tests {
//		testname := tt.target
//		t.Run(testname, func(t *testing.T) {
//			key, value := ParseDirectiveKeyValueWithAutoload(tt.target, tt.doAutoLoad)
//			if key != tt.wantKey || value != tt.wantValue {
//				t.Errorf("got (%v, %v), want (%v, %v)", key, value, tt.wantKey, tt.wantValue)
//			}
//		})
//	}
//}
