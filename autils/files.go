package autils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

// Define file permissions constants with comments explaining their purpose.
// See SetProcessUmask to set "0" globally prior to this function to have consistent behavior.
const (
	PATH_CHMOD_DIR             os.FileMode = 0755 // Default for generic dirs
	PATH_CHMOD_FILE            os.FileMode = 0644 // Default for generic files
	PATH_CHMOD_DIR_LIMIT       os.FileMode = 0744 // Limited dir perms
	PATH_CHMOD_DIR_FULL_PERMS  os.FileMode = 0777 // Full perms (rare)
	PATH_CHMOD_SCRIPTS         os.FileMode = 0744 // Scripts executable
	PATH_CHMOD_DIR_SECRETS     os.FileMode = 0700 // Secrets: owner only
	PATH_CHMOD_FILE_SECRETS    os.FileMode = 0600 // Secret files
	PATH_CHMOD_DIR_OWNER_GROUP os.FileMode = 0750 // Owner rwx, group rx
	PATH_CHMOD_PATH_RO         os.FileMode = 0444 // Read-only access for any path
)

// SetProcessUmask sets the process-wide file mode creation mask (umask).
//
// The umask controls the default permission bits for **newly created files and directories**
// in the current process. It is **subtracted** from the default permissions requested by
// system calls like `os.Mkdir`, `os.Create`, or `net.Listen("unix")`.
//
// Typical uses:
//
//   - ✅ Secure runtime temp files: Set `umask(0077)` so all created files/dirs are owner-only (`700` or `600`).
//   - ✅ Allow group collaboration: Set `umask(0027)` so files are owner+group accessible (`750` or `640`).
//   - ✅ Sockets: UNIX domain sockets inherit the umask — adjust to control who can connect.
//
// **Example:** A daemon might tighten its umask early during `main()` to ensure any temp config,
// secrets, or IPC sockets are not accidentally world-readable.
//
// Note: This affects only the **current process** and child processes that inherit the environment.
// It does **not** affect other processes.
//
// For permanent system-wide behavior, the OS or container runtime usually sets the default umask.
func SetProcessUmask(umask int) {
	syscall.Umask(umask)
}

// Exists checks if the given path exists.
func Exists(target string) bool {
	_, err := os.Stat(target)
	return err == nil
}

// FileExists checks if the given file exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// DirExists checks if the given directory exists.
func DirExists(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

// ResolveFile checks if the target is a file and returns its clean path.
func ResolveFile(target string) (string, error) {
	if target == "" {
		return "", errors.New("file path not found")
	}
	target = path.Clean(target)
	info, err := os.Stat(target)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("expected a file but found a directory")
	}
	return target, nil
}

// ErrNotDirectory is the error returned when a directory is expected but not found.
var ErrNotDirectory = errors.New("path is not a directory")

// ResolveDirectory checks if the target is a directory and returns its clean path.
func ResolveDirectory(target string) (string, error) {
	if target == "" {
		return "", errors.New("directory path not found")
	}
	target = filepath.Clean(target)
	info, err := os.Stat(target)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", ErrNotDirectory // Use the ErrNotDirectory error here.
	}
	return target, nil
}

