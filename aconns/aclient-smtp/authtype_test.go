package aclient_smtp

import (
	"testing"
)

func TestAuthType_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		at   AuthType
		want bool
	}{
		{"Empty AuthType", AuthType(""), true},
		{"Non-empty AuthType", AuthType("plain"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.IsEmpty(); got != tt.want {
				t.Errorf("AuthType.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthType_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		at   AuthType
		want AuthType
	}{
		{"Trim spaces", AuthType("  plain  "), AuthType("plain")},
		{"No spaces", AuthType("plain"), AuthType("plain")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.TrimSpace(); got != tt.want {
				t.Errorf("AuthType.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthType_String(t *testing.T) {
	tests := []struct {
		name string
		at   AuthType
		want string
	}{
		{"String representation", AuthType("plain"), "plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.String(); got != tt.want {
				t.Errorf("AuthType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthType_Matches(t *testing.T) {
	tests := []struct {
		name string
		at   AuthType
		s    string
		want bool
	}{
		{"Matches", AuthType("plain"), "plain", true},
		{"Does not match", AuthType("plain"), "none", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.Matches(tt.s); got != tt.want {
				t.Errorf("AuthType.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthType_ToStringTrimLower(t *testing.T) {
	tests := []struct {
		name string
		at   AuthType
		want string
	}{
		{"Trim and lower", AuthType("  PLAIN  "), "plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.ToStringTrimLower(); got != tt.want {
				t.Errorf("AuthType.ToStringTrimLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		at      AuthType
		wantErr bool
	}{
		{"Valid AuthType", AuthType("plain"), false},
		{"Empty AuthType", AuthType(""), true},
		{"Invalid characters", AuthType("plain!"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.at.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AuthType.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthTypes_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		want bool
	}{
		{"Empty AuthTypes", AuthTypes{}, true},
		{"Non-empty AuthTypes", AuthTypes{"plain"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.IsEmpty(); got != tt.want {
				t.Errorf("AuthTypes.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTypes_String(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		want string
	}{
		{"String representation", AuthTypes{"plain", "none"}, "plain, none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.String(); got != tt.want {
				t.Errorf("AuthTypes.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTypes_ToStringArray(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		want []string
	}{
		{"ToStringArray", AuthTypes{"plain", "none"}, []string{"plain", "none"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.ToStringArray(); !equalStringSlices(got, tt.want) {
				t.Errorf("AuthTypes.ToStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTypes_Find(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		at   AuthType
		want AuthType
	}{
		{"Find existing", AuthTypes{"plain", "none"}, AuthType("plain"), AuthType("plain")},
		{"Find non-existing", AuthTypes{"plain", "none"}, AuthType("identity"), AuthType("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.Find(tt.at); got != tt.want {
				t.Errorf("AuthTypes.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTypes_HasKey(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		s    AuthType
		want bool
	}{
		{"HasKey existing", AuthTypes{"plain", "none"}, AuthType("plain"), true},
		{"HasKey non-existing", AuthTypes{"plain", "none"}, AuthType("identity"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.HasKey(tt.s); got != tt.want {
				t.Errorf("AuthTypes.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTypes_Matches(t *testing.T) {
	tests := []struct {
		name string
		ats  AuthTypes
		s    string
		want bool
	}{
		{"Matches existing", AuthTypes{"plain", "none"}, "plain", true},
		{"Matches non-existing", AuthTypes{"plain", "none"}, "identity", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ats.Matches(tt.s); got != tt.want {
				t.Errorf("AuthTypes.Matches() = %v, want %v", got, tt.want)
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
