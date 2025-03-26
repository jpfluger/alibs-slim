package aclient_http

// HOB (HTTP Object) encapsulates all options for HTTP requests.
type HOB struct {
	Path              string
	Payload           interface{}
	Raw               []byte
	ContentType       string
	ExpectedType      string
	ConnectionTimeout int
}

// NewHOBGet creates a new HOB for GET requests.
func NewHOBGet(path string) *HOB {
	return &HOB{
		Path:              path,
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBPost creates a new HOB for POST requests.
func NewHOBPost(path string, payload interface{}) *HOB {
	return &HOB{
		Path:              path,
		Payload:           payload,
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBPostRaw creates a new HOB with raw content for POST requests.
func NewHOBPostRaw(path string, raw []byte) *HOB {
	return &HOB{
		Path:              path,
		Raw:               raw,
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBWithJSON creates a new HOB for JSON requests.
func NewHOBWithJSON(path string, payload interface{}) *HOB {
	return &HOB{
		Path:              path,
		Payload:           payload,
		ContentType:       "application/json",
		ExpectedType:      "application/json",
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBWithJSONRaw creates a new HOB with raw JSON content for JSON requests.
func NewHOBWithJSONRaw(path string, raw []byte) *HOB {
	return &HOB{
		Path:              path,
		Raw:               raw,
		ContentType:       "application/json",
		ExpectedType:      "application/json",
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBWithXML creates a new HOB for XML requests.
func NewHOBWithXML(path string, payload interface{}) *HOB {
	return &HOB{
		Path:              path,
		Payload:           payload,
		ContentType:       "application/xml",
		ExpectedType:      "application/xml",
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}

// NewHOBWithXMLRaw creates a new HOB with raw XML content for XML requests.
func NewHOBWithXMLRaw(path string, raw []byte) *HOB {
	return &HOB{
		Path:              path,
		Raw:               raw,
		ContentType:       "application/xml",
		ExpectedType:      "application/xml",
		ConnectionTimeout: HTTP_CONNECTION_TIMEOUT,
	}
}