// IsFile checks if the target is a file.
func IsFile(target string) (bool, error) {
	if target == "" {
		return false, errors.New("file path not found")
	}
	target = path.Clean(target)
	info, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// IsDirectory checks if the target is a directory.
func IsDirectory(target string) (bool, error) {
	if target == "" {
		return false, errors.New("directory path not found")
	}
	target = path.Clean(target)
	info, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// DirectoryRecreate removes and recreates a directory with the specified mode.
func DirectoryRecreate(dir string, mode os.FileMode) error {
	if Exists(dir) {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("cannot delete the existing directory: %v", err)
		}
	}
	return os.Mkdir(dir, mode)
}

// DeleteDirectory deletes a directory if the deleteDirectory flag is true.
func DeleteDirectory(dir string, dirCommonNameForErrorMessage string, deleteDirectory bool) (string, error) {
	if dirCommonNameForErrorMessage == "" {
		dirCommonNameForErrorMessage = "directory"
	}
	if dir == "" {
		return "", fmt.Errorf("the %s has not been specified", dirCommonNameForErrorMessage)
	}
	dir, err := ResolveDirectory(dir)
	if err != nil {
		return "", err
	}
	if !deleteDirectory {
		return "", fmt.Errorf("the %s already exists at %s", dirCommonNameForErrorMessage, dir)
	}
	if err := os.RemoveAll(dir); err != nil {
		return "", fmt.Errorf("error encountered when deleting the %s at \"%s\": %s", dirCommonNameForErrorMessage, dir, err.Error())
	}
	return path.Clean(dir), nil
}

// TempDirOptions defines options for creating temporary directories.
type TempDirOptions struct {
	DirRoot      string // Root directory for the temp dir, defaults to the system temp dir if empty.
	Name         string // Name of the temp dir, auto-created as "tmp-UUID" if empty.
	AppendUUIDv4 bool   // If true and name is not empty, then append "-UUID" to the name.
}

// CreateTempDir creates a temporary directory with default options.
func CreateTempDir() (string, error) {
	return CreateTempDirWithOptions(nil)
}

// CreateTempDirWithOptions creates a temporary directory with the specified options.
func CreateTempDirWithOptions(options *TempDirOptions) (string, error) {
	if options == nil {
		options = &TempDirOptions{}
	}
	dir := strings.TrimSpace(options.DirRoot)
	name := strings.TrimSpace(options.Name)
	if name == "" {
		name = "tmp-" + NewUUIDAsString()
	} else if options.AppendUUIDv4 {
		name += "-" + NewUUIDAsString()
	}
	return os.MkdirTemp(dir, name)
}

// FileNameParts extracts the full name, name without extension, and extension from a file path.
func FileNameParts(target string) (fullname, name, ext string) {
	// Trim any leading and trailing whitespace from the target path.
	target = strings.TrimSpace(target)

	// Return empty strings if the target is empty or just the root directory.
	if target == "" || target == "/" {
		return "", "", ""
	}

	// If the target contains directory separators, extract just the file name.
	if strings.Contains(target, "/") {
		target = filepath.Base(target)
	}

	// Initialize variables to hold the parts of the file name.
	var text string
	newname := target

	// Count the number of dots in the file name to handle multiple extensions.
	max := strings.Count(target, ".")

	// Iterate over each dot, extracting the extension each time.
	for ii := 0; ii < max; ii++ {
		text = filepath.Ext(newname)
		newname = strings.TrimSuffix(newname, text)
		ext = text + ext // Prepend the extracted extension to build the full extension.
	}

	// The remaining part of the file name is the name without the extension.
	name = newname
	fullname = target // The full name is the original target path.

	return fullname, name, ext
}

// GetFileNamePartExt extracts the extension from a file path.
// It handles multiple extensions (e.g., .tar.gz) and returns the full extension.
func GetFileNamePartExt(target string) (ext string) {
	// Trim any leading and trailing whitespace from the target path.
	target = strings.TrimSpace(target)

	// Return an empty string if the target is empty or just the root directory.
	if target == "" || target == "/" {
		return ""
	}

	// If the target contains directory separators, extract just the file name.
	if strings.Contains(target, "/") {
		target = filepath.Base(target)
	}

	// Count the number of dots in the file name to handle multiple extensions.
	max := strings.Count(target, ".")
	var text string
	newname := target

	// Iterate over each dot, extracting the extension each time.
	for ii := 0; ii < max; ii++ {
		text = filepath.Ext(newname)
		newname = strings.TrimSuffix(newname, text)
		ext = text + ext // Prepend the extracted extension to build the full extension.
	}

	return ext
}

// QuickBaseName removes a single-dot extension from a file name.
func QuickBaseName(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n == -1 {
		return s
	}
	return s[:n]
}

// GetFileNamePartExtNoDotPrefixToLower extracts the extension from a file path,
// converts it to lowercase, and removes the dot prefix.
func GetFileNamePartExtNoDotPrefixToLower(target string) (ext string) {
	// Use the GetFileNamePartExt function to extract the extension from the target path.
	ext = GetFileNamePartExt(target)

	// If an extension is found, convert it to lowercase and remove the dot prefix.
	if ext != "" {
		ext = strings.ToLower(ext)
		if strings.HasPrefix(ext, ".") {
			ext = strings.TrimPrefix(ext, ".")
		}
	}

	return ext
}

// GetFileNamePartName extracts the file name without the extension from a file path.
func GetFileNamePartName(target string) (name string) {
	// Trim any leading and trailing whitespace from the target path.
	target = strings.TrimSpace(target)

	// Return an empty string if the target is empty or just the root directory.
	if target == "" || target == "/" {
		return ""
	}

	// If the target contains directory separators, extract just the file name.
	if strings.Contains(target, "/") {
		target = filepath.Base(target)
	}

	// Initialize variables to hold the parts of the file name.
	var text, ext string
	newname := target

	// Count the number of dots in the file name to handle multiple extensions.
	max := strings.Count(target, ".")

	// Iterate over each dot, extracting the extension each time.
	for ii := 0; ii < max; ii++ {
		text = filepath.Ext(newname)
		newname = strings.TrimSuffix(newname, text)
		ext = text + ext // Accumulate the extension parts.
	}

	// The remaining part of the file name is the name without the extension.
	name = newname
	return
}

// StripExtensionPrefix removes the dot prefix from a file extension.
func StripExtensionPrefix(input string) string {
	// Trim any leading and trailing whitespace.
	input = strings.TrimSpace(input)

	if strings.HasSuffix(input, ".") {
		return ""
	}

	// Extract the extension if the input contains a file name.
	if idx := strings.LastIndex(input, "."); idx != -1 {
		input = input[idx+1:]
	}

	return strings.TrimSpace(input)
}

// SanitizeName ensures that a given string is safe for use as a file name and url.
// It performs the following operations:
// - Replaces spaces with dashes.
// - Removes invalid characters such as < > : " / \ | ? *.
// - Trims trailing dashes, periods, or whitespace.
// The resulting string is suitable for use in file systems and URLs.
func SanitizeName(name string) string {
	// Replace spaces with dashes
	sanitizedName := strings.ReplaceAll(name, " ", "-")

	// Remove characters not allowed in file names
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitizedName = invalidChars.ReplaceAllString(sanitizedName, "")

	// Consolidate multiple dashes into a single dash
	multipleDashes := regexp.MustCompile(`-+`)
	sanitizedName = multipleDashes.ReplaceAllString(sanitizedName, "-")

	// Trim any trailing dashes, periods, or whitespace
	sanitizedName = strings.Trim(sanitizedName, "-. ")

	return sanitizedName
}

// CopyFile copies a file from src to dst, setting permissions to PATH_CHMOD_FILE.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, PATH_CHMOD_FILE)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// MoveFileWithPerm copies a file from srcPath to destPath with the specified permissions, then deletes the srcPath.
func MoveFileWithPerm(srcPath, destPath string, doOverwrite bool, fileMode os.FileMode, includeNonRegularFiles bool) error {
	if err := CopyFileWithPerm(srcPath, destPath, doOverwrite, fileMode, includeNonRegularFiles); err != nil {
		return err
	}
	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("cannot remove existing file at \"%s\": %v", srcPath, err)
	}
	return nil
}

