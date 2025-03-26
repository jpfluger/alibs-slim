package acontact

import (
	"errors"
	"fmt"
	"strings"
)

// Name represents a person or entity's name details.
type Name struct {
	First   string `json:"first,omitempty"`   // First name
	Last    string `json:"last,omitempty"`    // Last name
	Company string `json:"company,omitempty"` // Company name (if applicable)
	Full    string `json:"full,omitempty"`    // Full name (optional; can be derived)
	Short   string `json:"short,omitempty"`   // Short name (optional; can be derived)
	Legal   string `json:"legal,omitempty"`   // Legal name (optional; use if needed)

	Title      string `json:"title,omitempty"`      // Title of the person (e.g., Mr., Dr.)
	Department string `json:"department,omitempty"` // Department within the company
}

// Validate ensures that the Name has at least one required field (First, Last, or Company).
func (n *Name) Validate() error {
	// Trim all fields to remove unnecessary spaces.
	n.First = strings.TrimSpace(n.First)
	n.Last = strings.TrimSpace(n.Last)
	n.Company = strings.TrimSpace(n.Company)
	n.Full = strings.TrimSpace(n.Full)
	n.Short = strings.TrimSpace(n.Short)
	n.Legal = strings.TrimSpace(n.Legal)
	n.Title = strings.TrimSpace(n.Title)
	n.Department = strings.TrimSpace(n.Department)

	// Ensure at least one identifying field is populated.
	if !n.IsPerson() && !n.IsEntity() {
		return errors.New("at least one of First, Last, or Company must be provided")
	}

	return nil
}

// IsPerson determines if the Name represents a person.
func (n *Name) IsPerson() bool {
	return n.First != "" || n.Last != ""
}

// IsEntity determines if the Name represents an entity.
func (n *Name) IsEntity() bool {
	return n.Company != ""
}

// IsBothPersonAndEntity checks if the Name represents both a person and an entity.
func (n *Name) IsBothPersonAndEntity() bool {
	return n.IsPerson() && n.IsEntity()
}

// GetFirstLastName returns the concatenation of First and Last names, if available.
func (n *Name) GetFirstLastName() string {
	if n.First == "" && n.Last == "" {
		return ""
	}
	if n.First == "" && n.Last != "" {
		return n.Last
	}
	if n.Last == "" && n.First != "" {
		return n.First
	}
	return fmt.Sprintf("%s %s", n.First, n.Last)
}

func (n *Name) GetCompany() string {
	return n.Company
}

// GetName returns the most appropriate name for display. Priority: First/Last -> Company.
func (n *Name) GetName() string {
	if name := n.GetFirstLastName(); name != "" {
		return name
	}
	return n.Company
}

// GetNamePlusCompany returns a combination of First/Last name and Company if both exist.
func (n *Name) GetNamePlusCompany() string {
	name := n.GetFirstLastName()
	if name != "" && n.Company != "" {
		return fmt.Sprintf("%s (%s)", name, n.Company)
	}
	if name != "" {
		return name
	}
	return n.Company
}

// MustGetFull ensures Full is returned, or falls back to derived options.
func (n *Name) MustGetFull() string {
	if n.Full != "" {
		return n.Full
	}
	if n.Legal != "" {
		return n.Legal
	}
	return n.GetName()
}

// MustGetShort ensures Short is returned, or falls back to derived options.
func (n *Name) MustGetShort() string {
	if n.Short != "" {
		return n.Short
	}
	return n.GetName()
}

// MustGetLegal ensures Legal is returned, or falls back to derived options.
func (n *Name) MustGetLegal() string {
	if n.Legal != "" {
		return n.Legal
	}
	return n.GetName()
}

// GetMail returns a prioritized list of names (Company and/or First/Last).
func (n *Name) GetMail() []string {
	var result []string
	if n.IsBothPersonAndEntity() {
		result = append(result, n.Company, n.GetFirstLastName())
	} else if n.IsPerson() {
		result = append(result, n.GetFirstLastName())
	} else if n.IsEntity() {
		result = append(result, n.Company)
	}
	return result
}

// GetMailWithTitleDepartment returns a prioritized list of names (Company and/or First/Last/Title).
func (n *Name) GetMailWithTitleDepartment() []string {
	var result []string

	// If both a person and an entity
	if n.IsBothPersonAndEntity() {
		if n.Title != "" {
			result = append(result, fmt.Sprintf("%s %s", n.Title, n.GetFirstLastName()))
		} else {
			result = append(result, n.GetFirstLastName())
		}
		result = append(result, n.Company)
	} else if n.IsPerson() {
		if n.Title != "" {
			result = append(result, fmt.Sprintf("%s %s", n.Title, n.GetFirstLastName()))
		} else {
			result = append(result, n.GetFirstLastName())
		}
		if n.Department != "" {
			result = append(result, n.Department)
		}
	} else if n.IsEntity() {
		result = append(result, n.Company)
	}

	return result
}

func (n *Name) CanSign(requireCompany bool) error {
	n.Title = strings.TrimSpace(n.Title)
	if n.Title == "" {
		return fmt.Errorf("title is empty")
	}
	if !n.IsPerson() {
		return fmt.Errorf("name is empty")
	}
	if requireCompany {
		if !n.IsEntity() {
			return fmt.Errorf("company is empty")
		}
	}
	return nil
}

func (n *Name) GetSignMap(requireCompany bool) (map[string]string, error) {
	if err := n.CanSign(requireCompany); err != nil {
		return nil, err
	}
	sign := map[string]string{}
	sign["title"] = strings.TrimSpace(n.Title)
	sign["name"] = n.GetFirstLastName()
	if requireCompany {
		sign["company"] = n.GetCompany()
	}
	return sign, nil
}
