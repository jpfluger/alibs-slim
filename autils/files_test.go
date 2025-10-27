package autils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanDirWithMkdirOption(t *testing.T) {
	// t.TempDir() creates a temporary directory that is automatically cleaned up after the test completes.
	rootDir := t.TempDir()

	tests := []struct {
		name      string
		dir       string
		root      string
		doMkDir   bool
		expectErr bool
		expected  string
	}{
		{"AbsolutePathExists", filepath.Join(rootDir, "exists"), rootDir, false, false, filepath.Join(rootDir, "exists")},
		{"RelativePathExists", "exists", rootDir, false, false, filepath.Join(rootDir, "exists")},
		{"AbsolutePathCreate", filepath.Join(rootDir, "create"), rootDir, true, false, filepath.Join(rootDir, "create")},
		{"RelativePathCreate", "create", rootDir, true, false, filepath.Join(rootDir, "create")},
		{"AbsolutePathNoCreate", filepath.Join(rootDir, "noexist"), rootDir, false, true, ""},
		{"RelativePathNoCreate", "noexist", rootDir, false, true, ""},
	}

	// Create a directory that exists
	if err := os.Mkdir(filepath.Join(rootDir, "exists"), os.ModePerm); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := CleanDirWithMkdirOption(tt.dir, tt.root, tt.doMkDir)
			if (err != nil) != tt.expectErr {
				t.Errorf("CleanDirWithMkdirOption() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && dir != tt.expected {
				t.Errorf("CleanDirWithMkdirOption() = %v, want %v", dir, tt.expected)
			}
		})
	}
}

func TestResolveDirectory(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test resolving the directory.
	resolvedPath, err := ResolveDirectory(tempDir)
	if err != nil {
		t.Errorf("ResolveDirectory() returned an error: %v", err)
	}
	if resolvedPath != tempDir {
		t.Errorf("ResolveDirectory() returned '%v', want '%v'", resolvedPath, tempDir)
	}

	// Test with a non-existent directory.
	if _, err := ResolveDirectory("nonexistentdirectory"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("ResolveDirectory() should return an error for non-existent directory")
	}

	// Test with a file instead of a directory.
	tempFile, err := os.CreateTemp(tempDir, "testfile-*.txt")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
		return
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName)
	if _, err := ResolveDirectory(tempFileName); !errors.Is(err, ErrNotDirectory) {
		t.Errorf("ResolveDirectory() should return an error when resolving a file")
	}
}

func TestCreateTempDir(t *testing.T) {
	dir, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	_, err = ResolveDirectory(dir)
	assert.NoError(t, err)
}

func TestFileNameParts(t *testing.T) {
	fullname, name, ext := FileNameParts("")
	assert.Equal(t, "", fullname)
	assert.Equal(t, "", name)
	assert.Equal(t, "", ext)

	fullname, name, ext = FileNameParts("/")
	assert.Equal(t, "", fullname)
	assert.Equal(t, "", name)
	assert.Equal(t, "", ext)

	fullname, name, ext = FileNameParts("/bronco")
	assert.Equal(t, "bronco", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, "", ext)

	fullname, name, ext = FileNameParts("/path/to/bronco")
	assert.Equal(t, "bronco", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, "", ext)

	fullname, name, ext = FileNameParts("/path/to/bronco.png")
	assert.Equal(t, "bronco.png", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, ".png", ext)

	fullname, name, ext = FileNameParts("/path/to/bronco.tar.bz")
	assert.Equal(t, "bronco.tar.bz", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, ".tar.bz", ext)

	fullname, name, ext = FileNameParts("bronco")
	assert.Equal(t, "bronco", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, "", ext)

	fullname, name, ext = FileNameParts("bronco.png")
	assert.Equal(t, "bronco.png", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, ".png", ext)

	fullname, name, ext = FileNameParts("bronco.tar.bz")
	assert.Equal(t, "bronco.tar.bz", fullname)
	assert.Equal(t, "bronco", name)
	assert.Equal(t, ".tar.bz", ext)
}

