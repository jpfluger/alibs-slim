package aimage

type ImageFilterOption struct {
	Types ImageTypes `json:"types,omitempty"` // List of allowed image types/extensions
	Tags  []string   `json:"tags,omitempty"`  // List of allowed tags
}

func (f *ImageFilterOption) HasOptions() bool {
	if f == nil {
		return false
	}
	if f.Types != nil && len(f.Types) > 0 {
		return true
	}
	if f.Tags != nil && len(f.Tags) > 0 {
		return true
	}
	return false
}
