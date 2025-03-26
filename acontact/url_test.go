package acontact

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/anetwork"
	"testing"
)

// TestUrl_Validate tests the Validate method of the Url type
func TestUrl_Validate(t *testing.T) {
	validNetURL := anetwork.MustParseNetURL("http://example.com/")
	invalidNetURL := anetwork.MustParseNetURL("")

	tests := []struct {
		url  *Url
		want error
	}{
		{&Url{Type: "home", Title: "My Homepage", Link: validNetURL}, nil},
		{&Url{Type: "", Title: "No Type", Link: validNetURL}, fmt.Errorf("url type is empty")},
		{&Url{Type: "work", Title: " ", Link: nil}, fmt.Errorf("url link is empty")},
		{&Url{Type: "work", Title: "Invalid URL", Link: invalidNetURL}, fmt.Errorf("url link is not a valid URL")},
	}

	for _, tt := range tests {
		if err := tt.url.Validate(); err != nil && err.Error() != tt.want.Error() {
			t.Errorf("Url.Validate() error = %v, wantErr %v", err, tt.want)
		}
	}
}

// TestUrl_GetLinkWithOptions tests the GetLinkWithOptions method of the Url type
func TestUrl_GetLinkWithOptions(t *testing.T) {
	netURL := anetwork.MustParseNetURL("http://example.com/profile")

	tests := []struct {
		url  *Url
		part string
		want string
	}{
		{&Url{Link: netURL}, "raw", "http://example.com/profile"},
		{&Url{Link: netURL}, "domain", "example.com"},
		{&Url{Link: netURL}, "path", "/profile"},
		{&Url{Link: netURL}, "port", ""},
		{&Url{Link: netURL}, "scheme", "http"},
		{&Url{Link: netURL}, "no-scheme", "example.com/profile"},
		{&Url{Link: netURL}, "invalid", "http://example.com/profile"},
	}

	for _, tt := range tests {
		if got := tt.url.GetLinkWithOptions(tt.part); got != tt.want {
			t.Errorf("Url.GetLinkWithOptions() = %v, want %v", got, tt.want)
		}
	}
}

// TestUrls_FindByType tests the FindByType method of the Urls type
func TestUrls_FindByType(t *testing.T) {
	urls := Urls{
		&Url{Type: "home", Link: anetwork.MustParseNetURL("http://home.com")},
		&Url{Type: "work", Link: anetwork.MustParseNetURL("http://work.com")},
	}

	if got := urls.FindByType("home"); got.Link.Host != "home.com" {
		t.Errorf("Urls.FindByType() = %v, want %v", got.Link.Host, "home.com")
	}
	if got := urls.FindByType("work"); got.Link.Host != "work.com" {
		t.Errorf("Urls.FindByType() = %v, want %v", got.Link.Host, "work.com")
	}
	if got := urls.FindByType("missing"); got != nil {
		t.Errorf("Urls.FindByType() = %v, want %v", got, nil)
	}
}
