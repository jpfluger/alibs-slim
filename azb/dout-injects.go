package azb

// DOUTInjects holds a collection of Data Output Injections.
type DOUTInjects struct {
	Injects DataInjects `json:"dInjects"` // The JSON key "dInjects" maps to a slice of DataInject pointers.
}

// DataInject represents a single data injection with label, HTML, and JS content.
type DataInject struct {
	Label string `json:"label"` // The label for the data injection.
	HTML  string `json:"html"`  // The HTML content to be injected.
	JS    string `json:"js"`    // The JavaScript code to be injected.
}

// DataInjects is a slice of pointers to DataInject, allowing for a collection of injections.
type DataInjects []*DataInject
