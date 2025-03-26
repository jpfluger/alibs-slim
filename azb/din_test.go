package azb

import (
	"testing"
)

// TestDINValidate checks if the Validate method correctly updates the ZAction's page limit.
func TestDINValidate(t *testing.T) {
	din := &DIN{
		ZAction: ZAction{
			PageLimit: 10,
			PageOn:    1,
		},
	}

	// Intentionally setting an incorrect PageLimit to trigger the validation logic.
	din.ZAction.PageLimit = 20

	if err := din.Validate(); err != nil {
		t.Errorf("Validate() returned an error: %v", err)
	}

	if din.ZAction.PageLimit != 20 {
		t.Errorf("Validate() did not update PageLimit correctly, got: %d, want: %d", din.ZAction.PageLimit, 20)
	}

	if din.ZAction.PageOn != 1 {
		t.Errorf("Validate() did not reset PageOn correctly, got: %d, want: %d", din.ZAction.PageOn, 1)
	}
}

// TestDINNewPaginate checks if the NewPaginate method correctly creates a new Paginate instance.
func TestDINNewPaginate(t *testing.T) {
	din := &DIN{
		ZAction: ZAction{
			PageOn: 2,
		},
	}

	totalItems, itemsPerPage := 50, 10
	paginate := din.NewPaginate(totalItems, itemsPerPage)

	if paginate.CurrentPage != din.ZAction.PageOn {
		t.Errorf("NewPaginate() CurrentPage = %d, want %d", paginate.CurrentPage, din.ZAction.PageOn)
	}

	if paginate.TotalItems != totalItems {
		t.Errorf("NewPaginate() TotalItems = %d, want %d", paginate.TotalItems, totalItems)
	}

	if paginate.ItemsPerPage != itemsPerPage {
		t.Errorf("NewPaginate() ItemsPerPage = %d, want %d", paginate.ItemsPerPage, itemsPerPage)
	}

	expectedTotalPages := 5 // 50 items, 10 per page, should result in 5 total pages.
	if paginate.TotalPages != expectedTotalPages {
		t.Errorf("NewPaginate() TotalPages = %d, want %d", paginate.TotalPages, expectedTotalPages)
	}
}
