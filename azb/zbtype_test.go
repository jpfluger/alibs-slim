package azb

import (
	"github.com/jpfluger/alibs-slim/ajson"
	"reflect"
	"testing"
)

// TestIsEmpty verifies that IsEmpty correctly identifies empty ZBType values.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   bool
	}{
		{"Empty", "", true},
		{"Whitespace", " ", true},
		{"NotEmpty", "azp", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.IsEmpty(); got != tt.want {
				t.Errorf("ZBType.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTrimSpace verifies that TrimSpace correctly trims whitespace from ZBType values.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   ZBType
	}{
		{"LeadingAndTrailingSpaces", " azp ", "azp"},
		{"NoSpaces", "azp", "azp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.TrimSpace(); got != tt.want {
				t.Errorf("ZBType.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestString verifies that String correctly converts ZBType to string.
func TestString(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   string
	}{
		{"SimpleString", "azp", "azp"},
		{"EmptyString", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.String(); got != tt.want {
				t.Errorf("ZBType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToStringTrimLower verifies that ToStringTrimLower correctly trims and lowers ZBType values.
func TestToStringTrimLower(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   string
	}{
		{"UpperCase", "AZP", "azp"},
		{"MixedCase", "AzP", "azp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.ToStringTrimLower(); got != tt.want {
				t.Errorf("ZBType.ToStringTrimLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToJsonKey verifies that ToJsonKey correctly converts ZBType to a JsonKey.
func TestToJsonKey(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   ajson.JsonKey
	}{
		{"SimpleKey", "azp", "azp"},
		{"ComplexKey", "AzP", "azp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.ToJsonKey(); got != tt.want {
				t.Errorf("ZBType.ToJsonKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasMatch verifies that HasMatch correctly identifies matching ZBType values.
func TestHasMatch(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		match  ZBType
		want   bool
	}{
		{"Match", "azp", "azp", true},
		{"NoMatch", "azp", "nope", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.HasMatch(tt.match); got != tt.want {
				t.Errorf("ZBType.HasMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMatchesOne verifies that MatchesOne correctly identifies if ZBType matches any of the provided ZBTypes.
func TestMatchesOne(t *testing.T) {
	tests := []struct {
		name    string
		zpType  ZBType
		matches []ZBType
		want    bool
	}{
		{"OneMatch", "azp", []ZBType{"azp", "nope"}, true},
		{"NoMatches", "azp", []ZBType{"nope", "nah"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.MatchesOne(tt.matches...); got != tt.want {
				t.Errorf("ZBType.MatchesOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasPrefix verifies that HasPrefix correctly identifies if ZBType has the specified prefix.
func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		prefix ZBType
		want   bool
	}{
		{"HasPrefix", "azpType", "azp", true},
		{"NoPrefix", "azpType", "nope", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.HasPrefix(tt.prefix); got != tt.want {
				t.Errorf("ZBType.HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasSuffix verifies that HasSuffix correctly identifies if ZBType has the specified suffix.
func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		suffix ZBType
		want   bool
	}{
		{"HasSuffix", "Typeazp", "azp", true},
		{"NoSuffix", "Typeazp", "nope", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.HasSuffix(tt.suffix); got != tt.want {
				t.Errorf("ZBType.HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetLeaf verifies that GetLeaf correctly extracts the last element in a JSON key path.
func TestGetLeaf(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   string
	}{
		{"SimpleLeaf", "azp", "azp"},
		{"ComplexLeaf", "azp.leaf", "leaf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.GetLeaf(); got != tt.want {
				t.Errorf("ZBType.GetLeaf() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetLeafToUpper verifies that GetLeafToUpper correctly converts the last element in a JSON key path to upper case.
func TestGetLeafToUpper(t *testing.T) {
	tests := []struct {
		name   string
		zpType ZBType
		want   string
	}{
		{"SimpleLeafUpper", "azp", "AZP"},
		{"ComplexLeafUpper", "azp.leaf", "LEAF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zpType.GetLeafToUpper(); got != tt.want {
				t.Errorf("ZBType.GetLeafToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasValues verifies that HasValues correctly identifies if ZBTypes slice has any values.
func TestHasValues(t *testing.T) {
	tests := []struct {
		name string
		rts  ZBTypes
		want bool
	}{
		{"HasValues", []ZBType{"azp"}, true},
		{"NoValues", []ZBType{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.HasValues(); got != tt.want {
				t.Errorf("ZBTypes.HasValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClone checks the Clone method.
func TestClone(t *testing.T) {
	tests := []struct {
		name string
		rts  ZBTypes
		want ZBTypes
	}{
		{"CloneNonEmpty", []ZBType{"azp", "nope"}, []ZBType{"azp", "nope"}},
		{"CloneEmpty", []ZBType{}, []ZBType{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZBTypes.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToArrStrings checks the ToArrStrings method.
func TestToArrStrings(t *testing.T) {
	tests := []struct {
		name string
		rts  ZBTypes
		want []string
	}{
		{"NonEmptyToStrings", []ZBType{"azp", "nope"}, []string{"azp", "nope"}},
		{"EmptyToStrings", []ZBType{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.ToArrStrings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZBTypes.ToArrStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIncludeIfInTargets checks the IncludeIfInTargets method.
func TestIncludeIfInTargets(t *testing.T) {
	tests := []struct {
		name    string
		rts     ZBTypes
		targets []ZBType
		want    ZBTypes
	}{
		{"IncludeMatching", []ZBType{"azp", "nope"}, []ZBType{"azp"}, []ZBType{"azp"}},
		{"ExcludeNonMatching", []ZBType{"azp", "nope"}, []ZBType{"nah"}, []ZBType{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.IncludeIfInTargets(tt.targets...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZBTypes.IncludeIfInTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClean checks the Clean method.
func TestClean(t *testing.T) {
	tests := []struct {
		name string
		rts  ZBTypes
		want ZBTypes
	}{
		{"CleanNonEmpty", []ZBType{"azp", "", "nope"}, []ZBType{"azp", "nope"}},
		{"CleanAllEmpty", []ZBType{"", ""}, []ZBType{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.Clean(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZBTypes.Clean() = %v, want %v", got, tt.want)
			}
		})
	}
}
