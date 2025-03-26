package acontact

import (
	"encoding/json"
	"sort"
)

// Mail represents a mailing address with a type and default status.
type Mail struct {
	Type      MailType    `json:"type,omitempty"`      // Type of the mailing address (e.g., home, business)
	Address   MailAddress `json:"address,omitempty"`   // The actual mailing address
	IsDefault bool        `json:"isDefault,omitempty"` // Indicates if this is the default mailing address
}

// Mails is a slice of Mail pointers, representing a collection of mailing addresses.
type Mails []*Mail

// FindByType searches for a mailing address by its type.
func (ms Mails) FindByType(mailType MailType) *Mail {
	return ms.findByType(mailType, false)
}

// FindByTypeOrDefault searches for a mailing address by its type or returns the default address.
func (ms Mails) FindByTypeOrDefault(mailType MailType) *Mail {
	return ms.findByType(mailType, true)
}

// findByType is a helper function that searches for a mailing address by type and optionally returns the default address.
func (ms Mails) findByType(mailType MailType, checkDefault bool) *Mail {
	var def *Mail
	for _, m := range ms {
		if m.Type.ToStringTrimLower() == mailType.ToStringTrimLower() {
			return m
		}
		if m.IsDefault {
			def = m
		}
	}
	if checkDefault && def != nil {
		return def
	}
	return nil
}

// HasType checks if a mailing address of the specified type exists in the collection.
func (ms Mails) HasType(mailType MailType) bool {
	return ms.FindByType(mailType) != nil
}

// HasTypeWithDefault checks if a mailing address of the specified type exists, or if there's a default address.
func (ms Mails) HasTypeWithDefault(mailType MailType, allowDefault bool) bool {
	return ms.findByType(mailType, allowDefault) != nil
}

// Clone creates a deep copy of the Mails collection.
func (ms Mails) Clone() Mails {
	b, err := json.Marshal(ms)
	if err != nil {
		return nil
	}
	var clone Mails
	if err := json.Unmarshal(b, &clone); err != nil {
		return nil
	}
	return clone
}

// MergeFrom adds mailing addresses from another collection that are not already present.
func (ms *Mails) MergeFrom(target Mails) {
	if ms == nil || target == nil {
		return
	}
	for _, t := range target {
		if t.Type.IsEmpty() {
			continue
		}
		isFound := false
		for _, m := range *ms {
			if m.Type.ToStringTrimLower() == t.Type.ToStringTrimLower() {
				isFound = true
				break
			}
		}
		if !isFound {
			*ms = append(*ms, t)
		}
	}
}

// Set adds or updates a mailing address in the collection.
func (ms *Mails) Set(mail *Mail) {
	if mail == nil || mail.Type.IsEmpty() {
		return
	}
	newMails := make(Mails, 0, len(*ms)+1)
	for _, m := range *ms {
		if m.Type.ToStringTrimLower() == mail.Type.ToStringTrimLower() {
			continue // Skip to replace the mailing address
		} else if m.IsDefault && mail.IsDefault {
			m.IsDefault = false // Unset the default if the new mailing address is the default
		}
		newMails = append(newMails, m)
	}
	newMails = append(newMails, mail) // Add the new mailing address

	// Sort the mailing addresses, placing the default address at the top
	sort.SliceStable(newMails, func(i, j int) bool {
		return newMails[i].IsDefault || newMails[i].Type < newMails[j].Type
	})

	*ms = newMails // Update the original collection
}

// Remove deletes a mailing address of the specified type from the collection.
func (ms *Mails) Remove(mailType MailType) {
	if mailType.IsEmpty() {
		return
	}
	newArr := make(Mails, 0)
	for _, m := range *ms {
		if m.Type.ToStringTrimLower() == mailType.ToStringTrimLower() {
			continue // Skip the mailing address to be removed
		}
		newArr = append(newArr, m)
	}
	*ms = newArr // Update the original collection with the remaining mailing addresses
}
