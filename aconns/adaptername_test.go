package aconns

import (
	"testing"
)

func TestAdapterName_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		an   AdapterName
		want bool
	}{
		{"Empty AdapterName", AdapterName(""), true},
		{"Non-empty AdapterName", AdapterName("adapter1"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.an.IsEmpty(); got != tt.want {
				t.Errorf("AdapterName.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterName_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		an   AdapterName
		want AdapterName
	}{
		{"Trim spaces", AdapterName(" adapter1 "), AdapterName("adapter1")},
		{"No spaces", AdapterName("adapter1"), AdapterName("adapter1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.an.TrimSpace(); got != tt.want {
				t.Errorf("AdapterName.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterName_String(t *testing.T) {
	an := AdapterName("adapter1")
	want := "adapter1"
	if got := an.String(); got != want {
		t.Errorf("AdapterName.String() = %v, want %v", got, want)
	}
}

func TestAdapterName_Matches(t *testing.T) {
	tests := []struct {
		name string
		an   AdapterName
		s    string
		want bool
	}{
		{"Matches", AdapterName("adapter1"), "adapter1", true},
		{"Does not match", AdapterName("adapter1"), "adapter2", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.an.Matches(tt.s); got != tt.want {
				t.Errorf("AdapterName.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterNames_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ans  AdapterNames
		want bool
	}{
		{"Empty AdapterNames", AdapterNames{}, true},
		{"Non-empty AdapterNames", AdapterNames{AdapterName("adapter1")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ans.IsEmpty(); got != tt.want {
				t.Errorf("AdapterNames.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterNames_String(t *testing.T) {
	tests := []struct {
		name string
		ans  AdapterNames
		want string
	}{
		{"Empty AdapterNames", AdapterNames{}, ""},
		{"Single AdapterName", AdapterNames{AdapterName("adapter1")}, "adapter1"},
		{"Multiple AdapterNames", AdapterNames{AdapterName("adapter1"), AdapterName("adapter2")}, "adapter1, adapter2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ans.String(); got != tt.want {
				t.Errorf("AdapterNames.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterNames_ToStringArray(t *testing.T) {
	ans := AdapterNames{AdapterName("adapter1"), AdapterName("adapter2")}
	want := []string{"adapter1", "adapter2"}
	if got := ans.ToStringArray(); !equalStringSlices(got, want) {
		t.Errorf("AdapterNames.ToStringArray() = %v, want %v", got, want)
	}
}

func TestAdapterNames_Find(t *testing.T) {
	ans := AdapterNames{AdapterName("adapter1"), AdapterName("adapter2")}
	tests := []struct {
		name string
		an   AdapterName
		want AdapterName
	}{
		{"Find existing", AdapterName("adapter1"), AdapterName("adapter1")},
		{"Find non-existing", AdapterName("adapter3"), AdapterName("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ans.Find(tt.an); got != tt.want {
				t.Errorf("AdapterNames.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterNames_HasKey(t *testing.T) {
	ans := AdapterNames{AdapterName("adapter1"), AdapterName("adapter2")}
	tests := []struct {
		name string
		an   AdapterName
		want bool
	}{
		{"Has key", AdapterName("adapter1"), true},
		{"Does not have key", AdapterName("adapter3"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ans.HasKey(tt.an); got != tt.want {
				t.Errorf("AdapterNames.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapterNames_Matches(t *testing.T) {
	ans := AdapterNames{AdapterName("adapter1"), AdapterName("adapter2")}
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
			if got := ans.Matches(tt.s); got != tt.want {
				t.Errorf("AdapterNames.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices of strings
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
