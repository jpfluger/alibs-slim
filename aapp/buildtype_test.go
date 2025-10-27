package aapp

import (
	"encoding/json"
	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildType_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		expected bool
	}{
		{
			name:     "empty",
			bt:       "",
			expected: true,
		},
		{
			name:     "whitespace",
			bt:       "  ",
			expected: true,
		},
		{
			name:     "non_empty",
			bt:       "release",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bt.IsEmpty(), "IsEmpty result mismatch")
		})
	}
}

func TestBuildType_String(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		expected string
	}{
		{
			name:     "release",
			bt:       "release",
			expected: "release",
		},
		{
			name:     "release_debug",
			bt:       "release:debug",
			expected: "release:debug",
		},
		{
			name:     "empty",
			bt:       "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bt.String(), "String result mismatch")
		})
	}
}

func TestBuildType_IsDebug(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		expected bool
	}{
		{
			name:     "debug_suffix",
			bt:       "release:debug",
			expected: true,
		},
		{
			name:     "no_debug_suffix",
			bt:       "release",
			expected: false,
		},
		{
			name:     "empty",
			bt:       "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bt.IsDebug(), "IsDebug result mismatch")
		})
	}
}

func TestBuildType_SetDebug(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		isDebug  bool
		expected BuildType
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "add_debug",
			bt:       "release",
			isDebug:  true,
			expected: "release:debug",
		},
		{
			name:     "remove_debug",
			bt:       "release:debug",
			isDebug:  false,
			expected: "release",
		},
		{
			name:     "no_change_debug",
			bt:       "release:debug",
			isDebug:  true,
			expected: "release:debug",
		},
		{
			name:     "no_change_no_debug",
			bt:       "release",
			isDebug:  false,
			expected: "release",
		},
		{
			name:    "empty_input",
			bt:      "",
			isDebug: true,
			wantErr: true,
			errMsg:  "build type cannot be empty",
		},
		{
			name:    "invalid_input",
			bt:      "release@invalid",
			isDebug: true,
			wantErr: true,
			errMsg:  "invalid build type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.bt.SetDebug(tc.isDebug)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, result, "SetDebug result mismatch")
		})
	}
}

func TestBuildType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		expected bool
	}{
		{
			name:     "valid_type",
			bt:       "release",
			expected: true,
		},
		{
			name:     "valid_type_with_debug",
			bt:       "release:debug",
			expected: true,
		},
		{
			name:     "invalid_type",
			bt:       "release@invalid",
			expected: false,
		},
		{
			name:     "empty_type",
			bt:       "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bt.IsValid(), "IsValid result mismatch")
		})
	}
}

func TestBuildTypes_Add(t *testing.T) {
	tests := []struct {
		name     string
		bts      BuildTypes
		bt       BuildType
		expected BuildTypes
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "add_new",
			bts:      BuildTypes{},
			bt:       "release",
			expected: BuildTypes{"release"},
		},
		{
			name:     "add_duplicate",
			bts:      BuildTypes{"release"},
			bt:       "release",
			expected: BuildTypes{"release"},
		},
		{
			name:     "add_multiple",
			bts:      BuildTypes{"release"},
			bt:       "debug",
			expected: BuildTypes{"release", "debug"},
		},
		{
			name:    "add_empty",
			bts:     BuildTypes{},
			bt:      "",
			wantErr: true,
			errMsg:  "cannot add empty build type",
		},
		{
			name:    "add_invalid",
			bts:     BuildTypes{},
			bt:      "release@invalid",
			wantErr: true,
			errMsg:  "invalid build type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.bts.Add(tc.bt)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, tc.bts, "Add result mismatch")
		})
	}
}

