package azb

import (
	"github.com/gofrs/uuid/v5"
	"testing"
	"time"
)

// TestIsDialog verifies that IsDialog correctly identifies dialog events.
func TestIsDialog(t *testing.T) {
	tests := []struct {
		name  string
		event ZBType
		want  bool
	}{
		{"DialogEvent", "zurl-dialog", true},
		{"NonDialogEvent", "zurl-action", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Event: tt.event}
			if got := za.IsDialog(); got != tt.want {
				t.Errorf("ZAction.IsDialog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasZId verifies that HasZId correctly identifies non-empty Ids.
func TestHasZId(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"HasId", "123", true},
		{"EmptyId", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Id: tt.id}
			if got := za.HasZId(); got != tt.want {
				t.Errorf("ZAction.HasZId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsPaginate verifies that IsPaginate correctly identifies pagination modules.
func TestIsPaginate(t *testing.T) {
	tests := []struct {
		name string
		mod  ZBType
		want bool
	}{
		{"IsPaginate", ZMOD_PAGINATE, true},
		{"IsNotPaginate", "other-mod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Mod: tt.mod}
			if got := za.IsPaginate(); got != tt.want {
				t.Errorf("ZAction.IsPaginate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasZIdParent verifies that HasZIdParent correctly identifies non-empty IdParents.
func TestHasZIdParent(t *testing.T) {
	tests := []struct {
		name     string
		idParent string
		want     bool
	}{
		{"HasIdParent", "123-parent", true},
		{"EmptyIdParent", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{IdParent: tt.idParent}
			if got := za.HasZIdParent(); got != tt.want {
				t.Errorf("ZAction.HasZIdParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToUUIDZId verifies that ToUUIDZId correctly converts Id to UUID.
func TestToUUIDZId(t *testing.T) {
	validUUID := uuid.Must(uuid.NewV4()).String()
	tests := []struct {
		name string
		id   string
		want uuid.UUID
	}{
		{"ValidUUID", validUUID, uuid.Must(uuid.FromString(validUUID))},
		{"InvalidUUID", "invalid-uuid", uuid.Nil},
		{"EmptyUUID", "", uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Id: tt.id}
			if got := za.ToUUIDZId(); got != tt.want {
				t.Errorf("ZAction.ToUUIDZId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToIntZId verifies that ToIntZId correctly converts Id to int.
func TestToIntZId(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want int
	}{
		{"ValidInt", "123", 123},
		{"InvalidInt", "abc", 0},
		{"EmptyInt", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Id: tt.id}
			if got := za.ToIntZId(); got != tt.want {
				t.Errorf("ZAction.ToIntZId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToFloatZId verifies that ToFloatZId correctly converts Id to float64.
func TestToFloatZId(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want float64
	}{
		{"ValidFloat", "123.456", 123.456},
		{"InvalidFloat", "abc", 0},
		{"EmptyFloat", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Id: tt.id}
			if got := za.ToFloatZId(); got != tt.want {
				t.Errorf("ZAction.ToFloatZId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToTimeZId verifies that ToTimeZId correctly converts Id to time.Time.
func TestToTimeZId(t *testing.T) {
	validTime := time.Now().Format(time.RFC3339Nano)
	parsedTime, _ := time.Parse(time.RFC3339Nano, validTime)
	tests := []struct {
		name string
		id   string
		want time.Time
	}{
		{"ValidTime", validTime, parsedTime},
		{"InvalidTime", "not-a-time", time.Time{}},
		{"EmptyTime", "", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Id: tt.id}
			if got := za.ToTimeZId(); !got.Equal(tt.want) {
				t.Errorf("ZAction.ToTimeZId() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToUUIDZIdParent verifies that ToUUIDZIdParent correctly converts IdParent to UUID.
func TestToUUIDZIdParent(t *testing.T) {
	validUUID := uuid.Must(uuid.NewV4()).String()
	tests := []struct {
		name     string
		idParent string
		want     uuid.UUID
	}{
		{"ValidUUID", validUUID, uuid.Must(uuid.FromString(validUUID))},
		{"InvalidUUID", "invalid-uuid", uuid.Nil},
		{"EmptyUUID", "", uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{IdParent: tt.idParent}
			if got := za.ToUUIDZIdParent(); got != tt.want {
				t.Errorf("ZAction.ToUUIDZIdParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToIntParent verifies that ToIntParent correctly converts IdParent to int.
func TestToIntParent(t *testing.T) {
	tests := []struct {
		name     string
		idParent string
		want     int
	}{
		{"ValidInt", "123", 123},
		{"InvalidInt", "abc", 0},
		{"EmptyInt", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{IdParent: tt.idParent}
			if got := za.ToIntParent(); got != tt.want {
				t.Errorf("ZAction.ToIntParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToFloatZIdParent verifies that ToFloatZIdParent correctly converts IdParent to float64.
func TestToFloatZIdParent(t *testing.T) {
	tests := []struct {
		name     string
		idParent string
		want     float64
	}{
		{"ValidFloat", "123.456", 123.456},
		{"InvalidFloat", "abc", 0},
		{"EmptyFloat", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{IdParent: tt.idParent}
			if got := za.ToFloatZIdParent(); got != tt.want {
				t.Errorf("ZAction.ToFloatZIdParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToTimeZIdParent verifies that ToTimeZIdParent correctly converts IdParent to time.Time.
func TestToTimeZIdParent(t *testing.T) {
	validTime := time.Now().Format(time.RFC3339Nano)
	parsedTime, _ := time.Parse(time.RFC3339Nano, validTime)
	tests := []struct {
		name     string
		idParent string
		want     time.Time
	}{
		{"ValidTime", validTime, parsedTime},
		{"InvalidTime", "not-a-time", time.Time{}},
		{"EmptyTime", "", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{IdParent: tt.idParent}
			if got := za.ToTimeZIdParent(); !got.Equal(tt.want) {
				t.Errorf("ZAction.ToTimeZIdParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToUUIDZSequence verifies that ToUUIDZSequence correctly converts Sequence to UUID.
func TestToUUIDZSequence(t *testing.T) {
	validUUID := uuid.Must(uuid.NewV4()).String()
	tests := []struct {
		name     string
		sequence string
		want     uuid.UUID
	}{
		{"ValidUUID", validUUID, uuid.Must(uuid.FromString(validUUID))},
		{"InvalidUUID", "invalid-uuid", uuid.Nil},
		{"EmptyUUID", "", uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Sequence: tt.sequence}
			if got := za.ToUUIDZSequence(); got != tt.want {
				t.Errorf("ZAction.ToUUIDZSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToIntZSequence verifies that ToIntZSequence correctly converts Sequence to int.
func TestToIntZSequence(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		want     int
	}{
		{"ValidInt", "123", 123},
		{"InvalidInt", "abc", 0},
		{"EmptyInt", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Sequence: tt.sequence}
			if got := za.ToIntZSequence(); got != tt.want {
				t.Errorf("ZAction.ToIntZSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToFloatZSequence verifies that ToFloatZSequence correctly converts Sequence to float64.
func TestToFloatZSequence(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		want     float64
	}{
		{"ValidFloat", "123.456", 123.456},
		{"InvalidFloat", "abc", 0},
		{"EmptyFloat", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Sequence: tt.sequence}
			if got := za.ToFloatZSequence(); got != tt.want {
				t.Errorf("ZAction.ToFloatZSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToTimeZSequence verifies that ToTimeZSequence correctly converts Sequence to time.Time.
func TestToTimeZSequence(t *testing.T) {
	validTime := time.Now().Format(time.RFC3339Nano)
	parsedTime, _ := time.Parse(time.RFC3339Nano, validTime)
	tests := []struct {
		name     string
		sequence string
		want     time.Time
	}{
		{"ValidTime", validTime, parsedTime},
		{"InvalidTime", "not-a-time", time.Time{}},
		{"EmptyTime", "", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			za := ZAction{Sequence: tt.sequence}
			if got := za.ToTimeZSequence(); !got.Equal(tt.want) {
				t.Errorf("ZAction.ToTimeZSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}
