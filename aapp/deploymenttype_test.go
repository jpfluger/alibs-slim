package aapp

import (
	"encoding/json"
	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeploymentType_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		expected bool
	}{
		{
			name:     "empty",
			dt:       "",
			expected: true,
		},
		{
			name:     "whitespace",
			dt:       "  ",
			expected: true,
		},
		{
			name:     "non_empty",
			dt:       "dev",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dt.IsEmpty(), "IsEmpty result mismatch")
		})
	}
}

func TestDeploymentType_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		expected string
	}{
		{
			name:     "dev",
			dt:       "dev",
			expected: "dev",
		},
		{
			name:     "dev_demo",
			dt:       "dev:demo",
			expected: "dev:demo",
		},
		{
			name:     "empty",
			dt:       "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dt.String(), "String result mismatch")
		})
	}
}

func TestDeploymentType_IsDemo(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		expected bool
	}{
		{
			name:     "demo_suffix",
			dt:       "dev:demo",
			expected: true,
		},
		{
			name:     "no_demo_suffix",
			dt:       "dev",
			expected: false,
		},
		{
			name:     "empty",
			dt:       "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dt.IsDemo(), "IsDemo result mismatch")
		})
	}
}

func TestDeploymentType_SetDemo(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		isDemo   bool
		expected DeploymentType
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "add_demo",
			dt:       "dev",
			isDemo:   true,
			expected: "dev:demo",
		},
		{
			name:     "remove_demo",
			dt:       "dev:demo",
			isDemo:   false,
			expected: "dev",
		},
		{
			name:     "no_change_demo",
			dt:       "dev:demo",
			isDemo:   true,
			expected: "dev:demo",
		},
		{
			name:     "no_change_no_demo",
			dt:       "dev",
			isDemo:   false,
			expected: "dev",
		},
		{
			name:    "empty_input",
			dt:      "",
			isDemo:  true,
			wantErr: true,
			errMsg:  "deployment type cannot be empty",
		},
		{
			name:    "invalid_input",
			dt:      "dev@invalid",
			isDemo:  true,
			wantErr: true,
			errMsg:  "invalid deployment type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.dt.SetDemo(tc.isDemo)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, result, "SetDemo result mismatch")
		})
	}
}

func TestDeploymentType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		expected bool
	}{
		{
			name:     "valid_type",
			dt:       "dev",
			expected: true,
		},
		{
			name:     "valid_type_with_demo",
			dt:       "dev:demo",
			expected: true,
		},
		{
			name:     "invalid_type",
			dt:       "dev@invalid",
			expected: false,
		},
		{
			name:     "empty_type",
			dt:       "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dt.IsValid(), "IsValid result mismatch")
		})
	}
}

func TestDeploymentTypes_Add(t *testing.T) {
	tests := []struct {
		name     string
		dts      DeploymentTypes
		dt       DeploymentType
		expected DeploymentTypes
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "add_new",
			dts:      DeploymentTypes{},
			dt:       "dev",
			expected: DeploymentTypes{"dev"},
		},
		{
			name:     "add_duplicate",
			dts:      DeploymentTypes{"dev"},
			dt:       "dev",
			expected: DeploymentTypes{"dev"},
		},
		{
			name:     "add_multiple",
			dts:      DeploymentTypes{"dev"},
			dt:       "qa",
			expected: DeploymentTypes{"dev", "qa"},
		},
		{
			name:    "add_empty",
			dts:     DeploymentTypes{},
			dt:      "",
			wantErr: true,
			errMsg:  "cannot add empty deployment type",
		},
		{
			name:    "add_invalid",
			dts:     DeploymentTypes{},
			dt:      "dev@invalid",
			wantErr: true,
			errMsg:  "invalid deployment type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.dts.Add(tc.dt)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, tc.dts, "Add result mismatch")
		})
	}
}

