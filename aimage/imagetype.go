package aimage

import (
	"fmt"
	"strings"
	"unicode"
)

// ImageType is a structured string that
// 1. Allow alpha A-Z, a-z, numbers 0-9, special characters (-_.)
// 2. Only alphas and numbers can begin or end the string
// 3. The period is the divider
type ImageType string

func (rt ImageType) IsEmpty() bool {
	rtNew := strings.TrimSpace(string(rt))
	return rtNew == ""
}

func (rt ImageType) TrimSpace() ImageType {
	rtNew := strings.TrimSpace(string(rt))
	return ImageType(rtNew)
}

func (rt ImageType) String() string {
	return string(rt)
}

func (rt ImageType) HasMatch(rtType ImageType) bool {
	return rt == rtType
}

func (rt ImageType) MatchesOne(rtTypes ...ImageType) bool {
	for _, rtType := range rtTypes {
		if rt == rtType {
			return true
		}
	}
	return false
}

func (rt ImageType) HasPrefix(rtType ImageType) bool {
	return strings.HasPrefix(rt.String(), rtType.String())
}

func (rt ImageType) Validate() error {
	// 1. Allow alpha A-Z, a-z, numbers 0-9, special characters (-_)
	// 2. Only alphas and numbers can begin or end the string
	// 3. The period is the divider
	target := rt.String()
	for ii, r := range target {
		if !unicode.IsLetter(r) {
			if !unicode.IsDigit(r) {
				if string(r) != "-" {
					if string(r) != "_" {
						if string(r) != "." {
							return fmt.Errorf("invalid char '%s'; hash allows alpha A-Z, a-z, numbers 0-9, special '-_.'", string(r))
						}
					}
				}
			}
		}
		if ii == 0 || ii == len(target)-1 {
			for _, hit := range []string{"-", "_", "."} {
				if string(r) == hit {
					return fmt.Errorf("invalid char '%s' at position %d; only alphas and numbers can begin or end the hash", string(r), ii)
				}
			}
		}
	}
	return nil
}

type ImageTypes []ImageType

func (rts ImageTypes) HasValues() bool {
	return rts != nil && len(rts) > 0
}

func (rts ImageTypes) HasMatch(rType ImageType) bool {
	if rts == nil || len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt == rType {
			return true
		}
	}
	return false
}

func (rts ImageTypes) HasPrefix(rType ImageType) bool {
	if rts == nil || len(rts) == 0 || rType.IsEmpty() {
		return false
	}
	for _, rt := range rts {
		if rt.HasPrefix(rType) {
			return true
		}
	}
	return false
}

func (rts ImageTypes) Clone() ImageTypes {
	arr := ImageTypes{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt)
		}
	}
	return arr
}

func (rts ImageTypes) ToArrStrings() []string {
	arr := []string{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if !rt.IsEmpty() {
			arr = append(arr, rt.String())
		}
	}
	return arr
}

func (rts ImageTypes) IncludeIfInTargets(targets ImageTypes) ImageTypes {
	arr := ImageTypes{}
	if rts == nil || len(rts) == 0 || targets == nil || len(targets) == 0 {
		return arr
	}
	for _, rt := range rts {
		if targets.HasMatch(rt) {
			arr = append(arr, rt)
		}
	}
	return arr
}

func (rts ImageTypes) Clean() ImageTypes {
	arr := ImageTypes{}
	if rts == nil || len(rts) == 0 {
		return arr
	}
	for _, rt := range rts {
		if rt.IsEmpty() {
			continue
		}
		arr = append(arr, rt)
	}
	return arr
}