// CopyFileWithPerm copies a file from srcPath to destPath with the specified permissions.
func CopyFileWithPerm(srcPath, destPath string, doOverwrite bool, fileMode os.FileMode, includeNonRegularFiles bool) error {
	// Check if source path is defined and exists.
	if srcPath == "" {
		return errors.New("source path is not defined")
	}
	srcPath, err := ResolveFile(srcPath)
	if err != nil {
		return fmt.Errorf("source path not found: %v", err)
	}

	// Check if destination path is defined.
	if destPath == "" {
		return errors.New("destination path is not defined")
	}

	// Check if destination path exists and handle overwriting.
	if Exists(destPath) && !doOverwrite {
		return errors.New("destination path already exists")
	}

	// Open source file.
	srcStream, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source path: %v", err)
	}
	defer srcStream.Close()

	// Open destination file.
	destStream, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return fmt.Errorf("cannot create destination path: %v", err)
	}
	defer destStream.Close()

	// Perform the file copy.
	if _, err := io.Copy(destStream, srcStream); err != nil {
		return fmt.Errorf("cannot copy from source to destination: %v", err)
	}

	// Flush to stable storage.
	if err := destStream.Sync(); err != nil {
		return fmt.Errorf("cannot flush to stable storage: %v", err)
	}

	return nil
}

// CopyDirOpts holds options for CopyDir.
type CopyDirOpts struct {
	IgnoreIfDestExists bool   // If true, skip copy if dest exists
	Timestamp          string // Optional timestamp suffix for dest (e.g., "20060102-150405"); empty means no suffix
	// Future options: e.g., PreserveTimes bool, Exclude []string
}