func TestDeploymentTypes_Remove(t *testing.T) {
	tests := []struct {
		name     string
		dts      DeploymentTypes
		dt       DeploymentType
		expected DeploymentTypes
	}{
		{
			name:     "remove_existing",
			dts:      DeploymentTypes{"dev", "qa"},
			dt:       "dev",
			expected: DeploymentTypes{"qa"},
		},
		{
			name:     "remove_non_existing",
			dts:      DeploymentTypes{"dev", "qa"},
			dt:       "prod",
			expected: DeploymentTypes{"dev", "qa"},
		},
		{
			name:     "remove_from_empty",
			dts:      DeploymentTypes{},
			dt:       "dev",
			expected: DeploymentTypes{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.dts.Remove(tc.dt)
			assert.Equal(t, tc.expected, tc.dts, "Remove result mismatch")
		})
	}
}

func TestDeploymentTypes_Contains(t *testing.T) {
	tests := []struct {
		name     string
		dts      DeploymentTypes
		dt       DeploymentType
		expected bool
	}{
		{
			name:     "contains_existing",
			dts:      DeploymentTypes{"dev", "qa"},
			dt:       "dev",
			expected: true,
		},
		{
			name:     "does_not_contain",
			dts:      DeploymentTypes{"dev", "qa"},
			dt:       "prod",
			expected: false,
		},
		{
			name:     "empty_slice",
			dts:      DeploymentTypes{},
			dt:       "dev",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dts.Contains(tc.dt), "Contains result mismatch")
		})
	}
}

func TestDeploymentTypes_IsKnownType(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeploymentType
		expected bool
	}{
		{
			name:     "known_type",
			dt:       DEPLOYMENTTYPE_DEV,
			expected: true,
		},
		{
			name:     "unknown_type",
			dt:       "staging",
			expected: false,
		},
		{
			name:     "known_type_with_demo",
			dt:       "dev:demo",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dts := DeploymentTypes{}
			assert.Equal(t, tc.expected, dts.IsKnownType(tc.dt), "IsKnownType result mismatch")
		})
	}
}

func TestDeploymentTypes_SelectPreferredDefault(t *testing.T) {
	tests := []struct {
		name     string
		dts      DeploymentTypes
		expected DeploymentType
	}{
		{
			name:     "prefer_prod",
			dts:      DeploymentTypes{"dev", "qa", "prod"},
			expected: DEPLOYMENTTYPE_PROD,
		},
		{
			name:     "prefer_qa",
			dts:      DeploymentTypes{"dev", "qa"},
			expected: DEPLOYMENTTYPE_QA,
		},
		{
			name:     "prefer_dev",
			dts:      DeploymentTypes{"dev"},
			expected: DEPLOYMENTTYPE_DEV,
		},
		{
			name:     "prefer_local",
			dts:      DeploymentTypes{"local"},
			expected: DEPLOYMENTTYPE_LOCAL,
		},
		{
			name:     "empty_slice",
			dts:      DeploymentTypes{},
			expected: DEPLOYMENTTYPE_LOCAL,
		},
		{
			name:     "unknown_types",
			dts:      DeploymentTypes{"staging"},
			expected: DEPLOYMENTTYPE_LOCAL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.dts.SelectPreferredDefault(), "SelectPreferredDefault result mismatch")
		})
	}
}

func TestDeploymentType_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"type": "dev",
		"types": ["dev", "qa", "prod:demo"],
		"invalid": "dev@invalid"
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
			name:     "single_deployment_type",
			target:   new(DeploymentType),
			expected: DeploymentType("dev"),
		},
		{
			name:     "deployment_types_slice",
			target:   new(DeploymentTypes),
			expected: DeploymentTypes{"dev", "qa", "prod:demo"},
		},
		{
			name:    "invalid_deployment_type",
			target:  new(DeploymentType),
			wantErr: true,
			errMsg:  "invalid deployment type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var key string
			switch tc.name {
			case "single_deployment_type":
				key = "type"
			case "deployment_types_slice":
				key = "types"
			case "invalid_deployment_type":
				key = "invalid"
			}

			err = ajson.UnmarshalRawMessage(m[key], tc.target)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			switch target := tc.target.(type) {
			case *DeploymentType:
				assert.Equal(t, tc.expected, *target, "Unmarshaled result mismatch")
			case *DeploymentTypes:
				assert.Equal(t, tc.expected, *target, "Unmarshaled result mismatch")
			default:
				t.Fatalf("Unexpected target type: %T", tc.target)
			}
		})
	}
}
