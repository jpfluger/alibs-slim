package atags

type Tags []TagKey

func (tags *Tags) Add(tag TagKey) {
	for _, t := range *tags {
		if t == tag {
			return
		}
	}
	*tags = append(*tags, tag)
}

func (tags Tags) HasTag(tag TagKey) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (tags *Tags) Remove(tag TagKey) {
	newTags := Tags{}
	for _, t := range *tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	*tags = newTags
}

func (tags Tags) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for _, t := range tags {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}
