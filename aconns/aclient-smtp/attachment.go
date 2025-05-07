package aclient_smtp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jhillyerd/enmime/v2"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Attachment can be used instead of enmime.Part.
// Why? Because different systems will store attachments differently.
// This attachment struct offers various ways to do so.
// If Key is ATTACHMENTKEY_NONE, then the Content is encoded directly as base64.
// If Key is ATTACHMENTKEY_FILE, then the target from key.GetParts() contains the file to load into Content.
// If Key is ATTACHMENTKEY_ID, then the target from key.GetParts() contains the id your app maps to content
// within your system.
type Attachment struct {
	Key         AttachmentKey `json:"key,omitempty"`
	ContentType string        `json:"contentType"` // Content-Type
	Name        string        `json:"name"`
	//Filename string `json:"filename"`
	// stored as base64, using base64.RawStdEncoding
	// https://golangbyexample.com/base64-golang/
	Content string `json:"content,omitempty"`

	// epart is optionally set and overrides
	// ContentType, Name and Content when using
	// GetContentType, GetName and GetContent
	epart *enmime.Part
}

// Global variable for the function to retrieve content by ID
var (
	contentByIDFunc func(id string) (data []byte, contentType string, name string, err error)
	mu              sync.Mutex
)

// SetContentByIDFunc sets the global function to retrieve content by ID
func SetContentByIDFunc(f func(id string) (data []byte, contentType string, name string, err error)) {
	mu.Lock()
	defer mu.Unlock()
	contentByIDFunc = f
}

func (a *Attachment) Validate() error {
	if a.epart != nil {
		// ContentType if empty by default is text in enmime
		// FileName is not guaranteed in enmime.
		if len(a.epart.Content) == 0 {
			return fmt.Errorf("attachment part content is empty")
		}
		return nil
	}
	if a.Content == "" {
		return fmt.Errorf("attachment content is empty")
	}
	return nil
}

// LoadContent loads the content based on the AttachmentKey
func (a *Attachment) LoadContent() error {
	if a.epart != nil {
		return a.Validate()
	}
	if a.Key.IsEmpty() {
		return a.Validate()
	}

	key, target, err := a.Key.GetParts()
	if err != nil {
		return err
	}

	switch key {
	case ATTACHMENTKEY_FILE:
		return a.loadContentFromFile(target)
	case ATTACHMENTKEY_ID:
		return a.loadContentFromID(target)
	default:
		return errors.New("unsupported attachment key")
	}
}

// loadContentFromFile loads the content from a file
func (a *Attachment) loadContentFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	a.SetContentFromBytes(data)

	if a.ContentType == "" {
		// used by enmime/builder.go, see "ctype := mime.TypeByExtension(filepath.Ext(name))"
		a.ContentType = mime.TypeByExtension(filepath.Ext(filePath))
	}
	if a.Name == "" {
		// used by enmime/builder.go, see "name := filepath.Base(path)"
		a.Name = filepath.Base(filePath)
	}

	return nil
}

// loadContentFromID loads the content from an ID using the global function
func (a *Attachment) loadContentFromID(id string) error {
	mu.Lock()
	defer mu.Unlock()
	if contentByIDFunc == nil {
		return errors.New("contentByIDFunc is not set")
	}
	data, contentType, name, err := contentByIDFunc(id)
	if err != nil {
		return err
	}
	a.SetContentFromBytes(data)

	if a.ContentType == "" {
		a.ContentType = contentType
	}
	if a.Name == "" {
		a.Name = name
	}

	return nil
}

// SetEnmimePart sets the enmime.Part for the attachment
func (a *Attachment) SetEnmimePart(part *enmime.Part) {
	a.epart = part
}

// GetEnmimePart returns the enmime.Part if it is set
func (a *Attachment) GetEnmimePart() *enmime.Part {
	return a.epart
}

// HasEnmimePart checks if the enmime.Part is set
func (a *Attachment) HasEnmimePart() bool {
	return a.epart != nil
}

// SetContentFromBytes sets the content from a byte array and encodes it as base64
func (a *Attachment) SetContentFromBytes(data []byte) {
	a.Content = base64.RawStdEncoding.EncodeToString(data)
}

// GetContent returns the content based on the enmime.Part if it is set
func (a *Attachment) GetContent() ([]byte, error) {
	if a.HasEnmimePart() {
		return a.epart.Content, nil
	}
	return base64.RawStdEncoding.DecodeString(a.Content)
}

// GetContentType returns the content type based on the enmime.Part if it is set
func (a *Attachment) GetContentType() string {
	if a.HasEnmimePart() {
		return a.epart.ContentType
	}
	return a.ContentType
}

// GetName returns the name based on the enmime.Part if it is set
func (a *Attachment) GetName() string {
	if a.HasEnmimePart() {
		return a.epart.FileName
	}
	return a.Name
}

// Clone creates a deep copy of an Attachment, including its enmime.Part if necessary.
func (att *Attachment) Clone() *Attachment {
	if att == nil || strings.TrimSpace(att.Content) == "" {
		return nil
	}
	clone := &Attachment{
		Key:         att.Key,
		ContentType: att.ContentType,
		Name:        att.Name,
		Content:     att.Content,
	}

	// Clone the enmime.Part if it's set
	if att.epart != nil {
		clone.epart = att.epart.Clone(nil) // assuming enmime.Part has a Clone method; otherwise, handle appropriately
	}

	return clone
}

type Attachments []*Attachment

// AddAttachmentFromFile creates a new Attachment from a file and adds it to the Attachments slice
func (as *Attachments) AddAttachmentFromFile(filePath string) error {
	attachment := &Attachment{
		Key: AttachmentKey("file:" + filePath),
	}
	if err := attachment.LoadContent(); err != nil {
		return err
	}
	*as = append(*as, attachment)
	return nil
}

// AddAttachmentsFromDirectory creates new Attachments from all files in a directory and adds them to the Attachments slice
func (as *Attachments) AddAttachmentsFromDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err := as.AddAttachmentFromFile(filePath); err != nil {
				return err
			}
		}
		return nil
	})
}

// Clone creates a deep copy of the Attachments slice, which contains Attachment elements.
func (as Attachments) Clone() Attachments {
	if as == nil || len(as) == 0 {
		return nil
	}
	clones := Attachments{}
	for _, attachment := range as {
		if clone := attachment.Clone(); clone != nil {
			clones = append(clones, clone)
		}
	}
	return clones
}
