package auuids

import (
	"testing"
)

func TestNewIDLabel(t *testing.T) {
	validLabel := UUIDLabel("testLabel")
	validID := NewUUID()
	nilLabel := UUIDLabel("")
	nilID := UUID{}

	t.Run("Valid IDLabel", func(t *testing.T) {
		idLabel := NewIDLabel(validLabel, validID)
		if idLabel == nil || idLabel.Label != validLabel || idLabel.Id != validID {
			t.Errorf("Expected valid IDLabel, got %v", idLabel)
		}
	})

	t.Run("Nil Label", func(t *testing.T) {
		idLabel := NewIDLabel(nilLabel, validID)
		if idLabel != nil {
			t.Errorf("Expected nil IDLabel for empty label, got %v", idLabel)
		}
	})

	t.Run("Nil ID", func(t *testing.T) {
		idLabel := NewIDLabel(validLabel, nilID)
		if idLabel != nil {
			t.Errorf("Expected nil IDLabel for nil ID, got %v", idLabel)
		}
	})
}

func TestValidateIDLabel(t *testing.T) {
	validLabel := UUIDLabel("testLabel")
	validID := NewUUID()

	t.Run("Valid IDLabel", func(t *testing.T) {
		idLabel := NewIDLabel(validLabel, validID)
		err := idLabel.Validate()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Nil IDLabel", func(t *testing.T) {
		var idLabel *IDLabel
		err := idLabel.Validate()
		if err == nil {
			t.Errorf("Expected error for nil IDLabel, got nil")
		}
	})

	t.Run("Empty Label", func(t *testing.T) {
		idLabel := NewIDLabel("", validID)
		if idLabel != nil {
			t.Errorf("Expected nil IDLabel for empty label")
		}
	})

	t.Run("Nil ID", func(t *testing.T) {
		idLabel := NewIDLabel(validLabel, UUID{})
		if idLabel != nil {
			t.Errorf("Expected nil IDLabel for nil ID")
		}
	})
}

func TestIDLabelHasMatch(t *testing.T) {
	label := UUIDLabel("testLabel")
	id := NewUUID()

	t.Run("Matching Labels", func(t *testing.T) {
		idLabel1 := NewIDLabel(label, id)
		idLabel2 := NewIDLabel(label, id)
		if !idLabel1.HasMatch(idLabel2) {
			t.Errorf("Expected HasMatch to return true, got false")
		}
	})

	t.Run("Non-Matching Labels", func(t *testing.T) {
		idLabel1 := NewIDLabel(label, id)
		idLabel2 := NewIDLabel(label, NewUUID())
		if idLabel1.HasMatch(idLabel2) {
			t.Errorf("Expected HasMatch to return false, got true")
		}
	})
}

func TestIDLabelsSet(t *testing.T) {
	label := UUIDLabel("testLabel")
	id := NewUUID()
	idLabels := IDLabels{}

	t.Run("Add New IDLabel", func(t *testing.T) {
		newLabel := NewIDLabel(label, id)
		updatedLabels, err := idLabels.Set(newLabel)
		if err != nil || len(updatedLabels) != 1 || updatedLabels[0] != newLabel {
			t.Errorf("Expected to add IDLabel, got %v, err: %v", updatedLabels, err)
		}
	})

	t.Run("Update Existing IDLabel", func(t *testing.T) {
		newLabel := NewIDLabel(label, id)
		idLabels, _ = idLabels.Set(newLabel)
		updatedLabel := NewIDLabel(label, NewUUID())
		updatedLabels, err := idLabels.SetMatchByLabel(updatedLabel)
		if err != nil || len(updatedLabels) != 1 || updatedLabels[0] != updatedLabel {
			t.Errorf("Expected to update IDLabel, got %v, err: %v", updatedLabels, err)
		}
	})
}

func TestIDLabelsRemove(t *testing.T) {
	label := UUIDLabel("testLabel")
	id := NewUUID()
	idLabels := IDLabels{NewIDLabel(label, id)}

	t.Run("Remove Existing IDLabel", func(t *testing.T) {
		target := NewIDLabel(label, id)
		updatedLabels, err := idLabels.Remove(target)
		if err != nil || len(updatedLabels) != 0 {
			t.Errorf("Expected to remove IDLabel, got %v, err: %v", updatedLabels, err)
		}
	})

	t.Run("Remove Non-Existent IDLabel", func(t *testing.T) {
		target := NewIDLabel(label, NewUUID())
		updatedLabels, err := idLabels.Remove(target)
		if err != nil || len(updatedLabels) != 1 {
			t.Errorf("Expected no removal, got %v, err: %v", updatedLabels, err)
		}
	})
}

func TestFromJSON(t *testing.T) {
	jsonStr := `[{"label":"Label1","id":"123e4567-e89b-12d3-a456-426614174000"}]`
	var labels IDLabels

	err := labels.FromJSON(jsonStr)
	if err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}

	if len(labels) != 1 || labels[0].Label != "Label1" {
		t.Errorf("Parsed labels mismatch")
	}
}

func TestHasMatch(t *testing.T) {
	id1 := ParseUUID("123e4567-e89b-12d3-a456-426614174000")
	id2 := ParseUUID("123e4567-e89b-12d3-a456-426614174001")
	id3 := ParseUUID("123e4567-e89b-12d3-a456-426614174002")

	label1 := &IDLabel{Label: "Label1", Id: id1}
	label2 := &IDLabel{Label: "Label2", Id: id2}
	label3 := &IDLabel{Label: "Label3", Id: id3}

	sourceLabels := IDLabels{label1, label2}

	// Test matching labels
	if !sourceLabels.HasMatch(label2) {
		t.Errorf("Expected HasMatch to return true, but got false")
	}

	// Test multiple targets
	if !sourceLabels.HasMatch(label2, label3) {
		t.Errorf("Expected HasMatch to return true for one matching target, but got false")
	}

	// Test non-matching labels
	if sourceLabels.HasMatch(label3) {
		t.Errorf("Expected HasMatch to return false, but got true")
	}

	// Test nil targets
	if sourceLabels.HasMatch(nil) {
		t.Errorf("Expected HasMatch to return false for nil targets, but got true")
	}

	// Test nil source labels
	if IDLabels(nil).HasMatch(label1) {
		t.Errorf("Expected HasMatch to return false for nil source labels, but got true")
	}
}