// CopyDir recursively copies a directory tree, preserving permissions.
// Use opts to customize (e.g., ignore existing dest or add timestamp for backups).
func CopyDir(src, dst string, opts ...CopyDirOpts) error {
	var opt CopyDirOpts
	if len(opts) > 0 {
		opt = opts[0] // Use first opts; extend for merging if needed
	}

	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// Append timestamp if provided (for unique backups)
	if opt.Timestamp != "" {
		dst = fmt.Sprintf("%s-%s", dst, opt.Timestamp)
	}

	// Check source
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return errors.New("source is not a directory")
	}

	// Handle destination
	if Exists(dst) {
		if opt.IgnoreIfDestExists {
			return nil // Skip copy
		}
		// Proceed to merge contents if not ignoring
	} else {
		if err := os.MkdirAll(dst, si.Mode()); err != nil {
			return err
		}
	}

	// Read source entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recurse for subdirs with same opts
			if err := CopyDir(srcPath, dstPath, opt); err != nil {
				return err
			}
		} else {
			// Skip symlinks
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			// Copy files
			if err := copyFileR(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func copyFileR(src, dst string) (err error) {
	// Open source
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create dest
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil && err == nil { // Only set if no prior err
			err = e
		}
	}()

	// Copy contents
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// Sync for data integrity
	if err = out.Sync(); err != nil {
		return err
	}

	// Preserve mode
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err = os.Chmod(dst, si.Mode()); err != nil {
		return err
	}
	return nil
}

// AppendDataNewLine appends data to a file, adding a newline at the end.
func AppendDataNewLine(path string, data []byte, fileMode os.FileMode) error {
	// Open file for appending.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	if err != nil {
		return fmt.Errorf("cannot open file at %s; %v", path, err)
	}
	defer f.Close()

	// Add a newline to the data.
	data = append(data, '\n')

	// Write the data to the file.
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("cannot write data at %s; %v", path, err)
	}

	return nil
}

// AppendFile appends data to a file with the specified mode.
func AppendFile(filepath string, data []byte, mode os.FileMode) error {
	// Resolve the file and get its mode if it exists.
	if _, err := ResolveFile(filepath); err != nil {
		if mode == 0 {
			mode = PATH_CHMOD_FILE
		}
	} else {
		if fi, err := os.Stat(filepath); err != nil {
			return err
		} else {
			mode = fi.Mode()
		}
	}

	// Open the file for appending.
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return fmt.Errorf("cannot open file at %s; %v", filepath, err)
	}
	defer file.Close()

	// Write the data to the file.
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("cannot write data at %s; %v", filepath, err)
	}

	return nil
}

// IsFileContentIdentical checks if the content of two files is identical.
func IsFileContentIdentical(file1 string, file2 string) (bool, error) {
	// Read the content of the first file.
	f1, err1 := os.ReadFile(file1)
	if err1 != nil {
		return false, fmt.Errorf("cannot read file1 at %s; %v", file1, err1)
	}

	// Read the content of the second file.
	f2, err2 := os.ReadFile(file2)
	if err2 != nil {
		return false, fmt.Errorf("cannot read file2 at %s; %v", file2, err2)
	}

	// Compare the content of the two files.
	return bytes.Equal(f1, f2), nil
}

// CopyDirFromTo contains the "to" path "from" which data is to be copied,.
type CopyDirFromTo struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// CopyDirsFromTo contains a slice of CopyDirFromTo structs for copying multiple directories.
type CopyDirsFromTo []*CopyDirFromTo

// RunCopy performs the copying of directories as specified in CopyDirsFromTo.
func (c CopyDirsFromTo) RunCopy() error {
	for _, cd := range c {
		// Ensure the source directory exists.
		if _, err := ResolveDirectory(cd.From); err != nil {
			return fmt.Errorf("cannot find dir to copy at '%s'; %v", cd.From, err)
		}

		// Remove the destination directory if it exists.
		if _, err := ResolveDirectory(cd.To); err == nil {
			if err := os.RemoveAll(cd.To); err != nil {
				return fmt.Errorf("cannot delete '%s'; %v", cd.To, err)
			}
		}

		// Perform the directory copy.
		if err := CopyDir(cd.From, cd.To); err != nil {
			return fmt.Errorf("cannot copy '%s' to '%s'; %v", cd.From, cd.To, err)
		}
	}
	return nil
}