func TestBuildTypes_Remove(t *testing.T) {
	tests := []struct {
		name     string
		bts      BuildTypes
		bt       BuildType
		expected BuildTypes
	}{
		{
			name:     "remove_existing",
			bts:      BuildTypes{"release", "debug"},
			bt:       "release",
			expected: BuildTypes{"debug"},
		},
		{
			name:     "remove_non_existing",
			bts:      BuildTypes{"release", "debug"},
			bt:       "profile",
			expected: BuildTypes{"release", "debug"},
		},
		{
			name:     "remove_from_empty",
			bts:      BuildTypes{},
			bt:       "release",
			expected: BuildTypes{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.bts.Remove(tc.bt)
			assert.Equal(t, tc.expected, tc.bts, "Remove result mismatch")
		})
	}
}

func TestBuildTypes_Contains(t *testing.T) {
	tests := []struct {
		name     string
		bts      BuildTypes
		bt       BuildType
		expected bool
	}{
		{
			name:     "contains_existing",
			bts:      BuildTypes{"release", "debug"},
			bt:       "release",
			expected: true,
		},
		{
			name:     "does_not_contain",
			bts:      BuildTypes{"release", "debug"},
			bt:       "profile",
			expected: false,
		},
		{
			name:     "empty_slice",
			bts:      BuildTypes{},
			bt:       "release",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bts.Contains(tc.bt), "Contains result mismatch")
		})
	}
}

func TestBuildTypes_IsKnownType(t *testing.T) {
	tests := []struct {
		name     string
		bt       BuildType
		expected bool
	}{
		{
			name:     "known_type",
			bt:       BUILDTYPE_RELEASE,
			expected: true,
		},
		{
			name:     "unknown_type",
			bt:       "custom",
			expected: false,
		},
		{
			name:     "known_type_with_debug",
			bt:       "release:debug",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bts := BuildTypes{}
			assert.Equal(t, tc.expected, bts.IsKnownType(tc.bt), "IsKnownType result mismatch")
		})
	}
}

func TestBuildTypes_SelectPreferredDefault(t *testing.T) {
	tests := []struct {
		name     string
		bts      BuildTypes
		expected BuildType
	}{
		{
			name:     "prefer_release",
			bts:      BuildTypes{"debug", "profile", "release"},
			expected: BUILDTYPE_RELEASE,
		},
		{
			name:     "prefer_profile",
			bts:      BuildTypes{"debug", "profile"},
			expected: BUILDTYPE_PROFILE,
		},
		{
			name:     "prefer_debug",
			bts:      BuildTypes{"debug"},
			expected: BUILDTYPE_DEBUG,
		},
		{
			name:     "empty_slice",
			bts:      BuildTypes{},
			expected: BUILDTYPE_DEBUG,
		},
		{
			name:     "unknown_types",
			bts:      BuildTypes{"custom"},
			expected: BUILDTYPE_DEBUG,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.bts.SelectPreferredDefault(), "SelectPreferredDefault result mismatch")
		})
	}
}

func TestBuildType_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"type": "release",
		"types": ["release", "debug", "profile:debug"],
		"invalid": "release@invalid"
	}`)

	var m map[string]json.RawMessage
	err := json.Unmarshal(jsonData, &m)
	assert.NoError(t, err, "Failed to unmarshal JSON map")

	tests := []struct {
		name     string
		target   interface{}
		expected interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "single_build_type",
			target:   new(BuildType),
			expected: BuildType("release"),
		},
		{
			name:     "build_types_slice",
			target:   new(BuildTypes),
			expected: BuildTypes{"release", "debug", "profile:debug"},
		},
		{
			name:    "invalid_build_type",
			target:  new(BuildType),
			wantErr: true,
			errMsg:  "invalid build type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var key string
			switch tc.name {
			case "single_build_type":
				key = "type"
			case "build_types_slice":
				key = "types"
			case "invalid_build_type":
				key = "invalid"
			}

			err = ajson.UnmarshalRawMessage(m[key], tc.target)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			switch target := tc.target.(type) {
			case *BuildType:
				assert.Equal(t, tc.expected, *target, "Unmarshaled result mismatch")
			case *BuildTypes:
				assert.Equal(t, tc.expected, *target, "Unmarshaled result mismatch")
			default:
				t.Fatalf("Unexpected target type: %T", tc.target)
			}
		})
	}
}
