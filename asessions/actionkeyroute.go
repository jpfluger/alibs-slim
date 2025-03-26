package asessions

// ActionKeyUrl represents a mapping between an ActionKey and a URL.
// It includes a flag to indicate if prefix checking is required.
type ActionKeyUrl struct {
	Key         ActionKey `json:"key,omitempty"`         // The action key.
	Url         string    `json:"url,omitempty"`         // The associated URL.
	CheckPrefix bool      `json:"checkPrefix,omitempty"` // Flag to check URL prefix.
}

// ActionKeyUrls is a slice of pointers to ActionKeyUrl.
type ActionKeyUrls []*ActionKeyUrl

// Find locates an ActionKeyUrl in the slice by its key.
// Returns nil if the key is empty or not found.
func (aks ActionKeyUrls) Find(key ActionKey) *ActionKeyUrl {
	if key.IsEmpty() {
		return nil // Return nil if the key is empty.
	}
	for _, ak := range aks {
		if ak.Key == key {
			return ak // Return the matching ActionKeyUrl.
		}
	}
	return nil // Return nil if no match is found.
}

// Has checks if an ActionKeyUrl with the specified key exists in the slice.
func (aks ActionKeyUrls) Has(key ActionKey) bool {
	return aks.Find(key) != nil // Utilize Find method to check existence.
}

// Add inserts a new ActionKeyUrl into the slice if it doesn't already exist.
func (aks *ActionKeyUrls) Add(ak *ActionKeyUrl) {
	if ak == nil || ak.Key.IsEmpty() {
		return // Do nothing if the input is nil or the key is empty.
	}
	if (*aks).Has(ak.Key) {
		return // Do nothing if the key already exists.
	}
	*aks = append(*aks, ak) // Append the new ActionKeyUrl to the slice.
}

// Remove deletes ActionKeyUrls from the slice based on a list of keys.
func (aks ActionKeyUrls) Remove(keys ...ActionKey) ActionKeyUrls {
	if keys == nil || len(keys) == 0 {
		return aks // Return the original slice if no keys are provided.
	}
	newActions := ActionKeyUrls{} // Create a new slice for the remaining ActionKeyUrls.
	for _, ak := range aks {
		if ak.Key.IsEmpty() {
			continue // Skip if the ActionKeyUrl's key is empty.
		}

		isFound := false
		for _, key := range keys {
			if key.IsEmpty() {
				continue // Skip if the key to remove is empty.
			}
			if ak.Key == key {
				isFound = true
				break // Mark as found and stop checking further.
			}
		}

		if !isFound {
			newActions = append(newActions, ak) // Add to new slice if not marked for removal.
		}
	}
	return newActions // Return the new slice with the specified keys removed.
}
