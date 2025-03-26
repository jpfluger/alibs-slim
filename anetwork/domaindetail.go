package anetwork

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
	"strings"
	"sync"
)

// DomainDetail holds the details of a domain including its name and certificate provider.
type DomainDetail struct {
	Name         string               `json:"name"`
	CertProvider ICertificateProvider `json:"certProvider"`
	mu           sync.RWMutex
}

// GetName returns the name of the domain.
func (d *DomainDetail) GetName() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Name
}

// GetCertProvider returns the certificate provider of the domain.
func (d *DomainDetail) GetCertProvider() ICertificateProvider {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.CertProvider
}

// Validate ensures the data is valid and includes an option to validate the certificate provider.
func (d *DomainDetail) Validate(dirCerts string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if d.CertProvider == nil {
		return fmt.Errorf("certProvider is nil")
	}
	return d.CertProvider.Validate(dirCerts)
}

// UnmarshalJSON custom unmarshals the JSON data into a DomainDetail object.
func (d *DomainDetail) UnmarshalJSON(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Alias to avoid recursion in UnmarshalJSON
	type Alias DomainDetail
	aux := &struct {
		Adapter json.RawMessage `json:"certProvider"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	// Unmarshal the main structure
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal DomainDetail: %v", err)
	}

	// Check if certProvider is empty
	if aux.Adapter == nil || len(aux.Adapter) == 0 {
		return fmt.Errorf("empty certProvider")
	}

	// Unmarshal certProvider to a map to extract the type
	var rawmap map[string]interface{}
	if err := json.Unmarshal(aux.Adapter, &rawmap); err != nil {
		return fmt.Errorf("failed to unmarshal certProvider: %v", err)
	}

	// Extract the type field
	rawType, ok := rawmap["type"].(string)
	if !ok {
		return fmt.Errorf("type field not found or is not a string in certProvider")
	}

	// Find the reflect type using the type manager
	rtype, err := areflect.TypeManager().FindReflectType(TYPEMANAGER_CERTIFICATEPROVIDERS, rawType)
	if err != nil {
		return fmt.Errorf("cannot find type struct '%s': %v", rawType, err)
	}

	// Create a new instance of the type and unmarshal the certProvider into it
	obj := reflect.New(rtype).Interface()
	if err = json.Unmarshal(aux.Adapter, obj); err != nil {
		return fmt.Errorf("failed to unmarshal certProvider where type is '%s': %v", rawType, err)
	}

	// Assert that the created object implements ICertificateProvider
	iCertProvider, ok := obj.(ICertificateProvider)
	if !ok {
		return fmt.Errorf("created object does not implement ICertificateProvider where type is '%s'", rawType)
	}
	d.CertProvider = iCertProvider

	return nil
}
