package anode

import (
	"testing"
	"time"
)

func TestAdminLock_IsPasswordLocked(t *testing.T) {
	lock := AdminLock{
		IsPasswordLocked: true,
	}

	if !lock.IsPasswordLocked {
		t.Errorf("expected IsPasswordLocked to be true, got %v", lock.IsPasswordLocked)
	}
}

func TestAdminLock_Date(t *testing.T) {
	now := time.Now()
	lock := AdminLock{
		Date: &now,
	}

	if lock.Date == nil || !lock.Date.Equal(now) {
		t.Errorf("expected Date to be %v, got %v", now, lock.Date)
	}
}

func TestAdminLock_Message(t *testing.T) {
	message := "Account locked due to suspicious activity"
	lock := AdminLock{
		Message: message,
	}

	if lock.Message != message {
		t.Errorf("expected Message to be %v, got %v", message, lock.Message)
	}
}

func TestAdminLock_RequestResetPassword(t *testing.T) {
	now := time.Now()
	lock := AdminLock{
		RequestResetPassword: &now,
	}

	if lock.RequestResetPassword == nil || !lock.RequestResetPassword.Equal(now) {
		t.Errorf("expected RequestResetPassword to be %v, got %v", now, lock.RequestResetPassword)
	}
}

func TestAdminLock_FullStruct(t *testing.T) {
	now := time.Now()
	message := "Account locked due to suspicious activity"
	lock := AdminLock{
		IsPasswordLocked:     true,
		Date:                 &now,
		Message:              message,
		RequestResetPassword: &now,
	}

	if !lock.IsPasswordLocked {
		t.Errorf("expected IsPasswordLocked to be true, got %v", lock.IsPasswordLocked)
	}
	if lock.Date == nil || !lock.Date.Equal(now) {
		t.Errorf("expected Date to be %v, got %v", now, lock.Date)
	}
	if lock.Message != message {
		t.Errorf("expected Message to be %v, got %v", message, lock.Message)
	}
	if lock.RequestResetPassword == nil || !lock.RequestResetPassword.Equal(now) {
		t.Errorf("expected RequestResetPassword to be %v, got %v", now, lock.RequestResetPassword)
	}
}
