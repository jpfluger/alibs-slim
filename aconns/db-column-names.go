package aconns

// DbColumnNames is a slice of DbColumnName that represents a list of database column names.
type DbColumnNames []DbColumnName

// NewDbColumnNames creates a new DbColumnNames slice from a variadic list of DbColumnName.
func NewDbColumnNames(names ...DbColumnName) DbColumnNames {
	return names
}

// Find searches for a DbColumnName in the slice and returns it if found.
func (cns DbColumnNames) Find(name DbColumnName) DbColumnName {
	if cns == nil || len(cns) == 0 || name.IsEmpty() {
		return ""
	}
	for _, cn := range cns {
		if cn == name {
			return cn
		}
	}
	return ""
}

// Has checks if a DbColumnName exists in the slice.
func (cns DbColumnNames) Has(name DbColumnName) bool {
	return cns.Find(name) != ""
}

// Length returns the number of DbColumnName in the slice.
func (cns DbColumnNames) Length() int {
	if cns == nil {
		return 0
	}
	return len(cns)
}

// Ensure ensures that the given DbColumnNames are present in the slice.
// It's an alias for the Add method for better readability.
func (cns DbColumnNames) Ensure(names ...DbColumnName) DbColumnNames {
	return cns.Add(names...)
}

// Add appends non-duplicate DbColumnNames to the slice.
func (cns DbColumnNames) Add(names ...DbColumnName) DbColumnNames {
	if names == nil || len(names) == 0 {
		return cns
	}
	cnsNew := make(DbColumnNames, len(cns))
	copy(cnsNew, cns)
	for _, name := range names {
		if !cnsNew.Has(name) {
			cnsNew = append(cnsNew, name)
		}
	}
	return cnsNew
}

// Remove deletes the specified DbColumnNames from the slice.
func (cns DbColumnNames) Remove(names ...DbColumnName) DbColumnNames {
	if names == nil || len(names) == 0 {
		return cns
	}
	cnsNew := DbColumnNames{}
	for _, cn := range cns {
		remove := false
		for _, name := range names {
			if cn == name {
				remove = true
				break
			}
		}
		if !remove {
			cnsNew = append(cnsNew, cn)
		}
	}
	return cnsNew
}
