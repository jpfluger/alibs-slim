package azb

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

// TestDOUTInjectsSerialization verifies that DOUTInjects can be correctly serialized to JSON.
func TestDOUTInjectsSerialization(t *testing.T) {
	injects := DataInjects{
		&DataInject{Label: "TestLabel1", HTML: "<div>Content1</div>", JS: "console.log('test1');"},
		&DataInject{Label: "TestLabel2", HTML: "<div>Content2</div>", JS: "console.log('test2');"},
	}
	doutInjects := DOUTInjects{Injects: injects}

	// want := `{"dInjects":[{"label":"TestLabel1","html":"\u003cdiv\u003eContent1\u003c/div\u003e","js":"console.log('test1');"},{"label":"TestLabel2","html":"\u003cdiv\u003eContent2\u003c/div\u003e","js":"console.log('test2');"}]}`
	var b bytes.Buffer
	json.HTMLEscape(&b, []byte(`{"dInjects":[{"label":"TestLabel1","html":"<div>Content1</div>","js":"console.log('test1');"},{"label":"TestLabel2","html":"<div>Content2</div>","js":"console.log('test2');"}]}`))
	escapedWant := b.String()

	bytes, err := json.Marshal(doutInjects)
	if err != nil {
		t.Fatalf("Failed to marshal DOUTInjects: %v", err)
	}

	if got := string(bytes); got != escapedWant {
		t.Errorf("DOUTInjects serialization = %v, want %v", got, escapedWant)
	}
}

// TestDOUTInjectsDeserialization verifies that DOUTInjects can be correctly deserialized from JSON.
func TestDOUTInjectsDeserialization(t *testing.T) {
	jsonData := `{"dInjects":[{"label":"TestLabel1","html":"<div>Content1</div>","js":"console.log('test1');"},{"label":"TestLabel2","html":"<div>Content2</div>","js":"console.log('test2');"}]}`
	want := DOUTInjects{
		Injects: DataInjects{
			&DataInject{Label: "TestLabel1", HTML: "<div>Content1</div>", JS: "console.log('test1');"},
			&DataInject{Label: "TestLabel2", HTML: "<div>Content2</div>", JS: "console.log('test2');"},
		},
	}

	var got DOUTInjects
	err := json.Unmarshal([]byte(jsonData), &got)
	if err != nil {
		t.Fatalf("Failed to unmarshal DOUTInjects: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DOUTInjects deserialization got = %v, want %v", got, want)
	}
}

// TestNewDataInject verifies that a new DataInject can be created with the expected values.
func TestNewDataInject(t *testing.T) {
	label := "TestLabel"
	html := "<div>TestContent</div>"
	js := "console.log('test');"

	dataInject := DataInject{Label: label, HTML: html, JS: js}

	if dataInject.Label != label {
		t.Errorf("DataInject label got = %v, want %v", dataInject.Label, label)
	}
	if dataInject.HTML != html {
		t.Errorf("DataInject HTML got = %v, want %v", dataInject.HTML, html)
	}
	if dataInject.JS != js {
		t.Errorf("DataInject JS got = %v, want %v", dataInject.JS, js)
	}
}
