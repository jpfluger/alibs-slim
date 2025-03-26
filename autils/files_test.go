package autils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