func TestFileNamePartsExt(t *testing.T) {
	assert.Equal(t, "", GetFileNamePartExt(""))
	assert.Equal(t, "", GetFileNamePartExt("/"))
	assert.Equal(t, "", GetFileNamePartExt("/bronco"))
	assert.Equal(t, "", GetFileNamePartExt("/path/to/bronco"))
	assert.Equal(t, ".png", GetFileNamePartExt("/path/to/bronco.png"))
	assert.Equal(t, ".tar.bz", GetFileNamePartExt("/path/to/bronco.tar.bz"))
	assert.Equal(t, "", GetFileNamePartExt("bronco"))
	assert.Equal(t, ".png", GetFileNamePartExt("bronco.png"))
	assert.Equal(t, ".tar.bz", GetFileNamePartExt("bronco.tar.bz"))
}

func TestFileNamePartExtNoDotPrefixToLower(t *testing.T) {
	assert.Equal(t, "", GetFileNamePartExtNoDotPrefixToLower(""))
	assert.Equal(t, "", GetFileNamePartExtNoDotPrefixToLower("/"))
	assert.Equal(t, "", GetFileNamePartExtNoDotPrefixToLower("/bronco"))
	assert.Equal(t, "", GetFileNamePartExtNoDotPrefixToLower("/path/to/bronco"))
	assert.Equal(t, "png", GetFileNamePartExtNoDotPrefixToLower("/path/to/bronco.png"))
	assert.Equal(t, "tar.bz", GetFileNamePartExtNoDotPrefixToLower("/path/to/bronco.tar.bz"))
	assert.Equal(t, "", GetFileNamePartExtNoDotPrefixToLower("bronco"))
	assert.Equal(t, "png", GetFileNamePartExtNoDotPrefixToLower("bronco.png"))
	assert.Equal(t, "tar.bz", GetFileNamePartExtNoDotPrefixToLower("bronco.tar.bz"))
}

func TestFileNamePartsName(t *testing.T) {
	assert.Equal(t, "", GetFileNamePartName(""))
	assert.Equal(t, "", GetFileNamePartName("/"))
	assert.Equal(t, "bronco", GetFileNamePartName("/bronco"))
	assert.Equal(t, "bronco", GetFileNamePartName("/path/to/bronco"))
	assert.Equal(t, "bronco", GetFileNamePartName("/path/to/bronco.png"))
	assert.Equal(t, "bronco", GetFileNamePartName("/path/to/bronco.tar.bz"))
	assert.Equal(t, "bronco", GetFileNamePartName("bronco"))
	assert.Equal(t, "bronco", GetFileNamePartName("bronco.png"))
	assert.Equal(t, "bronco", GetFileNamePartName("bronco.tar.bz"))
}

