package acontact

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/anetwork" // Custom network utilities
	"github.com/jpfluger/alibs-slim/autils"   // Custom utility functions
	"sort"
	"strings"
)

// Url represents a URL with its type, title, link, and default status.
type Url struct {
	Type      UrlType          `json:"type,omitempty"`      // The type of URL (e.g., home, work)
	Title     string           `json:"title,omitempty"`     // The title of the URL
	Link      *anetwork.NetURL `json:"link,omitempty"`      // The actual URL link
	IsDefault bool             `json:"isDefault,omitempty"` // Indicates if this is the default URL
}

// Validate checks if the Url fields are valid.
func (u *Url) Validate() error {
	if u.Type.IsEmpty() {
		return fmt.Errorf("url type is empty")
	}
	u.Title = strings.TrimSpace(u.Title) // Trim spaces from the title
	if u.Link == nil {
		return fmt.Errorf("url link is empty")
	}
	if !u.Link.IsUrl() {
		return fmt.Errorf("url link is not a valid URL")
	}
	return nil
}

// GetLinkWithOptions returns specific parts of the URL based on the input parameter.
// If no option is specified, the full Link is returned.
func (u *Url) GetLinkWithOptions(part string) string {
	if u.Link == nil {
		return ""
	}
	if part == "raw" {
		return u.Link.String()
	}
	if !u.Link.IsUrl() {
		return ""
	}

	// Return specific parts of the URL based on the requested part
	switch part {
	case "domain":
		return u.Link.Host
	case "path":
		return u.Link.Path
	case "port":
		return u.Link.Port()
	case "scheme":
		return u.Link.Scheme
	case "no-scheme":
		scheme := u.Link.Scheme
		if !strings.Contains(scheme, "://") {
			scheme += "://"
		}
		return strings.TrimPrefix(u.Link.String(), scheme)
	default:
		return u.Link.String()
	}
}

// Urls is a slice of Url pointers, representing a collection of URLs.
type Urls []*Url

// FindByType searches for a URL by its type.
func (us Urls) FindByType(urlType UrlType) *Url {
	return us.findByType(urlType, false)
}

// FindByTypeOrDefault searches for a URL by its type or returns the default URL.
func (us Urls) FindByTypeOrDefault(urlType UrlType) *Url {
	return us.findByType(urlType, true)
}

// findByType is a helper function that searches for a URL by type and optionally returns the default URL.
func (us Urls) findByType(urlType UrlType, checkDefault bool) *Url {
	var defaultUrl *Url
	for _, u := range us {
		if u.Type.ToStringTrimLower() == urlType.ToStringTrimLower() {
			return u
		}
		if u.IsDefault {
			defaultUrl = u
		}
	}
	if checkDefault {
		return defaultUrl
	}
	return nil
}

// FindByLink searches for a URL by its address.
func (us Urls) FindByLink(address string, checkValidOnly bool) *Url {
	address = autils.ToStringTrimLower(address)
	if address == "" {
		return nil
	}
	for _, u := range us {
		if u.Link != nil {
			if checkValidOnly && !u.Link.IsUrl() {
				continue
			}
			if autils.ToStringTrimLower(u.Link.String()) == address {
				return u
			}
		}
	}
	return nil
}

// HasType checks if a URL of the specified type exists in the collection.
func (us Urls) HasType(urlType UrlType) bool {
	return us.FindByType(urlType) != nil
}

// HasTypeWithDefault checks if a url of the specified type exists, or if there's a default url.
func (us Urls) HasTypeWithDefault(urlType UrlType, allowDefault bool) bool {
	return us.findByType(urlType, allowDefault) != nil
}

// HasTypeOrDefault checks if a URL of the specified type exists, or if there's a default URL.
func (us Urls) HasTypeOrDefault(urlType UrlType) bool {
	return us.FindByTypeOrDefault(urlType) != nil
}

// HasLink checks if a URL with the specified address exists in the collection.
func (us Urls) HasLink(address string) bool {
	return us.FindByLink(address, true) != nil
}

// Clone creates a deep copy of the Urls collection.
func (us Urls) Clone() Urls {
	b, err := json.Marshal(us)
	if err != nil {
		return nil
	}
	var clone Urls
	if err := json.Unmarshal(b, &clone); err != nil {
		return nil
	}
	return clone
}

// MergeFrom adds URLs from another collection that are not already present.
func (us *Urls) MergeFrom(target Urls) {
	if us == nil || target == nil {
		return
	}
	for _, t := range target {
		if t.Type.IsEmpty() {
			continue
		}
		isFound := false
		for _, u := range *us {
			if u.Type.ToStringTrimLower() == t.Type.ToStringTrimLower() {
				isFound = true
				break
			}
		}
		if !isFound {
			*us = append(*us, t)
		}
	}
}

// Set adds or updates a URL in the collection.
func (us *Urls) Set(url *Url) {
	if url == nil || url.Type.IsEmpty() || url.Link == nil || !url.Link.IsUrl() {
		return
	}
	// Create a new slice for the updated URLs
	newUrls := Urls{}
	for _, u := range *us {
		if u.Type.ToStringTrimLower() == url.Type.ToStringTrimLower() {
			continue // Skip the URL of the same type to replace it
		} else if u.IsDefault && url.IsDefault {
			u.IsDefault = false // Unset the default if the new URL is the default
		}
		newUrls = append(newUrls, u)
	}
	newUrls = append(newUrls, url) // Add the new URL

	// Sort the URLs, placing the default URL at the top
	sort.SliceStable(newUrls, func(i, j int) bool {
		return newUrls[i].IsDefault || newUrls[i].Type < newUrls[j].Type
	})

	*us = newUrls // Update the original collection
}

// Remove deletes a URL of the specified type from the collection.
func (us *Urls) Remove(urlType UrlType) {
	if urlType.IsEmpty() {
		return
	}
	newUrls := Urls{}
	for _, u := range *us {
		if u.Type.ToStringTrimLower() == urlType.ToStringTrimLower() {
			continue // Skip the URL of the type to be removed
		}
		newUrls = append(newUrls, u)
	}
	*us = newUrls // Update the original collection with the remaining URLs
}
