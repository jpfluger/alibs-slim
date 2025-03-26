package azb

import (
	"bytes"
	"html/template"
	"testing"
)

func TestNewPaginateValidations(t *testing.T) {
	tests := []struct {
		currentPage  int
		totalItems   int
		itemsPerPage int
		want         *Paginate
	}{
		{currentPage: -1, totalItems: 50, itemsPerPage: 10, want: &Paginate{CurrentPage: 1, TotalItems: 50, ItemsPerPage: 10, TotalPages: 5}},
		{currentPage: 2, totalItems: -50, itemsPerPage: 10, want: &Paginate{CurrentPage: 2, TotalItems: 0, ItemsPerPage: 10, TotalPages: 0}},
		{currentPage: 2, totalItems: 50, itemsPerPage: -10, want: &Paginate{CurrentPage: 2, TotalItems: 50, ItemsPerPage: 25, TotalPages: 2}},
	}

	for _, tt := range tests {
		got := NewPaginate(tt.currentPage, tt.totalItems, tt.itemsPerPage)
		if *got != *tt.want {
			t.Errorf("NewPaginate(%d, %d, %d) = %+v, want %+v", tt.currentPage, tt.totalItems, tt.itemsPerPage, got, tt.want)
		}
	}
}

// TestSetCurrentPage verifies that SetCurrentPage correctly sets the current page.
func TestSetCurrentPage(t *testing.T) {
	paginate := &Paginate{TotalPages: 5}

	paginate.SetCurrentPage(3)
	if paginate.CurrentPage != 3 {
		t.Errorf("SetCurrentPage() = %d, want %d", paginate.CurrentPage, 3)
	}

	// Test boundary conditions
	paginate.SetCurrentPage(0)
	if paginate.CurrentPage != 1 {
		t.Errorf("SetCurrentPage() with below range = %d, want %d", paginate.CurrentPage, 1)
	}

	paginate.SetCurrentPage(6)
	if paginate.CurrentPage != 5 {
		t.Errorf("SetCurrentPage() with above range = %d, want %d", paginate.CurrentPage, 5)
	}
}

// TestOffset verifies that Offset calculates the correct query offset.
func TestOffset(t *testing.T) {
	paginate := &Paginate{CurrentPage: 2, ItemsPerPage: 10}
	expectedOffset := 10 // Page 2 with 10 items per page should have an offset of 10.

	if paginate.Offset() != expectedOffset {
		t.Errorf("Offset() = %d, want %d", paginate.Offset(), expectedOffset)
	}
}

// TestLimit verifies that Limit returns the correct items per page.
func TestLimit(t *testing.T) {
	paginate := &Paginate{ItemsPerPage: 10}

	if paginate.Limit() != paginate.ItemsPerPage {
		t.Errorf("Limit() = %d, want %d", paginate.Limit(), paginate.ItemsPerPage)
	}
}

// Test navigation functions (NavNext, NavPrev, NavFirst, NavLast).
func TestNavigation(t *testing.T) {
	paginate := &Paginate{CurrentPage: 2, TotalPages: 5}

	paginate.NavNext()
	if paginate.CurrentPage != 3 {
		t.Errorf("NavNext() = %d, want %d", paginate.CurrentPage, 3)
	}

	paginate.NavPrev()
	if paginate.CurrentPage != 2 {
		t.Errorf("NavPrev() = %d, want %d", paginate.CurrentPage, 2)
	}

	paginate.NavFirst()
	if paginate.CurrentPage != 1 {
		t.Errorf("NavFirst() = %d, want %d", paginate.CurrentPage, 1)
	}

	paginate.NavLast()
	if paginate.CurrentPage != 5 {
		t.Errorf("NavLast() = %d, want %d", paginate.CurrentPage, 5)
	}
}

// Test PeekNext and PeekPrev functions.
func TestPeekFunctions(t *testing.T) {
	paginate := &Paginate{CurrentPage: 2, TotalPages: 5}

	nextPage := paginate.PeekNext()
	if nextPage != 3 {
		t.Errorf("PeekNext() = %d, want %d", nextPage, 3)
	}

	prevPage := paginate.PeekPrev()
	if prevPage != 1 {
		t.Errorf("PeekPrev() = %d, want %d", prevPage, 1)
	}

	// Test boundary conditions
	paginate.CurrentPage = 5
	if paginate.PeekNext() != 5 {
		t.Errorf("PeekNext() at last page = %d, want %d", paginate.PeekNext(), 5)
	}

	paginate.CurrentPage = 1
	if paginate.PeekPrev() != 1 {
		t.Errorf("PeekPrev() at first page = %d, want %d", paginate.PeekPrev(), 1)
	}
}

// Define the Go HTML template for the navigation panel.
const navPanelTemplate = `
<nav class="pagination-nav">
    {{- if gt .CurrentPage 1 }}
    <a href="?cp=1">First</a>
    <a href="?cp={{ .PeekPrev }}">Prev</a>
    {{- end }}

    {{- range .PageNumbers }}
    {{- if eq . $.CurrentPage }}
    <span class="current">{{ . }}</span>
    {{- else }}
    <a href="?cp={{ . }}">{{ . }}</a>
    {{- end }}
    {{- end }}

    {{- if lt .CurrentPage .TotalPages }}
    <a href="?cp={{ .PeekNext }}">Next</a>
    <a href="?cp={{ .TotalPages }}">Last</a>
    {{- end }}
</nav>
`

// PageNumbers generates a slice of page numbers for pagination display.
func (p *Paginate) PageNumbers() []int {
	var pages []int
	for i := 1; i <= p.TotalPages; i++ {
		pages = append(pages, i)
	}
	return pages
}

// TestRenderNavigationPanel verifies that the navigation panel is rendered correctly.
func TestRenderNavigationPanel(t *testing.T) {
	// Create a Paginate struct with known values.
	paginate := &Paginate{
		CurrentPage:  2,
		TotalItems:   50,
		ItemsPerPage: 10,
		TotalPages:   5,
	}

	// Define the expected output.
	expectedOutput := `
<nav class="pagination-nav">
    <a href="?cp=1">First</a>
    <a href="?cp=1">Prev</a>
    <a href="?cp=1">1</a>
    <span class="current">2</span>
    <a href="?cp=3">3</a>
    <a href="?cp=4">4</a>
    <a href="?cp=5">5</a>
    <a href="?cp=3">Next</a>
    <a href="?cp=5">Last</a>
</nav>
`

	// Parse the template.
	tmpl, err := template.New("navPanel").Parse(navPanelTemplate)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Execute the template with the Paginate struct.
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, paginate)
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Compare the output to the expected output.
	if buf.String() != expectedOutput {
		t.Errorf("Rendered navigation panel was incorrect, got: %s, want: %s.", buf.String(), expectedOutput)
	}
}
