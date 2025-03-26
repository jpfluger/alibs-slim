package aconns

import (
	"testing"
)

func TestAdapter_GetType(t *testing.T) {
	tests := []struct {
		name    string
		adapter Adapter
		want    AdapterType
	}{
		{"GetType with non-empty type", Adapter{Type: AdapterType("type1")}, AdapterType("type1")},
		{"GetType with empty type", Adapter{Type: AdapterType("")}, AdapterType("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.adapter.GetType(); got != tt.want {
				t.Errorf("Adapter.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapter_GetName(t *testing.T) {
	tests := []struct {
		name    string
		adapter Adapter
		want    AdapterName
	}{
		{"GetName with non-empty name", Adapter{Name: AdapterName("name1")}, AdapterName("name1")},
		{"GetName with empty name", Adapter{Name: AdapterName("")}, AdapterName("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.adapter.GetName(); got != tt.want {
				t.Errorf("Adapter.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdapter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		adapter Adapter
		wantErr bool
	}{
		{"Valid Adapter", Adapter{Type: AdapterType("type1"), Name: AdapterName("name1"), Host: "localhost"}, false},
		{"Empty Type", Adapter{Type: AdapterType(""), Name: AdapterName("name1"), Host: "localhost"}, true},
		{"Empty Name", Adapter{Type: AdapterType("type1"), Name: AdapterName(""), Host: "localhost"}, true},
		{"Empty Type and Name", Adapter{Type: AdapterType(""), Name: AdapterName(""), Host: "localhost"}, true},
		{"Empty Host", Adapter{Type: AdapterType("type1"), Name: AdapterName("name1"), Host: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.adapter.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Adapter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
