package atemplates

import (
	"strings"
	ttemplate "text/template"

	"github.com/dustin/go-humanize"
	"github.com/jpfluger/alibs-slim/atime"
	"github.com/jpfluger/alibs-slim/atypeconvert"
)

// TemplateTextMap is a map that associates string keys with text templates.
type TemplateTextMap map[string]*ttemplate.Template

// Find retrieves a template associated with the given key from the map.
func (ttm TemplateTextMap) Find(key string) *ttemplate.Template {
	if ttm == nil {
		return nil
	}
	return ttm[key]
}

// GetTextTemplateFunctions returns a map of common template functions.
func GetTextTemplateFunctions(fmapType TemplateFunctions) *ttemplate.FuncMap {
	switch fmapType {
	case TEMPLATE_FUNCTIONS_COMMON:
		return &ttemplate.FuncMap{
			// Date and time formatting functions.
			"EnsureDateTime":            atime.EnsureDateTime,
			"FormatDateTime":            atime.FormatDateTime,
			"FormatDateTimeRFC3339":     atime.FormatDateTimeRFC3339,
			"FormatDateTimeRFC3339Nano": atime.FormatDateTimeRFC3339Nano,
			"FormatDateTimeRFC1123":     atime.FormatDateTimeRFC1123,
			"FormatDateTimeRFC1123Z":    atime.FormatDateTimeRFC1123Z,
			"ConvertToTimeZone":         atime.ConvertToTimeZone,
			"IsDateBeforeNow":           atime.IsDateBeforeNow,
			"IsDateAfterNow":            atime.IsDateAfterNow,
			"IsDateBefore":              atime.IsDateBefore,
			"IsDateAfter":               atime.IsDateAfter,
			"FormatDateTimeElse":        atime.FormatDateTimeElse,
			"IfDateEmptyElse":           atime.IfDateEmptyElse,
			"CurrentYear":               atime.CurrentYear,
			"CurrentYearUTC":            atime.CurrentYearUTC,

			// String manipulation functions.
			"ToUpper":      strings.ToUpper,
			"ToLower":      strings.ToLower,
			"ToUpperFirst": ToUpperFirst,
			"Title":        Title,

			// Conversion functions.
			"ToInt":       ToInt,
			"ToInt64":     ToInt64,
			"ToUInt64":    ToUInt64,
			"ToFloat64":   ToFloat64,
			"ToString":    ToString,
			"ToBool":      ToBool,
			"ToUUIDEmpty": ToUUIDEmpty,

			// Logical operations.
			"IfBoolThen":                 IfBoolThen,
			"IfStringNotEmptyElse":       IfStringNotEmptyElse,
			"IfStringNotEmptyFormatElse": IfStringNotEmptyFormatElse,
			"IfStringEmptyThen":          IfStringEmptyThen,
			"IfStringFormatElse":         IfStringFormatElse,
			"IfStringCompareThen":        IfStringCompareThen,
			"IfStringArrContains":        IfStringArrContains,
			"IfUUIDNilElse":              IfUUIDNilElse,
			"IfUUIDCompareThen":          IfUUIDCompareThen,
			"IfIntThen":                  IfIntThen,
			"IfIntCompareThen":           IfIntCompareThen,
			"IfFloatCompareThen":         IfFloatCompareThen,
			"IsNotNil":                   IsNotNil,
			"IsNil":                      IsNil,
			"AddInteger":                 AddInteger,
			"SubtractInteger":            SubtractInteger,
			"UntilInteger":               UntilInteger,

			// Dictionary
			"Dict": Dict,

			// Array
			"JoinString":          JoinString,
			"ArrayContains":       ArrayContains,
			"ArrayContainsString": ArrayContainsString,

			// Formatting functions.
			"FormatIntegerComma":                 FormatIntegerComma,
			"FormatFloatComma":                   FormatFloatComma,
			"FormatFloatCommaDecimal":            FormatFloatCommaDecimal,
			"FormatFloatNoTrailingZeroes":        FormatFloatNoTrailingZeroes,
			"FormatFloatDecimalNoTrailingZeroes": FormatFloatDecimalNoTrailingZeroes,
			"FormatIntegerOrdinal":               FormatIntegerOrdinal,
			"FormatDateTimeRelative":             atime.FormatDateTimeRelative,
			"FormatDateTimeAgo":                  atime.FormatDateTimeAgo,
			"FormatBytes":                        FormatBytes,
		}
	}
	return &ttemplate.FuncMap{}
}

// FormatIntegerComma formats an integer with commas.
func FormatIntegerComma(v int) string {
	return humanize.Comma(int64(v))
}

// FormatBytes formats a byte count into a human-readable string.
func FormatBytes(v int) string {
	return humanize.Bytes(uint64(v))
}

// ToInt converts an interface{} to an int, returning 0 on error.
func ToInt(v interface{}) int {
	val, err := atypeconvert.ConvertToIntFrom(v)
	if err != nil {
		return 0
	}
	return val
}

// ToInt64 converts an interface{} to an int64, returning 0 on error.
func ToInt64(v interface{}) int64 {
	val, err := atypeconvert.ConvertToIntFrom(v)
	if err != nil {
		return 0
	}
	return int64(val)
}

// ToUInt64 converts an interface{} to a uint64, returning 0 on error.
func ToUInt64(v interface{}) uint64 {
	val, err := atypeconvert.ConvertToIntFrom(v)
	if err != nil {
		return 0
	}
	return uint64(val)
}

// ToFloat64 converts an interface{} to a float64, returning 0.0 on error.
func ToFloat64(v interface{}) float64 {
	val, err := atypeconvert.ConvertToFloatFrom(v)
	if err != nil {
		return 0
	}
	return val
}

// ToBool converts an interface{} to a bool, returning false on error.
func ToBool(v interface{}) bool {
	val, err := atypeconvert.ConvertToBoolFrom(v)
	if err != nil {
		return false
	}
	return val
}

// ToString converts an interface{} to a string, returning an empty string on error.
func ToString(v interface{}) string {
	val, err := atypeconvert.ConvertToStringFrom(v)
	if err != nil {
		return ""
	}
	return val
}
