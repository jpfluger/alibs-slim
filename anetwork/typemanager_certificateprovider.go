package anetwork

import (
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
)

const TYPEMANAGER_CERTIFICATEPROVIDERS = "certificateproviders"

// init registers the contact types with the type manager upon package initialization.
func init() {
	// Ignoring the error as per the original code, but consider handling it.
	_ = areflect.TypeManager().Register(TYPEMANAGER_CERTIFICATEPROVIDERS, "anetwork", returnTypeManagerCertificateProviders)
}

// returnTypeManagerCertificateProviders returns the reflect.Type corresponding to the provided typeName.
func returnTypeManagerCertificateProviders(typeName string) (reflect.Type, error) {
	var rtype reflect.Type
	switch CertificateProviderType(typeName) {
	case CERTIFICATEPROVIDERTYPE_SELFSIGN:
		rtype = reflect.TypeOf(CertificateProviderSelfSign{})
	case CERTIFICATEPROVIDERTYPE_SIGNED:
		rtype = reflect.TypeOf(CertificateProviderSigned{})
	default:
		return nil, nil
	}
	return rtype, nil
}
