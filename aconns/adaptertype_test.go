package aconns

import (
	"testing"
)

func TestAdapterType_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		at   AdapterType
		want bool
	}{
		{"Empty AdapterType", AdapterType(""), true},
		{"Non-empty AdapterType", AdapterType("adapter1"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.IsEmpty(); got != tt.want {
				t.Errorf("AdapterType.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterType_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		at   AdapterType
		want AdapterType
	}{
		{"Trim spaces", AdapterType(" adapter1 "), AdapterType("adapter1")},
		{"No spaces", AdapterType("adapter1"), AdapterType("adapter1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.TrimSpace(); got != tt.want {
				t.Errorf("AdapterType.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterType_String(t *testing.T) {
	at := AdapterType("adapter1")
	want := "adapter1"
	if got := at.String(); got != want {
		t.Errorf("AdapterType.String() = %v, want %v", got, want)
	}
}

func TestAdapterType_Matches(t *testing.T) {
	tests := []struct {
		name string
		at   AdapterType
		s    string
		want bool
	}{
		{"Matches", AdapterType("adapter1"), "adapter1", true},
		{"Does not match", AdapterType("adapter1"), "adapter2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.Matches(tt.s); got != tt.want {
				t.Errorf("AdapterType.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterType_ToStringTrimLower(t *testing.T) {
	tests := []struct {
		name string
		at   AdapterType
		want string
	}{
		{"Trim and lower", AdapterType(" Adapter1 "), "adapter1"},
		{"Already trimmed and lower", AdapterType("adapter1"), "adapter1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.ToStringTrimLower(); got != tt.want {
				t.Errorf("AdapterType.ToStringTrimLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		at      AdapterType
		wantErr bool
	}{
		{"Valid AdapterType", AdapterType("adapter1"), false},
		{"Empty AdapterType", AdapterType(""), true},
		{"Invalid characters", AdapterType("adapter!"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.at.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AdapterType.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdapterTypes_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ats  AdapterTypes
		want bool
	}{
		{"Empty AdapterTypes", AdapterTypes{}, true},
		{"Non-empty AdapterTypes", AdapterTypes{AdapterType("adapter1")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.IsEmpty(); got != tt.want {
				t.Errorf("AdapterTypes.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterTypes_String(t *testing.T) {
	tests := []struct {
		name string
		ats  AdapterTypes
		want string
	}{
		{"Empty AdapterTypes", AdapterTypes{}, ""},
		{"Single AdapterType", AdapterTypes{AdapterType("type1")}, "type1"},
		{"Multiple AdapterTypes", AdapterTypes{AdapterType("type1"), AdapterType("type2")}, "type1, type2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.String(); got != tt.want {
				t.Errorf("AdapterTypes.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterTypes_ToStringArray(t *testing.T) {
	ats := AdapterTypes{AdapterType("adapter1"), AdapterType("adapter2")}
	want := []string{"adapter1", "adapter2"}
	if got := ats.ToStringArray(); !equalStringSlices(got, want) {
		t.Errorf("AdapterTypes.ToStringArray() = %v, want %v", got, want)
	}
}

func TestAdapterTypes_Find(t *testing.T) {
	ats := AdapterTypes{AdapterType("adapter1"), AdapterType("adapter2")}
	tests := []struct {
		name string
		at   AdapterType
		want AdapterType
	}{
		{"Find existing", AdapterType("adapter1"), AdapterType("adapter1")},
		{"Find non-existing", AdapterType("adapter3"), AdapterType("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ats.Find(tt.at); got != tt.want {
				t.Errorf("AdapterTypes.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterTypes_HasKey(t *testing.T) {
	ats := AdapterTypes{AdapterType("adapter1"), AdapterType("adapter2")}
	tests := []struct {
		name string
		at   AdapterType
		want bool
	}{
		{"Has key", AdapterType("adapter1"), true},
		{"Does not have key", AdapterType("adapter3"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ats.HasKey(tt.at); got != tt.want {
				t.Errorf("AdapterTypes.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterTypes_Matches(t *testing.T) {
	ats := AdapterTypes{AdapterType("adapter1"), AdapterType("adapter2")}
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"Matches", "adapter1", true},
		{"Does not match", "adapter3", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ats.Matches(tt.s); got != tt.want {
				t.Errorf("AdapterTypes.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
