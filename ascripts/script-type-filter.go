package ascripts

// ScriptTypeFilter represents a filter based on script type with an optional title.
type ScriptTypeFilter struct {
	Type  ScriptType `json:"type,omitempty"`  // The script type to filter by
	Title string     `json:"title,omitempty"` // The human-readable title for the filter
}

// ScriptTypeFilters is a slice of pointers to ScriptTypeFilter.
type ScriptTypeFilters []*ScriptTypeFilter

// Find searches for a ScriptTypeFilter in the slice that matches the given ScriptType.
func (stfs ScriptTypeFilters) Find(sType ScriptType) *ScriptTypeFilter {
	// Return nil if the slice is nil, empty, or the ScriptType is empty
	if stfs == nil || len(stfs) == 0 || sType.IsEmpty() {
		return nil
	}
	// Iterate over the filters to find a match
	for _, stf := range stfs {
		if stf.Type.TrimSpaceToLower() == sType.TrimSpaceToLower() {
			return stf // Return the matching filter
		}
	}
	return nil // Return nil if no match is found
}

// Has checks if there is a ScriptTypeFilter in the slice that matches the given ScriptType.
func (stfs ScriptTypeFilters) Has(sType ScriptType) bool {
	return stfs.Find(sType) != nil // Utilize the Find method to check for existence
}

// GetScriptTypeFilterDefaults creates a slice of default ScriptTypeFilters based on the provided ScriptTypes.
func GetScriptTypeFilterDefaults(scriptTypes ScriptTypes) ScriptTypeFilters {
	stfs := ScriptTypeFilters{} // Initialize an empty slice of ScriptTypeFilters
	// Return the empty slice if the provided ScriptTypes is nil or empty
	if scriptTypes == nil || len(scriptTypes) == 0 {
		return stfs
	}
	// Iterate over the provided ScriptTypes to create default filters
	for _, sType := range scriptTypes {
		switch sType {
		case SCRIPTTYPE_GO:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_GO, Title: "Go"})
		case SCRIPTTYPE_HTML:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_HTML, Title: "HTML"})
		case SCRIPTTYPE_MARKDOWN:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_MARKDOWN, Title: "Markdown"})
		case SCRIPTTYPE_MARKDOWN_HTML:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_MARKDOWN_HTML, Title: "Markdown with HTML"})
		case SCRIPTTYPE_CSS:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_CSS, Title: "CSS"})
		case SCRIPTTYPE_JS:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_JS, Title: "Javascript"})
		case SCRIPTTYPE_TEXT:
			stfs = append(stfs, &ScriptTypeFilter{Type: SCRIPTTYPE_TEXT, Title: "Text"})
			// Add cases for other script types if necessary
		}
	}
	return stfs // Return the populated slice of ScriptTypeFilters
}
