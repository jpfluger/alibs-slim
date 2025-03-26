package anode

import (
	"github.com/jpfluger/alibs-slim/alegal"
	"github.com/jpfluger/alibs-slim/atime"
	"testing"
)

func TestUserAccount_AddDeviceLogin(t *testing.T) {
	ua := UserAccount{}
	device := "Test Device"
	ip := "192.168.1.1"
	maxHistory := 5

	// Add a device login
	ua.AddDeviceLogin(device, ip, maxHistory)

	if len(ua.Logins) != 1 {
		t.Errorf("expected 1 login, got %d", len(ua.Logins))
	}

	if ua.Logins[0].Device != device {
		t.Errorf("expected device to be %s, got %s", device, ua.Logins[0].Device)
	}

	if ua.Logins[0].IP != ip {
		t.Errorf("expected IP to be %s, got %s", ip, ua.Logins[0].IP)
	}

	// Add more logins to exceed maxHistory
	for i := 0; i < maxHistory+1; i++ {
		ua.AddDeviceLogin(device, ip, maxHistory)
	}

	if len(ua.Logins) != maxHistory {
		t.Errorf("expected maxHistory logins, got %d", len(ua.Logins))
	}
}

func TestUserAccount_AddLDSByKey(t *testing.T) {
	ua := UserAccount{}
	key := alegal.LEGALKEY_TERMS
	effectiveDate := atime.GetNowUTCPointer()

	// Add a legal document signature
	hasChanges := ua.AddLDSByKey(key, false, effectiveDate)

	if !hasChanges {
		t.Errorf("expected hasChanges to be true, got %v", hasChanges)
	}

	if len(ua.LDS) != 1 {
		t.Errorf("expected 1 legal document signature, got %d", len(ua.LDS))
	}

	if ua.LDS[0].Key != key {
		t.Errorf("expected key to be %v, got %v", key, ua.LDS[0].Key)
	}

	if ua.LDS[0].EffectiveDate == nil || !ua.LDS[0].EffectiveDate.Equal(*effectiveDate) {
		t.Errorf("expected EffectiveDate to be %v, got %v", effectiveDate, ua.LDS[0].EffectiveDate)
	}

	// Try to add the same key again with appendIfFound set to false
	hasChanges = ua.AddLDSByKey(key, false, effectiveDate)

	if hasChanges {
		t.Errorf("expected hasChanges to be false, got %v", hasChanges)
	}

	// Try to add the same key again with appendIfFound set to true
	hasChanges = ua.AddLDSByKey(key, true, effectiveDate)

	if !hasChanges {
		t.Errorf("expected hasChanges to be true, got %v", hasChanges)
	}

	if len(ua.LDS) != 2 {
		t.Errorf("expected 2 legal document signatures, got %d", len(ua.LDS))
	}
}
