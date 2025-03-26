package aclient_smtp

import (
	"net/mail"
	"testing"
)

func TestMailPiece_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mail    MailPiece
		wantErr bool
	}{
		{
			name: "Valid mail piece",
			mail: MailPiece{
				From: mail.Address{
					Name:    "Sender",
					Address: "sender@example.com",
				},
				To: []mail.Address{
					{
						Name:    "Recipient",
						Address: "recipient@example.com",
					},
				},
				Subject: "Test Subject",
				Text:    "Test Body",
			},
			wantErr: false,
		},
		{
			name: "Invalid from address",
			mail: MailPiece{
				From: mail.Address{
					Name:    "Sender",
					Address: "invalid-email",
				},
				To: []mail.Address{
					{
						Name:    "Recipient",
						Address: "recipient@example.com",
					},
				},
				Subject: "Test Subject",
				Text:    "Test Body",
			},
			wantErr: true,
		},
		{
			name: "No recipients",
			mail: MailPiece{
				From: mail.Address{
					Name:    "Sender",
					Address: "sender@example.com",
				},
				Subject: "Test Subject",
				Text:    "Test Body",
			},
			wantErr: true,
		},
		{
			name: "Empty subject not allowed",
			mail: MailPiece{
				From: mail.Address{
					Name:    "Sender",
					Address: "sender@example.com",
				},
				To: []mail.Address{
					{
						Name:    "Recipient",
						Address: "recipient@example.com",
					},
				},
				AllowEmptySubject: false,
				Text:              "Test Body",
			},
			wantErr: true,
		},
		{
			name: "No text or HTML content",
			mail: MailPiece{
				From: mail.Address{
					Name:    "Sender",
					Address: "sender@example.com",
				},
				To: []mail.Address{
					{
						Name:    "Recipient",
						Address: "recipient@example.com",
					},
				},
				Subject: "Test Subject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mail.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MailPiece.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create a sample MailPiece for testing
func createTestMailPiece() *MailPiece {
	return &MailPiece{
		From:              mail.Address{Name: "Alice", Address: "alice@example.com"},
		To:                []mail.Address{{Name: "Bob", Address: "bob@example.com"}},
		CC:                []mail.Address{{Name: "Charlie", Address: "charlie@example.com"}},
		BCC:               []mail.Address{{Name: "David", Address: "david@example.com"}},
		AllowEmptySubject: true,
		Subject:           "Test Subject",
		AllowNoTextHMTL:   true,
		Text:              "This is a text body.",
		HTML:              "<p>This is an HTML body.</p>",
		AllowAttachments:  true,
		Attachments: Attachments{
			&Attachment{
				Key:         "attachment1",
				ContentType: "text/plain",
				Name:        "file1.txt",
				Content:     "SGVsbG8gd29ybGQ=", // "Hello world" in base64
			},
		},
		Inlines: Attachments{
			&Attachment{
				Key:         "inline1",
				ContentType: "image/png",
				Name:        "image1.png",
				Content:     "iVBORw0KGgoAAAANSUhEUgAAAAUA", // Mock base64 content
			},
		},
	}
}

func TestMailPieceClone(t *testing.T) {
	original := createTestMailPiece()
	clone := original.Clone()

	// Check if basic fields are copied correctly
	if clone.From != original.From {
		t.Errorf("From field not copied correctly")
	}
	if clone.AllowEmptySubject != original.AllowEmptySubject {
		t.Errorf("AllowEmptySubject field not copied correctly")
	}
	if clone.Subject != original.Subject {
		t.Errorf("Subject field not copied correctly")
	}
	if clone.AllowNoTextHMTL != original.AllowNoTextHMTL {
		t.Errorf("AllowNoTextHMTL field not copied correctly")
	}
	if clone.Text != original.Text {
		t.Errorf("Text field not copied correctly")
	}
	if clone.HTML != original.HTML {
		t.Errorf("HTML field not copied correctly")
	}
	if clone.AllowAttachments != original.AllowAttachments {
		t.Errorf("AllowAttachments field not copied correctly")
	}

	// Check if To, CC, and BCC slices are deep copied
	if len(clone.To) != len(original.To) || clone.To[0] != original.To[0] {
		t.Errorf("To slice not deep copied correctly")
	}
	if len(clone.CC) != len(original.CC) || clone.CC[0] != original.CC[0] {
		t.Errorf("CC slice not deep copied correctly")
	}
	if len(clone.BCC) != len(original.BCC) || clone.BCC[0] != original.BCC[0] {
		t.Errorf("BCC slice not deep copied correctly")
	}

	// Check if Attachments and Inlines are deep copied
	if len(clone.Attachments) != len(original.Attachments) || clone.Attachments[0] == original.Attachments[0] {
		t.Errorf("Attachments not deep copied correctly")
	}
	if len(clone.Inlines) != len(original.Inlines) || clone.Inlines[0] == original.Inlines[0] {
		t.Errorf("Inlines not deep copied correctly")
	}

	// Check if changes to the clone do not affect the original
	clone.Subject = "Modified Subject"
	if original.Subject == clone.Subject {
		t.Errorf("Modifying clone affected the original Subject field")
	}

	clone.To[0].Name = "Modified Name"
	if original.To[0].Name == clone.To[0].Name {
		t.Errorf("Modifying clone's To slice affected the original")
	}

	clone.Attachments[0].Name = "Modified File"
	if original.Attachments[0].Name == clone.Attachments[0].Name {
		t.Errorf("Modifying clone's Attachments affected the original")
	}
}
