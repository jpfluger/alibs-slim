package atags

//type Taggable interface {
//	MarshalTag() (string, error)
//	UnmarshalTag(string) error
//}
//
//type StringTag struct {
//	Value string
//}
//
//func (s *StringTag) MarshalTag() (string, error) {
//	return s.Value, nil
//}
//
//func (s *StringTag) UnmarshalTag(data string) error {
//	s.Value = data
//	return nil
//}
//
//func (s *StringTag) ToInt64() (int64, error) {
//	// Conversion logic
//}
//
//type TaggingSystem struct {
//	tags map[string]Taggable
//}
//
//func NewTaggingSystem() *TaggingSystem {
//	return &TaggingSystem{
//		tags: make(map[string]Taggable),
//	}
//}
//
//func (t *TaggingSystem) SetTag(key string, tag Taggable) {
//	t.tags[key] = tag
//}
//
//func (t *TaggingSystem) GetTag(key string) Taggable {
//	return t.tags[key]
//}
