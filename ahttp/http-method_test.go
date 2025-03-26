package ahttp

import (
	"testing"
)

// TestHttpMethod_IsEmpty tests the IsEmpty method of HttpMethod.
func TestHttpMethod_IsEmpty(t *testing.T) {
	var hm HttpMethod
	hm = "  "
	if !hm.IsEmpty() {
		t.Errorf("Expected HttpMethod to be empty")
	}

	hm = "GET"
	if hm.IsEmpty() {
		t.Errorf("Expected HttpMethod not to be empty")
	}
}

// TestHttpMethod_TrimSpace tests the TrimSpace method of HttpMethod.
func TestHttpMethod_TrimSpace(t *testing.T) {
	hm := HttpMethod("  POST  ")
	expected := HttpMethod("POST")
	if hm.TrimSpace() != expected {
		t.Errorf("TrimSpace() = %v, want %v", hm.TrimSpace(), expected)
	}
}

// TestHttpMethod_String tests the String method of HttpMethod.
func TestHttpMethod_String(t *testing.T) {
	hm := HttpMethod("DELETE")
	if hm.String() != "DELETE" {
		t.Errorf("String() = %v, want %v", hm.String(), "DELETE")
	}
}

// TestHttpMethod_ToStringTrimUpper tests the ToStringTrimUpper method of HttpMethod.
func TestHttpMethod_ToStringTrimUpper(t *testing.T) {
	hm := HttpMethod("  patch  ")
	if hm.ToStringTrimUpper() != "PATCH" {
		t.Errorf("ToStringTrimUpper() = %v, want %v", hm.ToStringTrimUpper(), "PATCH")
	}
}

// TestHttpMethod_TrimSpaceToUpper tests the TrimSpaceToUpper method of HttpMethod.
func TestHttpMethod_TrimSpaceToUpper(t *testing.T) {
	hm := HttpMethod("  put  ")
	expected := HttpMethod("PUT")
	if hm.TrimSpaceToUpper() != expected {
		t.Errorf("TrimSpaceToUpper() = %v, want %v", hm.TrimSpaceToUpper(), expected)
	}
}

// TestHttpMethods_Find tests the Find method of HttpMethods.
func TestHttpMethods_Find(t *testing.T) {
	hms := HttpMethods{HTTPMETHOD_GET, HTTPMETHOD_POST}
	target := HttpMethod("POST")
	if hms.Find(target) != target {
		t.Errorf("Find() did not find the HttpMethod %v", target)
	}

	target = HttpMethod("PUT")
	if hms.Find(target) != "" {
		t.Errorf("Find() should not find the HttpMethod %v", target)
	}
}

// TestHttpMethods_Has tests the Has method of HttpMethods.
func TestHttpMethods_Has(t *testing.T) {
	hms := HttpMethods{HTTPMETHOD_GET, HTTPMETHOD_POST}
	if !hms.Has(HTTPMETHOD_GET) {
		t.Errorf("Has() should have found HttpMethod GET")
	}

	if hms.Has(HTTPMETHOD_DELETE) {
		t.Errorf("Has() should not have found HttpMethod DELETE")
	}
}
