package aimage

import (
	"strings"
)

// Images is a slice of pointers to Image objects.
type Images []*Image

// Count returns the number of images in the slice.
func (ims Images) Count() int {
	return len(ims)
}

// Clean returns a new slice of Images with all elements that have data.
func (ims Images) Clean() Images {
	var imsNew Images
	for _, im := range ims {
		if im.HasData() {
			imsNew = append(imsNew, im)
		}
	}
	return imsNew
}

// FindByName searches for an image by its name and returns it if found.
func (ims Images) FindByName(name string) *Image {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, im := range ims {
		if strings.ToLower(strings.TrimSpace(im.Name)) == name {
			return im
		}
	}
	return nil
}

// LoadFileImports loads image data from files for all images in the slice.
func (ims Images) LoadFileImports(dirOptions string) error {
	for _, im := range ims {
		if err := im.LoadFileImports(dirOptions); err != nil {
			return err
		}
	}
	return nil
}
