package acrypt

import "testing"

// Unit tests for SecretsItem and SecretsItems.
func TestSecretsItem_GetDecodedValue(t *testing.T) {
	key := SecretsKey("test_key")
	password := "test_password"
	encodedValue := "d;plain;aes256;test_data"
	si := NewSecretsItem(key, encodedValue, ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	si.SetValueDecrypted("test_data", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)

	decodedValue, err := si.GetDecodedValue(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(decodedValue) != "test_data" {
		t.Errorf("expected %s, got %s", "test_data", string(decodedValue))
	}
}

func TestSecretsItem_SetValueDecrypted(t *testing.T) {
	key := SecretsKey("test_key")
	si := NewSecretsItem(key, "", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)

	si.SetValueDecrypted("updated_data", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	decodedValue, err := si.GetDecodedValue("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(decodedValue) != "updated_data" {
		t.Errorf("expected %s, got %s", "updated_data", string(decodedValue))
	}
}

func TestSecretsItems_Find(t *testing.T) {
	key := SecretsKey("find_key")
	keyNotFound := SecretsKey("missing_key")
	si := NewSecretsItem(key, "value", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	sis := SecretsItems{si}

	found := sis.Find(key)
	if found == nil || found.GetKey() != key {
		t.Errorf("expected to find key %s", key)
	}

	notFound := sis.Find(keyNotFound)
	if notFound != nil {
		t.Errorf("did not expect to find key %s", keyNotFound)
	}
}

func TestSecretsItems_SetAndRemove(t *testing.T) {
	key := SecretsKey("set_key")
	si := NewSecretsItem(key, "value", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	sis := SecretsItems{}

	err := sis.Set(si)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := sis.Find(key)
	if found == nil {
		t.Errorf("expected to find key %s after set", key)
	}

	sis.Remove(key)
	found = sis.Find(key)
	if found != nil {
		t.Errorf("did not expect to find key %s after remove", key)
	}
}