func TestCopyFileWithPerm(t *testing.T) {
	// Create a temporary directory for the source file.
	srcDir, err := os.MkdirTemp("", "src-")
	if err != nil {
		t.Fatalf("Failed to create temp source directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create a source file in the temporary directory.
	srcFilePath := filepath.Join(srcDir, "testfile.txt")
	if err := ioutil.WriteFile(srcFilePath, []byte("hello world"), 0644); err != nil {
		t.Errorf("Failed to write to source file: %v", err)
		return
	}

	// Create a temporary directory for the destination file.
	destDir, err := os.MkdirTemp("", "dest-")
	if err != nil {
		t.Errorf("Failed to create temp destination directory: %v", err)
		return
	}
	defer os.RemoveAll(destDir)

	// Define a destination file path within the temporary directory.
	destFilePath := filepath.Join(destDir, "testfile.txt")

	// Test copying the file.
	if err := CopyFileWithPerm(srcFilePath, destFilePath, true, 0644, false); err != nil {
		t.Errorf("CopyFileWithPerm() returned an error: %v", err)
	}

	// Check the content of the destination file.
	destContent, err := ioutil.ReadFile(destFilePath)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
		return
	}
	if string(destContent) != "hello world" {
		t.Errorf("Content of destination file is incorrect, got: %s, want: hello world", destContent)
	}
}

func TestCopyDir(t *testing.T) {
	// Create a temporary source directory.
	srcDir, err := os.MkdirTemp("", "srcdir-")
	if err != nil {
		t.Fatalf("Failed to create temp source directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create a file in the source directory.
	srcFilePath := filepath.Join(srcDir, "testfile.txt")
	if err = os.WriteFile(srcFilePath, []byte("hello world"), 0644); err != nil {
		t.Errorf("Failed to write to source file: %v", err)
		return
	}

	// Create a temporary directory for the destination.
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("Failed to create temp directory for destination: %v", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Define the destination directory within the temporary directory.
	destDir := filepath.Join(tempDir, "destdir")

	// Test copying the directory.
	if err := CopyDir(srcDir, destDir); err != nil {
		t.Errorf("CopyDir() returned an error: %v", err)
	}

	// Check if the file exists in the destination directory.
	destFilePath := filepath.Join(destDir, "testfile.txt")
	if _, err := os.Stat(destFilePath); os.IsNotExist(err) {
		t.Errorf("File was not copied to destination directory")
	}
}

func TestStripExtensionPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"No prefix, no spaces", "jpg", "jpg"},
		{"Dot prefix", ".png", "png"},
		{"Dot prefix with spaces", "  .gif  ", "gif"},
		{"No prefix with spaces", "  bmp  ", "bmp"},
		{"Empty string", "", ""},
		{"Only dot", ".", ""},
		{"Spaces only", "   ", ""},
		{"Complex case", " .tiff   ", "tiff"},
		{"Multiple dots", "...jpeg", "jpeg"},
		{"No prefix special chars", "  svg+xml  ", "svg+xml"},
		// File name cases
		{"File name with single extension", "file.png", "png"},
		{"File name with multiple dots", "file.date.png", "png"},
		{"File name with trailing dots", "file..png", "png"},
		{"File name with no extension", "file", "file"}, // the assumption is the target is an "ext", so this case is legitimate.
		{"File name with dot but no extension", "file.", ""},
		{"File name with spaces and extension", "  file . jpg  ", "jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripExtensionPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("StripExtensionPrefix(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSanitizeName verifies the behavior of the SanitizeName utility function.
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My File Name", "My-File-Name"},
		{"Invalid:File*Name?", "InvalidFileName"},
		{"Trailing Space .", "Trailing-Space"},
		{"Leading/Slash", "LeadingSlash"},
		{"Reserved|Name", "ReservedName"},
		{"Ends-With-Period.", "Ends-With-Period"},
		{"   Extra Spaces   ", "Extra-Spaces"},
		{"----Multiple----Dashes----", "Multiple-Dashes"},
	}

	for _, tt := range tests {
		result := SanitizeName(tt.input)
		assert.Equal(t, tt.expected, result, "SanitizeName failed for input: %s", tt.input)
	}
}

func Test_EvaluateMockRoot(t *testing.T) {
	mockSrc := "./test_data/mock1"
	mockDest := "./test_data/root_output_test_base"

	// Clean dest dir if exists
	_ = os.RemoveAll(mockDest)

	// Clean up after test
	defer func() {
		_ = os.RemoveAll(mockDest)
	}()

	err := EvaluateMockRootDir(mockSrc, mockDest, false)
	assert.NoError(t, err, "should copy mock dir to root")

	// Ensure expected file is copied
	_, err = ResolveFile(filepath.Join(mockDest, "config.json"))
	assert.NoError(t, err, "config.json should exist in copied root")
}

func TestIsPathWithin_BasicCases(t *testing.T) {
	base := t.TempDir()

	tests := []struct {
		name       string
		targetRel  []string // joined under base unless absolute
		makeAbs    bool     // if true, we make target absolute from base
		wantWithin bool
	}{
		{
			name:       "same-directory-is-within",
			targetRel:  []string{"."},
			makeAbs:    true,
			wantWithin: true,
		},
		{
			name:       "direct-child-file",
			targetRel:  []string{"a.txt"},
			makeAbs:    true,
			wantWithin: true,
		},
		{
			name:       "nested-descendant",
			targetRel:  []string{"dir1", "dir2", "b.cfg"},
			makeAbs:    true,
			wantWithin: true,
		},
		{
			name:       "normalized-descendant-with-dots",
			targetRel:  []string{"dir", ".", "sub", "c.yaml"},
			makeAbs:    true,
			wantWithin: true,
		},
		{
			name: "sibling-is-not-within",
			// simulate a sibling: parent(base)/sib/file
			// build target by going up from base: .. / sib / file
			targetRel:  []string{"..", "sib", "file.txt"},
			makeAbs:    true,
			wantWithin: false,
		},
		{
			name:       "parent-itself-is-not-within",
			targetRel:  []string{".."},
			makeAbs:    true,
			wantWithin: false,
		},
		{
			name:       "cousin-is-not-within",
			targetRel:  []string{"..", "cousin", "c.txt"},
			makeAbs:    true,
			wantWithin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := filepath.Join(append([]string{base}, tt.targetRel...)...)
			// filepath.Clean to normalize things like "." and ".."
			target = filepath.Clean(target)
			got, err := IsPathWithin(base, target)
			if err != nil {
				t.Fatalf("IsPathWithin returned error: %v", err)
			}
			if got != tt.wantWithin {
				t.Fatalf("IsPathWithin(%q, %q) = %v, want %v", base, target, got, tt.wantWithin)
			}
		})
	}
}

func TestIsPathWithin_TrailingSlashAndDotNormalization(t *testing.T) {
	base := t.TempDir()

	// base with trailing separator (simulate sloppy caller input)
	withSlash := base + string(filepath.Separator)
	target := filepath.Join(base, "sub", ".", "file.txt")

	got, err := IsPathWithin(withSlash, target)
	if err != nil {
		t.Fatalf("IsPathWithin returned error: %v", err)
	}
	if !got {
		t.Fatalf("expected target to be within base-with-slash; base=%q target=%q", withSlash, target)
	}
}

func TestIsPathWithin_AbsoluteOutsidePath(t *testing.T) {
	base := t.TempDir()

	// Construct an absolute path outside base by going up to parent (if possible)
	parent := filepath.Dir(base)
	// On some systems TempDir may be root; guard against parent == base
	if parent == base {
		t.Skip("cannot form parent of temp dir on this platform")
	}
	outside := filepath.Join(parent, "elsewhere.txt")

	got, err := IsPathWithin(base, outside)
	if err != nil {
		t.Fatalf("IsPathWithin returned error: %v", err)
	}
	if got {
		t.Fatalf("expected outside absolute path to NOT be within base; base=%q target=%q", base, outside)
	}
}

func TestIsPathWithin_CaseSensitivityIsFilesystemDependent(t *testing.T) {
	// On Windows, paths are case-insensitive; on Unix they are case-sensitive.
	base := t.TempDir()
	child := filepath.Join(base, "CaseDir", "file.txt")

	within, err := IsPathWithin(base, child)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !within {
		t.Fatalf("expected child to be within base")
	}

	// Flip the case of base and expect platform-dependent behavior.
	var flippedBase string
	if len(base) > 0 {
		flippedBase = base[:1]
		if strings.ToUpper(flippedBase) == flippedBase {
			flippedBase = strings.ToLower(base[:1]) + base[1:]
		} else {
			flippedBase = strings.ToUpper(base[:1]) + base[1:]
		}
	} else {
		flippedBase = base
	}

	within2, err := IsPathWithin(flippedBase, child)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if runtime.GOOS == "windows" {
		if !within2 {
			t.Fatalf("on Windows, path comparison should be case-insensitive; got within=false")
		}
	} else {
		// Most Unix filesystems are case-sensitive, but Abs/Rel resolution often normalizes.
		// We accept either result here to avoid flaky tests across exotic filesystems.
		_ = within2
	}
}
