package anode

import (
	"testing"
	"time"
)

func TestUserAccountIsOnMFA_IsVerified(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		uat      *UserAccountIsOnMFA
		expected bool
	}{
		{"Nil struct", nil, false},
		{"Nil Verified", &UserAccountIsOnMFA{}, false},
		{"Zero Verified", &UserAccountIsOnMFA{Verified: &time.Time{}}, false},
		{"Valid Verified", &UserAccountIsOnMFA{Verified: &now}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uat.IsVerified(); got != tt.expected {
				t.Errorf("UserAccountIsOnMFA.IsVerified() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUserAccountIsOnMFA_Fields(t *testing.T) {
	now := time.Now()
	uat := UserAccountIsOnMFA{
		IsOn:     true,
		Created:  &now,
		Verified: &now,
	}

	if !uat.IsOn {
		t.Errorf("expected IsOn to be true, got %v", uat.IsOn)
	}

	if uat.Created == nil || !uat.Created.Equal(now) {
		t.Errorf("expected Created to be %v, got %v", now, uat.Created)
	}

	if uat.Verified == nil || !uat.Verified.Equal(now) {
		t.Errorf("expected Verified to be %v, got %v", now, uat.Verified)
	}
}
