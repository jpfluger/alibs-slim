package aclient_smtp

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Mock function for contentByIDFunc
func mockContentByIDFunc(id string) ([]byte, string, string, error) {
	if id == "valid_id" {
		return []byte("mock content"), "text/plain", "mock.txt", nil
	}
	return nil, "", "", errors.New("content not found")
}

func TestAttachment_SetContentFromBytes(t *testing.T) {
	a := &Attachment{}
	data := []byte("test content")
	a.SetContentFromBytes(data)
	expected := base64.RawStdEncoding.EncodeToString(data)
	if a.Content != expected {
		t.Errorf("SetContentFromBytes() = %v, want %v", a.Content, expected)
	}
}

func TestAttachment_GetContent(t *testing.T) {
	a := &Attachment{}
	data := []byte("test content")
	a.SetContentFromBytes(data)
	got, err := a.GetContent()
	if err != nil {
		t.Errorf("GetContent() error = %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("GetContent() = %v, want %v", string(got), string(data))
	}
}

func TestAttachment_LoadContentFromFile(t *testing.T) {
	expectFilePath := "test_data/plain.txt"
	expectContent := "file content"
	expectContentType := "text/plain; charset=utf-8"
	expectName := "plain.txt"

	a := &Attachment{
		Key: AttachmentKey("file:" + expectFilePath),
	}

	if err := a.LoadContent(); err != nil {
		t.Errorf("LoadContent() error = %v", err)
	}

	got, err := a.GetContent()
	if err != nil {
		t.Errorf("GetContent() error = %v", err)
	}
	if string(got) != expectContent {
		t.Errorf("GetContent() = %v, want %v", string(got), expectContent)
	}

	if a.ContentType != expectContentType {
		t.Errorf("ContentType = %v, want %v", a.ContentType, expectContentType)
	}
	if a.Name != filepath.Base(expectName) {
		t.Errorf("Name = %v, want %v", a.Name, filepath.Base(expectName))
	}
}

func TestAttachment_LoadContentFromID(t *testing.T) {
	SetContentByIDFunc(mockContentByIDFunc)

	a := &Attachment{
		Key: AttachmentKey("id:valid_id"),
	}

	if err := a.LoadContent(); err != nil {
		t.Errorf("LoadContent() error = %v", err)
	}

	got, err := a.GetContent()
	if err != nil {
		t.Errorf("GetContent() error = %v", err)
	}
	expectedContent := []byte("mock content")
	if string(got) != string(expectedContent) {
		t.Errorf("GetContent() = %v, want %v", string(got), string(expectedContent))
	}

	if a.ContentType != "text/plain" {
		t.Errorf("ContentType = %v, want %v", a.ContentType, "text/plain")
	}
	if a.Name != "mock.txt" {
		t.Errorf("Name = %v, want %v", a.Name, "mock.txt")
	}
}

func TestAttachment_LoadContentFromID_NoFunc(t *testing.T) {
	mu.Lock()
	contentByIDFunc = nil
	mu.Unlock()

	a := &Attachment{
		Key: AttachmentKey("id:valid_id"),
	}

	if err := a.LoadContent(); err == nil {
		t.Errorf("LoadContent() error = %v, want %v", err, "contentByIDFunc is not set")
	}
}

func TestAttachments_AddAttachmentFromFile(t *testing.T) {
	// Use the test_data/plain.txt file
	filePath := "test_data/plain.txt"
	content := []byte("file content")

	attachments := &Attachments{}
	if err := attachments.AddAttachmentFromFile(filePath); err != nil {
		t.Errorf("AddAttachmentFromFile() error = %v", err)
	}

	if len(*attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(*attachments))
	}

	attachment := (*attachments)[0]
	got, err := attachment.GetContent()
	if err != nil {
		t.Errorf("GetContent() error = %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("GetContent() = %v, want %v", string(got), string(content))
	}

	if attachment.ContentType != "text/plain; charset=utf-8" {
		t.Errorf("ContentType = %v, want %v", attachment.ContentType, "text/plain; charset=utf-8")
	}
	if attachment.Name != "plain.txt" {
		t.Errorf("Name = %v, want %v", attachment.Name, "plain.txt")
	}
}

func TestAttachments_AddAttachmentsFromDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create temporary files in the directory
	files := []string{"file1.txt", "file2.txt"}
	content := []byte("file content")
	for _, file := range files {
		tmpFile, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tmpFile.Write(content); err != nil {
			t.Fatal(err)
		}
		tmpFile.Close()
	}

	attachments := &Attachments{}
	if err := attachments.AddAttachmentsFromDirectory(tmpDir); err != nil {
		t.Errorf("AddAttachmentsFromDirectory() error = %v", err)
	}

	if len(*attachments) != len(files) {
		t.Errorf("Expected %d attachments, got %d", len(files), len(*attachments))
	}

	for i, attachment := range *attachments {
		got, err := attachment.GetContent()
		if err != nil {
			t.Errorf("GetContent() error = %v", err)
		}
		if string(got) != string(content) {
			t.Errorf("GetContent() = %v, want %v", string(got), string(content))
		}

		if attachment.ContentType != "text/plain; charset=utf-8" {
			t.Errorf("ContentType = %v, want %v", attachment.ContentType, "text/plain; charset=utf-8")
		}
		if attachment.Name != files[i] {
			t.Errorf("Name = %v, want %v", attachment.Name, files[i])
		}
	}
}
