package aclient_smtp

import (
	"testing"
)

func TestAttachmentKey_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ak   AttachmentKey
		want bool
	}{
		{"Empty key", ATTACHMENTKEY_NONE, true},
		{"Non-empty key", ATTACHMENTKEY_FILE, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ak.IsEmpty(); got != tt.want {
				t.Errorf("AttachmentKey.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ak      AttachmentKey
		wantErr bool
	}{
		{"Valid key", AttachmentKey("file:example.txt"), false},
		{"Empty key", ATTACHMENTKEY_NONE, false},
		{"Invalid key - no colon", AttachmentKey("fileexample.txt"), true},
		{"Invalid key - empty parts", AttachmentKey("file:"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ak.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AttachmentKey.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttachmentKey_IsFile(t *testing.T) {
	tests := []struct {
		name string
		ak   AttachmentKey
		want bool
	}{
		{"File key", AttachmentKey("file:example.txt"), true},
		{"Non-file key", AttachmentKey("id:12345"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ak.IsFile(); got != tt.want {
				t.Errorf("AttachmentKey.IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKey_IsId(t *testing.T) {
	tests := []struct {
		name string
		ak   AttachmentKey
		want bool
	}{
		{"ID key", AttachmentKey("id:12345"), true},
		{"Non-ID key", AttachmentKey("file:example.txt"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ak.IsId(); got != tt.want {
				t.Errorf("AttachmentKey.IsId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKey_GetParts(t *testing.T) {
	tests := []struct {
		name    string
		ak      AttachmentKey
		wantKey AttachmentKey
		want    string
		wantErr bool
	}{
		{"Valid key", AttachmentKey("file:example.txt"), ATTACHMENTKEY_FILE, "example.txt", false},
		{"Empty key", ATTACHMENTKEY_NONE, ATTACHMENTKEY_NONE, "", false},
		{"Invalid key - no colon", AttachmentKey("fileexample.txt"), "", "", true},
		{"Invalid key - empty parts", AttachmentKey("file:"), "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, got, err := tt.ak.GetParts()
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachmentKey.GetParts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("AttachmentKey.GetParts() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if got != tt.want {
				t.Errorf("AttachmentKey.GetParts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKeys_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		aks  AttachmentKeys
		want bool
	}{
		{"Empty keys", AttachmentKeys{}, true},
		{"Non-empty keys", AttachmentKeys{ATTACHMENTKEY_FILE}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aks.IsEmpty(); got != tt.want {
				t.Errorf("AttachmentKeys.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKeys_HasKey(t *testing.T) {
	tests := []struct {
		name string
		aks  AttachmentKeys
		key  AttachmentKey
		want bool
	}{
		{"Has key", AttachmentKeys{ATTACHMENTKEY_FILE}, ATTACHMENTKEY_FILE, true},
		{"Does not have key", AttachmentKeys{ATTACHMENTKEY_FILE}, ATTACHMENTKEY_ID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aks.HasKey(tt.key); got != tt.want {
				t.Errorf("AttachmentKeys.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentKeys_Matches(t *testing.T) {
	tests := []struct {
		name string
		aks  AttachmentKeys
		s    string
		want bool
	}{
		{"Matches key", AttachmentKeys{ATTACHMENTKEY_FILE}, "file", true},
		{"Does not match key", AttachmentKeys{ATTACHMENTKEY_FILE}, "id", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aks.Matches(tt.s); got != tt.want {
				t.Errorf("AttachmentKeys.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
