package auser

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aemail"
	"github.com/jpfluger/alibs-slim/atags"
	"strings"
)

// RecordUserIdentity represents a user identity within an information system.
type RecordUserIdentity struct {
	UID UID                `json:"uid,omitempty" jsonschema:"type=string,format=uuid,title=User ID,description=Unique identifier for the user,example=123e4567-e89b-12d3-a456-426614174000"`
	IDs atags.TagMapString `json:"ids,omitempty" jsonschema:"type=array,items=object,title=IDs,description=Collection of ID labels associated with the record"`
}

// NewRecordUserIdentity creates a new RecordUserIdentity with the provided values.
func NewRecordUserIdentity(uid UID, ids atags.TagMapString) RecordUserIdentity {
	return RecordUserIdentity{
		UID: uid,
		IDs: ids,
	}
}

// NewRecordUserIdentityByEmail creates a new RecordUserIdentity with the email and altId.
func NewRecordUserIdentityByEmail(email aemail.EmailAddress) RecordUserIdentity {
	return RecordUserIdentity{
		IDs: atags.TagMapString{"email": email.String()},
	}
}

// NewRecordUserIdentityByUID creates a new RecordUserIdentity with UUID.
func NewRecordUserIdentityByUID(uid UID) RecordUserIdentity {
	return RecordUserIdentity{
		UID: uid,
	}
}

// NewRecordUserIdentityById creates a new RecordUserIdentity with the altId.
func NewRecordUserIdentityById(label atags.TagKey, id string) RecordUserIdentity {
	return RecordUserIdentity{
		IDs: atags.TagMapString{label: id},
	}
}

// IsEmpty checks if the RecordUserIdentity is empty.
func (rui RecordUserIdentity) IsEmpty() bool {
	return rui.UID.IsNil() && rui.IDs.IsEmpty()
}

// String returns a comma-separated key-value representation of non-empty fields.
func (rui RecordUserIdentity) String() string {
	parts := []string{}

	if !rui.UID.IsNil() {
		parts = append(parts, fmt.Sprintf("uid=%s", rui.UID.String()))
	}
	if !rui.IDs.IsEmpty() {
		for key, val := range rui.IDs {
			parts = append(parts, fmt.Sprintf("%s=%s", key.String(), val))
		}
	}

	return fmt.Sprintf("%s", strings.Join(parts, ","))
}

// MarshalJSON customizes the JSON marshaling of RecordUserIdentity.
func (rui RecordUserIdentity) MarshalJSON() ([]byte, error) {
	// Use an anonymous struct for standard JSON marshaling
	type alias RecordUserIdentity
	return json.Marshal(&struct {
		alias
	}{
		alias: alias(rui),
	})
}

//// UnmarshalJSON customizes the JSON unmarshaling of RecordUserIdentity.
//func (rui *RecordUserIdentity) UnmarshalJSON(data []byte) error {
//	// Unmarshal JSON data to a string
//	var s string
//	if err := json.Unmarshal(data, &s); err != nil {
//		return err
//	}
//
//	// Use the ParseRecordUserIdentity function to populate the fields
//	parsed, err := ParseRecordUserIdentity(s)
//	if err != nil {
//		return err
//	}
//
//	*rui = parsed // Assign parsed values to the receiver
//	return nil
//}

func (rui *RecordUserIdentity) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// Parse the string format
		parsed, err := ParseRecordUserIdentity(s)
		if err != nil {
			return err
		}
		*rui = parsed
		return nil
	}

	// If it's not a string, try unmarshaling as a JSON object
	var obj struct {
		UID string             `json:"uid"`
		IDs atags.TagMapString `json:"ids"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	rui.UID = ParseUID(obj.UID)
	rui.IDs = obj.IDs
	return nil
}

// ParseRecordUserIdentity parses a comma-separated key-value representation
// and returns a new RecordUserIdentity.
func ParseRecordUserIdentity(input string) (RecordUserIdentity, error) {
	identity := RecordUserIdentity{
		IDs: make(atags.TagMapString), // Ensure the map is initialized
	}
	pairs := strings.Split(input, ",")

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return identity, fmt.Errorf("invalid key-value pair format in '%s'", pair)
		}

		key, value := kv[0], kv[1]
		switch key {
		case "uid":
			identity.UID = ParseUID(value)
			if identity.UID.IsNil() {
				return identity, fmt.Errorf("invalid UID")
			}
		default:
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key != "" && value != "" {
				identity.IDs.Set(atags.TagKey(key), value)
			}
		}
	}

	return identity, nil
}

// GetEmail locates the email id and returns it
func (rui RecordUserIdentity) GetEmail() aemail.EmailAddress {
	if rui.IDs.IsEmpty() {
		return ""
	}
	return aemail.EmailAddress(rui.IDs.Value("email"))
}

func (rui RecordUserIdentity) FindLabel(key atags.TagKey) string {
	if key.IsEmpty() {
		return ""
	}
	return rui.IDs.Value(key)
}

func (rui RecordUserIdentity) HasMatch(user RecordUserIdentity) bool {
	// Check if the UID matches and is not nil
	if !rui.UID.IsNil() && rui.UID == user.UID {
		return true
	}

	// Check if any keys in the IDs match
	for key, value := range rui.IDs {
		if value != "" && value == user.IDs.Value(key) {
			return true
		}
	}

	// No matches found
	return false
}
