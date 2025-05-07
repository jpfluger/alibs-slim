package aconns

import (
	"testing"

	"github.com/jpfluger/alibs-slim/auuids"
)

func TestConnTenantInfo_Validate(t *testing.T) {
	valid := ConnTenantInfo{
		Region:   "us-west",
		TenantId: auuids.NewUUID(),
		Priority: 0,
		Label:    "US-West Primary",
	}

	if err := valid.Validate(); err != nil {
		t.Errorf("Expected valid ConnTenantInfo, got error: %v", err)
	}

	tests := []struct {
		name    string
		input   ConnTenantInfo
		wantErr bool
	}{
		{
			name: "missing tenant ID",
			input: ConnTenantInfo{
				Region:   "us-east",
				TenantId: auuids.UUID{}, // zero UUID
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "missing region",
			input: ConnTenantInfo{
				Region:   "  ",
				TenantId: auuids.NewUUID(),
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "negative priority",
			input: ConnTenantInfo{
				Region:   "eu-central",
				TenantId: auuids.NewUUID(),
				Priority: -2,
			},
			wantErr: true,
		},
		{
			name: "valid full case",
			input: ConnTenantInfo{
				Region:   "ap-south",
				TenantId: auuids.NewUUID(),
				Priority: 2,
				Label:    "India Cluster",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnTenantInfos_Validate(t *testing.T) {
	valid := ConnTenantInfos{
		{
			Region:   "us-west",
			TenantId: auuids.NewUUID(),
			Priority: 1,
		},
		{
			Region:   "us-east",
			TenantId: auuids.NewUUID(),
			Priority: 0,
		},
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("Expected no error for valid ConnTenantInfos, got: %v", err)
	}

	invalid := ConnTenantInfos{
		{
			Region:   "us-west",
			TenantId: auuids.NewUUID(),
			Priority: 1,
		},
		{
			Region:   "",
			TenantId: auuids.NewUUID(),
			Priority: 0,
		},
	}
	if err := invalid.Validate(); err == nil {
		t.Error("Expected error for invalid ConnTenantInfos, got none")
	}
}

func TestConnTenantInfos_HasTenant(t *testing.T) {
	id1 := auuids.NewUUID()
	id2 := auuids.NewUUID()
	infos := ConnTenantInfos{
		{Region: "us-west", TenantId: id1, Priority: 0},
	}

	if !infos.HasTenant(id1) {
		t.Errorf("Expected HasTenant(%v) to return true", id1.String())
	}
	if infos.HasTenant(id2) {
		t.Errorf("Expected HasTenant(%v) to return false", id2.String())
	}
}

func TestConnTenantInfos_GetByTenantId(t *testing.T) {
	id := auuids.NewUUID()
	infos := ConnTenantInfos{
		{Region: "us-west", TenantId: id, Priority: 0},
	}

	got, found := infos.GetByTenantId(id)
	if !found {
		t.Errorf("Expected GetByTenantId to find %v", id.String())
	}
	if got.TenantId != id {
		t.Errorf("Expected TenantId to match, got %v", got.TenantId.String())
	}
}

func TestConnTenantInfos_FilterByRegion(t *testing.T) {
	infos := ConnTenantInfos{
		{Region: "us-west", TenantId: auuids.NewUUID(), Priority: 1},
		{Region: "US-WEST", TenantId: auuids.NewUUID(), Priority: 2},
		{Region: "eu-central", TenantId: auuids.NewUUID(), Priority: 3},
	}

	filtered := infos.FilterByRegion("us-west")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 results for region 'us-west', got %d", len(filtered))
	}

	filteredNone := infos.FilterByRegion("asia-east")
	if len(filteredNone) != 0 {
		t.Errorf("Expected 0 results for region 'asia-east', got %d", len(filteredNone))
	}
}
