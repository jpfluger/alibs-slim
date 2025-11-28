package azb

import (
	"reflect"
	"testing"
)

// TestDINValidate verifies default handling and unchanged valid values.
func TestDINValidate(t *testing.T) {
	tests := []struct {
		name      string
		initial   ZAction
		wantOn    int
		wantLimit int
	}{
		{
			name:      "Defaults for invalid values",
			initial:   ZAction{PageOn: -1, PageLimit: 0},
			wantOn:    1,
			wantLimit: 25,
		},
		{
			name:      "Unchanged valid values",
			initial:   ZAction{PageOn: 2, PageLimit: 20},
			wantOn:    2,
			wantLimit: 20,
		},
		{
			name:      "Zero to defaults",
			initial:   ZAction{PageOn: 0, PageLimit: 0},
			wantOn:    1,
			wantLimit: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DIN{ZAction: tt.initial}
			if err := d.Validate(); err != nil {
				t.Errorf("Validate() error: %v", err)
			}
			if d.ZAction.PageOn != tt.wantOn {
				t.Errorf("PageOn = %d, want %d", d.ZAction.PageOn, tt.wantOn)
			}
			if d.ZAction.PageLimit != tt.wantLimit {
				t.Errorf("PageLimit = %d, want %d", d.ZAction.PageLimit, tt.wantLimit)
			}
		})
	}
}

// TestDINNewPaginate verifies pagination creation, including cursor and errors.
func TestDINNewPaginate(t *testing.T) {
	tests := []struct {
		name       string
		din        *DIN
		totalItems int
		cursor     string
		wantErr    bool
		want       *Paginate
	}{
		{
			name:       "Valid with defaults",
			din:        &DIN{ZAction: ZAction{PageOn: 0, PageLimit: 0}},
			totalItems: 50,
			cursor:     "",
			wantErr:    false,
			want: &Paginate{
				CurrentPage:    1,
				TotalItems:     50,
				ItemsPerPage:   25,
				TotalPages:     2,
				Cursor:         "",
				PerPageOptions: DefaultPerPageOptions,
			},
		},
		{
			name:       "Valid with cursor",
			din:        &DIN{ZAction: ZAction{PageOn: 2, PageLimit: 10}},
			totalItems: 50,
			cursor:     "uuid-123",
			wantErr:    false,
			want: &Paginate{
				CurrentPage:    2,
				TotalItems:     50,
				ItemsPerPage:   10,
				TotalPages:     5,
				Cursor:         "uuid-123",
				PerPageOptions: DefaultPerPageOptions,
			},
		},
		{
			name:       "Negative totalItems error",
			din:        &DIN{ZAction: ZAction{PageOn: 1, PageLimit: 10}},
			totalItems: -1,
			cursor:     "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.din.NewPaginate(tt.totalItems, tt.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaginate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaginate() = %+v, want %+v", got, tt.want)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.din.Paginate, got) {
				t.Errorf("DIN.Paginate not set: %+v", tt.din.Paginate)
			}
		})
	}
}

// TestIDINPaginate verifies interface implementation.
func TestIDINPaginate(t *testing.T) {
	var idp IDINPaginate = &DIN{}
	if err := idp.Validate(); err != nil {
		t.Errorf("Validate() error: %v", err)
	}
	if _, err := idp.NewPaginate(50, ""); err != nil {
		t.Errorf("NewPaginate() error: %v", err)
	}
}