// ReadFileTrimSpace reads a file, trims the whitespace, and returns the content as a string.
func ReadFileTrimSpace(filePath string) string {
	if strings.TrimSpace(filePath) != "" {
		if _, err := ResolveFile(filePath); err == nil {
			if b, err := os.ReadFile(filePath); err == nil {
				return strings.TrimSpace(string(b))
			}
		}
	}
	return ""
}

// ErrEmptyPath indicates that the file path provided is empty.
var ErrEmptyPath = errors.New("file path is empty")

// ErrFileNotResolved indicates that the file could not be resolved.
var ErrFileNotResolved = errors.New("file could not be resolved")

// ReadFileTrimSpaceWithError reads a file, trims the whitespace, and returns the content as a string along with an error if any.
func ReadFileTrimSpaceWithError(filePath string) (string, error) {
	trimmedPath := strings.TrimSpace(filePath)
	if trimmedPath == "" {
		return "", ErrEmptyPath
	}

	if _, err := ResolveFile(trimmedPath); err != nil {
		return "", ErrFileNotResolved
	}

	b, err := os.ReadFile(trimmedPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// ReadByFilePathWithDirOption reads a file by combining a directory option with a partial path.
func ReadByFilePathWithDirOption(targetPathAbsOrPartial string, dirOption string) ([]byte, error) {
	targetPath, err := CleanFilePathWithDirOption(targetPathAbsOrPartial, dirOption)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file where path=%s", targetPath)
	}

	return b, nil
}

// CleanFilePathWithDirOption combines a directory option with a partial path and cleans the result.
func CleanFilePathWithDirOption(targetPathAbsOrPartial string, dirOption string) (string, error) {
	var targetPath string
	targetPathAbsOrPartial = strings.TrimSpace(targetPathAbsOrPartial)
	if targetPathAbsOrPartial == "" {
		return "", fmt.Errorf("target path is empty")
	}

	// Check if the targetPath is an absolute path.
	if path.IsAbs(targetPathAbsOrPartial) {
		cleaned, err := ResolveFile(targetPathAbsOrPartial)
		if err != nil {
			return "", fmt.Errorf("invalid target path where path=%s; %v", targetPathAbsOrPartial, err)
		}
		targetPath = cleaned
	} else {
		// Try relative path
		cleaned, err := ResolveFile(targetPathAbsOrPartial)
		if err == nil {
			return cleaned, nil
		}

		dirOption = strings.TrimSpace(dirOption)
		if dirOption == "" {
			return "", fmt.Errorf("dir option is empty")
		}

		// Combine with the directory option.
		cleaned, err = ResolveFile(path.Join(dirOption, targetPathAbsOrPartial))
		if err != nil {
			return "", fmt.Errorf("invalid target path where dirOption=%s and path=%s; %v", dirOption, targetPathAbsOrPartial, err)
		}
		targetPath = cleaned
	}

	return targetPath, nil
}

// CleanDirWithMkdirOption validates and combines a directory path with the root directory.
func CleanDirWithMkdirOption(dir, root string, doMkDir bool) (string, error) {
	dir = strings.TrimSpace(dir)
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(root, dir)
	}
	if !DirExists(dir) {
		if !doMkDir {
			return "", fmt.Errorf("dir does not exist; dir=%s", dir)
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory '%s': %v", dir, err)
		}
	}
	return dir, nil
}

// CleanDirsWithMkdirOption validates and combines a directory path with the root directory.
func CleanDirsWithMkdirOption(dirList []string, root string, doMkDir bool) error {
	for _, dir := range dirList {
		_, err := CleanDirWithMkdirOption(dir, root, doMkDir)
		if err != nil {
			return fmt.Errorf("failed to process directory '%s': %v", dir, err)
		}
	}
	return nil
}

// FNScanLine is a function type for processing a line of text.
type FNScanLine func(string) error

// ScanFileByLine reads a file line by line and processes each line using the provided FNScanLine function.
func ScanFileByLine(filePath string, fnScanLine FNScanLine) error {
	if fnScanLine == nil {
		return fmt.Errorf("fnScanLine is nil")
	}

	// Open the file for reading.
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a scanner to read the file line by line.
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		// Process each line using the provided function.
		if err = fnScanLine(sc.Text()); err != nil {
			return err
		}
	}
	if err = sc.Err(); err != nil {
		return fmt.Errorf("scan file error: %v", err)
	}
	return nil
}

// ListFilenamesInDir returns a list of filenames (not full paths) in the given directory.
func ListFilenamesInDir(dirPath string) []string {
	return ListFilenamesInDirWithExtensions(dirPath)
}

