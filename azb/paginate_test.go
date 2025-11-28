package azb

import (
	"bytes"
	"html/template"
	"reflect"
	"strings"
	"testing"
)

// TestNewPaginateValidations checks constructor validations, including cursor.
func TestNewPaginateValidations(t *testing.T) {
	tests := []struct {
		name         string
		currentPage  int
		totalItems   int
		itemsPerPage int
		cursor       string
		want         *Paginate
	}{
		{
			name:         "Negative currentPage",
			currentPage:  -1,
			totalItems:   50,
			itemsPerPage: 10,
			cursor:       "",
			want: &Paginate{
				CurrentPage:    1,
				TotalItems:     50,
				ItemsPerPage:   10,
				TotalPages:     5,
				PerPageOptions: DefaultPerPageOptions,
			},
		},
		{
			name:         "Negative totalItems",
			currentPage:  2,
			totalItems:   -50,
			itemsPerPage: 10,
			cursor:       "",
			want: &Paginate{
				CurrentPage:    2,
				TotalItems:     0,
				ItemsPerPage:   10,
				TotalPages:     0,
				PerPageOptions: DefaultPerPageOptions,
			},
		},
		{
			name:         "Negative itemsPerPage",
			currentPage:  2,
			totalItems:   50,
			itemsPerPage: -10,
			cursor:       "",
			want: &Paginate{
				CurrentPage:    2,
				TotalItems:     50,
				ItemsPerPage:   25,
				TotalPages:     2,
				PerPageOptions: DefaultPerPageOptions,
			},
		},
		{
			name:         "With cursor",
			currentPage:  1,
			totalItems:   50,
			itemsPerPage: 10,
			cursor:       "uuid-123",
			want: &Paginate{
				CurrentPage:    1,
				TotalItems:     50,
				ItemsPerPage:   10,
				TotalPages:     5,
				Cursor:         "uuid-123",
				PerPageOptions: DefaultPerPageOptions,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPaginate(tt.currentPage, tt.totalItems, tt.itemsPerPage, tt.cursor)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaginate() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestSetCurrentPage verifies page bounding.
func TestSetCurrentPage(t *testing.T) {
	p := NewPaginate(1, 50, 10, "")
	p.TotalPages = 5 // Set explicitly for test.

	tests := []struct {
		name string
		page int
		want int
	}{
		{"Valid page", 3, 3},
		{"Below range", 0, 1},
		{"Above range", 6, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.SetCurrentPage(tt.page)
			if p.CurrentPage != tt.want {
				t.Errorf("SetCurrentPage(%d) = %d, want %d", tt.page, p.CurrentPage, tt.want)
			}
		})
	}
}

// TestOffset verifies offset calculation, including cursor mode.
func TestOffset(t *testing.T) {
	tests := []struct {
		name         string
		currentPage  int
		itemsPerPage int
		cursor       string
		want         int
	}{
		{"Standard offset", 2, 10, "", 10},
		{"Cursor mode", 2, 10, "uuid-123", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPaginate(tt.currentPage, 50, tt.itemsPerPage, tt.cursor)
			if got := p.Offset(); got != tt.want {
				t.Errorf("Offset() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestLimit verifies limit return.
func TestLimit(t *testing.T) {
	p := NewPaginate(1, 50, 10, "")
	if got := p.Limit(); got != 10 {
		t.Errorf("Limit() = %d, want %d", got, 10)
	}
}

// TestNavigation verifies nav methods.
func TestNavigation(t *testing.T) {
	p := NewPaginate(2, 50, 10, "")
	p.TotalPages = 5

	tests := []struct {
		name string
		op   func()
		want int
	}{
		{"NavNext", p.NavNext, 3},
		{"NavPrev", p.NavPrev, 2},
		{"NavFirst", p.NavFirst, 1},
		{"NavLast", p.NavLast, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.op()
			if p.CurrentPage != tt.want {
				t.Errorf("%s() set CurrentPage to %d, want %d", tt.name, p.CurrentPage, tt.want)
			}
		})
	}
}

// TestPeekFunctions verifies peek methods.
func TestPeekFunctions(t *testing.T) {
	p := NewPaginate(2, 50, 10, "")
	p.TotalPages = 5

	if got := p.PeekNext(); got != 3 {
		t.Errorf("PeekNext() = %d, want 3", got)
	}
	if got := p.PeekPrev(); got != 1 {
		t.Errorf("PeekPrev() = %d, want 1", got)
	}

	// Boundaries
	p.CurrentPage = 5
	if got := p.PeekNext(); got != 5 {
		t.Errorf("PeekNext() at last = %d, want 5", got)
	}
	p.CurrentPage = 1
	if got := p.PeekPrev(); got != 1 {
		t.Errorf("PeekPrev() at first = %d, want 1", got)
	}
}

// TestPerPageOptions verifies custom options and methods.
func TestPerPageOptions(t *testing.T) {
	p := NewPaginate(1, 50, 10, "")

	// Default check
	if len(p.PerPageOptions) != 4 || p.PerPageOptions[1].Limit != 25 {
		t.Errorf("Default PerPageOptions mismatch: %+v", p.PerPageOptions)
	}

	// Set custom
	custom := []PaginateLimit{{Label: "10", Limit: 10, IsDefault: true}}
	p.SetPerPageOptions(custom)
	if !reflect.DeepEqual(p.PerPageOptions, custom) {
		t.Errorf("SetPerPageOptions() = %+v, want %+v", p.PerPageOptions, custom)
	}

	// Has and Get
	if !p.HasPerPageLimit(10) {
		t.Error("HasPerPageLimit(10) = false, want true")
	}
	if p.HasPerPageLimit(25) {
		t.Error("HasPerPageLimit(25) = true, want false after custom")
	}
	if got := p.GetPerPageLimitElseDefault(25); got != 10 {
		t.Errorf("GetPerPageLimitElseDefault(25) = %d, want 10 (default)", got)
	}
	if got := p.GetPerPageLimitElseDefault(10); got != 10 {
		t.Errorf("GetPerPageLimitElseDefault(10) = %d, want 10", got)
	}
}

// TestGetZClick verifies UI class helper.
func TestGetZClick(t *testing.T) {
	p := NewPaginate(1, 50, 10, "")
	if got := p.GetZClick(); got != "" {
		t.Errorf("GetZClick() false = %q, want empty", got)
	}

	p.AddZClick = true
	if got := p.GetZClick(); got != " zclick" {
		t.Errorf("GetZClick() true = %q, want ' zclick'", got)
	}
}

// TestPageNumbers verifies page slice generation.
func TestPageNumbers(t *testing.T) {
	p := NewPaginate(1, 50, 10, "")
	p.TotalPages = 5
	want := []int{1, 2, 3, 4, 5}
	if got := p.PageNumbers(); !reflect.DeepEqual(got, want) {
		t.Errorf("PageNumbers() = %v, want %v", got, want)
	}
}

// navPanelTemplate (use HasMatch method on Controls).
const navPanelTemplate = `
<nav class="pagination-nav{{ .GetZClick }}">
    {{- if not .LinkRender.NoStart }}{{ if gt .CurrentPage 1 }}<a href="{{ .ZUrl }}?cp=1" class="za-link">First</a>{{ end }}{{ end }}
    {{- if not .LinkRender.NoPrev }}{{ if gt .CurrentPage 1 }}<a href="{{ .ZUrl }}?cp={{ .PeekPrev }}" class="za-link">Prev</a>{{ end }}{{ end }}

    {{- range .PageNumbers }}
    {{- if eq . $.CurrentPage }}
    <span class="current">{{ . }}</span>
    {{- else }}
    <a href="{{ $.ZUrl }}?cp={{ . }}" class="za-link">{{ . }}</a>
    {{- end }}
    {{- end }}

    {{- if not .LinkRender.NoNext }}{{ if lt .CurrentPage .TotalPages }}<a href="{{ .ZUrl }}?cp={{ .PeekNext }}" class="za-link">Next</a>{{ end }}{{ end }}
    {{- if not .LinkRender.NoEnd }}{{ if lt .CurrentPage .TotalPages }}<a href="{{ .ZUrl }}?cp={{ .TotalPages }}" class="za-link">Last</a>{{ end }}{{ end }}

    {{- if .Label }}<span class="label">{{ .Label }}</span>{{ end }}
</nav>

<!-- Per-page selector if Controls include show-per-page -->
{{ if .Controls.HasMatch "show-per-page" }}
<select class="per-page-select">
    {{- range .PerPageOptions }}
    <option value="{{ .Limit }}" {{ if .IsDefault }}selected{{ end }}>{{ .Label }}</option>
    {{- end }}
</select>
{{ end }}
`

// TestRenderNavigationPanel verifies template rendering with new features.
func TestRenderNavigationPanel(t *testing.T) {
	p := NewPaginate(2, 50, 10, "")
	p.ZUrl = "/test"
	p.AddZClick = true
	p.Label = "Results"
	p.Controls = ZBTypes{ZBType("show-per-page")} // Use ZBType for elements.
	p.LinkRender.NoPrev = false

	tmpl, err := template.New("navPanel").Parse(navPanelTemplate)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := buf.String()
	expectedContains := []string{
		`pagination-nav zclick`,
		`/test?cp=1`,
		`Prev`,
		`<span class="current">2</span>`,
		`/test?cp=3`, // Next page link example.
		`Results`,
		`<select class="per-page-select">`,
		`<option value="25" selected>25</option>`, // From defaults.
	}
	for _, exp := range expectedContains {
		if !strings.Contains(got, exp) {
			t.Errorf("Rendered output missing '%s':\n%s", exp, got)
		}
	}
}

// TestDINMethods verifies DIN validation and pagination creation.
func TestDINMethods(t *testing.T) {
	d := &DIN{
		ZAction: ZAction{
			PageOn:    -1, // Invalid.
			PageLimit: 0,  // Invalid.
		},
	}

	// Validate fixes defaults.
	if err := d.Validate(); err != nil {
		t.Errorf("Validate() error: %v", err)
	}
	if d.ZAction.PageOn != 1 || d.ZAction.PageLimit != 25 {
		t.Errorf("Validate() set PageOn=%d (want 1), PageLimit=%d (want 25)", d.ZAction.PageOn, d.ZAction.PageLimit)
	}

	// NewPaginate.
	p, err := d.NewPaginate(50, "uuid-123")
	if err != nil {
		t.Errorf("NewPaginate() error: %v", err)
	}
	if p.ItemsPerPage != 25 || p.Cursor != "uuid-123" || d.Paginate != p {
		t.Errorf("NewPaginate() = %+v, mismatched", p)
	}

	// Error case.
	_, err = d.NewPaginate(-1, "")
	if err == nil {
		t.Error("NewPaginate() with negative totalItems should error")
	}
}

// TestInterfaces verifies implementations.
func TestInterfaces(t *testing.T) {
	p := NewPaginate(1, 100, 25, "")
	var ip IPaginate = p
	ip.NavNext()

	d := &DIN{}
	var idp IDINPaginate = d
	if err := idp.Validate(); err != nil {
		t.Error(err)
	}
	idp.NewPaginate(100, "")
}