// ListFilenamesInDirWithExtensions returns a list of filenames (not full paths) in the given directory
// and filters them by the given extensions. If no extensions are passed, it returns all files.
func ListFilenamesInDirWithExtensions(dirPath string, exts ...string) []string {
	var filenames []string

	if strings.TrimSpace(dirPath) == "" {
		return filenames
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return filenames
	}

	extMap := make(map[string]struct{})
	for _, ext := range exts {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext != "" {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			extMap[ext] = struct{}{}
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(extMap) == 0 {
				filenames = append(filenames, name)
			} else {
				ext := strings.ToLower(filepath.Ext(name))
				if _, ok := extMap[ext]; ok {
					filenames = append(filenames, name)
				}
			}
		}
	}

	return filenames
}

// EvaluateMockRootDir sets up the real RootDir based on a mock seed directory.
// - If DeleteMockRoot is true, RootDir will be forcefully wiped first.
// - If RootDir does not exist, MockRootDir will be copied into it.
func EvaluateMockRootDir(mockDir string, rootDir string, deleteRoot bool) error {
	mockDir = strings.TrimSpace(mockDir)
	rootDir = strings.TrimSpace(rootDir)

	if mockDir == "" || rootDir == "" {
		return nil // nothing to do
	}

	// Check if mockDir exists
	_, err := ResolveDirectory(mockDir)
	if err != nil {
		return fmt.Errorf("mock root dir does not exist: %v", err)
	}

	// Check if rootDir exists
	dirRootExists := false
	if _, err = ResolveDirectory(rootDir); err == nil {
		dirRootExists = true
	}

	// Optionally delete existing root
	if dirRootExists && deleteRoot {
		if err = os.RemoveAll(rootDir); err != nil {
			return fmt.Errorf("failed to delete existing rootDir (%s): %w", rootDir, err)
		}
		dirRootExists = false
	}

	// Copy mock -> root if root is missing
	if !dirRootExists {
		if err = CopyDir(mockDir, rootDir); err != nil {
			return fmt.Errorf("failed to copy mock root (%s) to rootDir (%s): %w", mockDir, rootDir, err)
		}
	}

	return nil
}

// IsPathWithin returns true if target is the same as base or a descendant of base.
// Both base and target may be relative or absolute; they will be resolved to absolute paths.
// Symlinks are not resolved here; if you need that, wrap Abs with EvalSymlinks.
func IsPathWithin(base, target string) (bool, error) {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return false, err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return false, err
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return false, err
	}

	// Normalize separators and case for Windows
	sep := string(filepath.Separator)
	if runtime.GOOS == "windows" {
		absBase = strings.ToLower(strings.ReplaceAll(absBase, "/", sep))
		absTarget = strings.ToLower(strings.ReplaceAll(absTarget, "/", sep))
		rel = strings.ToLower(strings.ReplaceAll(rel, "/", sep))
	}

	if rel == "." {
		return true, nil
	}
	// rel starting with ".." means target is outside base
	if rel == ".." || strings.HasPrefix(rel, ".."+sep) {
		return false, nil
	}
	return true, nil
}

// IsUnderRoot checks if the given path is under (or equal to) the root directory.
// It resolves both paths to absolute and checks if the relative path does not start with "..".
// Returns true if path is under root, false otherwise, and an error if resolution fails.
func IsUnderRoot(root, path string) (bool, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return false, fmt.Errorf("failed to resolve root %s: %w", root, err)
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("failed to resolve path %s: %w", path, err)
	}

	rel, err := filepath.Rel(rootAbs, pathAbs)
	if err != nil {
		return false, nil // Not under root if Rel fails
	}

	return !strings.HasPrefix(rel, ".."), nil
}

// IsDirEmpty checks if the directory at the given path exists and contains no files or subdirectories.
// It returns true if the directory is empty, false if it exists but is not empty, and an error if the path
// does not exist or cannot be accessed.
func IsDirEmpty(path string) (bool, error) {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // Does not exist, treat as "empty" for auto-creation purposes
		}
		return false, fmt.Errorf("failed to open directory %s: %w", path, err)
	}
	defer dir.Close()

	// Read one entry to check if empty
	_, err = dir.Readdirnames(1)
	if err == io.EOF {
		return true, nil // No entries, empty
	} else if err != nil {
		return false, fmt.Errorf("failed to read directory %s: %w", path, err)
	}
	return false, nil // Has at least one entry
}
