package atemplates

import (
	"fmt"
	htemplate "html/template"
	"strings"
	"unicode"

	"github.com/dustin/go-humanize"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/atime"
)

// TemplateHtmlMap is a map that associates string keys with HTML templates.
type TemplateHtmlMap map[string]*htemplate.Template

// Find retrieves a template associated with the given key from the map.
func (ttm TemplateHtmlMap) Find(key string) *htemplate.Template {
	if ttm == nil {
		return nil
	}
	return ttm[key]
}

// GetHTMLTemplateFunctions returns a map of common HTML template functions.
func GetHTMLTemplateFunctions(fmapType TemplateFunctions) *htemplate.FuncMap {
	switch fmapType {
	case TEMPLATE_FUNCTIONS_COMMON:
		// Define common HTML template functions used across various templates.
		return &htemplate.FuncMap{
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
			"ToUpper":        strings.ToUpper,
			"ToLower":        strings.ToLower,
			"ToUpperFirst":   ToUpperFirst,
			"Title":          Title,
			"TruncateString": TruncateString,
			"SliceString":    SliceString,

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

			// Dictionary
			"Dict": Dict,

			// Array
			"JoinString":          JoinString,
			"ArrayContains":       ArrayContains,
			"ArrayContainsString": ArrayContainsString,

			// HTML-specific functions to ensure safe rendering.
			"SafeURL":      SafeURL,
			"SafeHtml":     SafeHtml,
			"SafeHtmlAttr": SafeHtmlAttr,
			"SafeJS":       SafeJS,
			"SafeCSS":      SafeCSS,
			// Deprecated
			"MustSnippetRenderHTML": MustSnippetRenderHTML,
			// Deprecated
			"MustSnippetRenderText": MustSnippetRenderText,
			"MustRenderHTML":        MustFSSnippetRenderHTML,
			"MustRenderText":        MustFSSnippetRenderText,
		}
	default:
		// Return an empty function map if no common functions are requested.
		return &htemplate.FuncMap{}
	}
}

// SafeURL marks a string as a safe URL within an HTML template.
func SafeURL(s string) htemplate.URL {
	return htemplate.URL(s)
}

// SafeHtml marks a string as safe HTML content within an HTML template.
func SafeHtml(s string) htemplate.HTML {
	return htemplate.HTML(s)
}

// SafeHtmlAttr marks a string as a safe HTML attribute within an HTML template.
func SafeHtmlAttr(s string) htemplate.HTMLAttr {
	return htemplate.HTMLAttr(s)
}

// SafeJS marks a string as safe JavaScript within an HTML template.
func SafeJS(s string) htemplate.JS {
	return htemplate.JS(s)
}

// SafeCSS marks a string as safe CSS within an HTML template.
func SafeCSS(s string) htemplate.CSS {
	return htemplate.CSS(s)
}

// IfBoolThen returns one of two strings based on a boolean condition.
func IfBoolThen(target bool, thenString string, elseString string) string {
	if target {
		return thenString
	}
	return elseString
}

// IfStringNotEmptyElse returns the target string if it's not empty, otherwise returns the elseString.
func IfStringNotEmptyElse(target string, elseString string) string {
	if strings.TrimSpace(target) == "" {
		return elseString
	}
	return target
}

// IfStringEmptyThen returns trueValue if the target string is empty, otherwise returns falseValue.
func IfStringEmptyThen(target string, trueValue string, falseValue string) string {
	if strings.TrimSpace(target) != "" {
		return trueValue
	}
	return falseValue
}

// IfStringNotEmptyFormatElse formats the target string if it's not empty, otherwise returns falseValue.
func IfStringNotEmptyFormatElse(target string, format string, falseValue string) string {
	target = strings.TrimSpace(target)
	if target != "" {
		if format == "" {
			return target
		}
		return fmt.Sprintf(format, target)
	}
	return falseValue
}

// IfStringFormatElse formats the target string if it's not empty, otherwise returns falseValue.
func IfStringFormatElse(target string, format string, falseValue string) string {
	if target != "" {
		if format == "" {
			return target
		}
		return fmt.Sprintf(format, target)
	}
	return falseValue
}

// IfStringCompareThen compares two strings and returns trueValue or falseValue based on the comparison result.
func IfStringCompareThen(a string, b string, operator string, trueValue string, falseValue string) string {
	switch operator {
	case "=", "==":
		if a == b {
			return trueValue
		}
	case "hasPrefix":
		if strings.HasPrefix(a, b) {
			return trueValue
		}
	case "hasSuffix":
		if strings.HasSuffix(a, b) {
			return trueValue
		}
	case "!=":
		if a != b {
			return trueValue
		}
	}
	return falseValue
}

// IfStringArrContains checks if a slice contains a specific string
func IfStringArrContains(slice []string, item string) bool {
	if slice == nil {
		return false
	}
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

// ToUUIDEmpty converts a UUID to a string, returning an empty string if the UUID is nil.
func ToUUIDEmpty(target uuid.UUID) string {
	if target == uuid.Nil {
		return ""
	}
	return target.String()
}

// IfUUIDNilElse checks if a UUID is nil and returns elseValue or the UUID string accordingly.
func IfUUIDNilElse(target uuid.UUID, elseValue string) string {
	if target == uuid.Nil {
		return elseValue
	}
	return target.String()
}

// IfUUIDCompareThen compares two UUIDs and returns trueValue or falseValue based on the comparison result.
func IfUUIDCompareThen(a uuid.UUID, b uuid.UUID, operator string, trueValue string, falseValue string) string {
	switch operator {
	case "=", "==":
		if a == b {
			return trueValue
		}
	case "!=":
		if a != b {
			return trueValue
		}
	}
	return falseValue
}

// IfIntThen returns one of two strings based on whether the target integer is non-zero.
func IfIntThen(target int, trueValue string, falseValue string) string {
	if target != 0 {
		return trueValue
	}
	return falseValue
}

// IfIntCompareThen compares two integers and returns trueValue or falseValue based on the comparison result.
func IfIntCompareThen(a int, b int, operator string, trueValue string, falseValue string) string {
	switch operator {
	case "=", "==":
		return IfBoolThen(a == b, trueValue, falseValue)
	case ">":
		if a > b {
			return trueValue
		}
	case ">=":
		if a >= b {
			return trueValue
		}
	case "<":
		if a < b {
			return trueValue
		}
	case "<=":
		if a <= b {
			return trueValue
		}
	case "!=":
		if a != b {
			return trueValue
		}
	}
	return falseValue
}

// IfFloatCompareThen compares two floating-point numbers and returns trueValue or falseValue based on the comparison result.
func IfFloatCompareThen(a float64, b float64, operator string, trueValue string, falseValue string) string {
	switch operator {
	case "=", "==":
		return IfBoolThen(a == b, trueValue, falseValue)
	case ">":
		if a > b {
			return trueValue
		}
	case ">=":
		if a >= b {
			return trueValue
		}
	case "<":
		if a < b {
			return trueValue
		}
	case "<=":
		if a <= b {
			return trueValue
		}
	case "!=":
		if a != b {
			return trueValue
		}
	}
	return falseValue
}

// IsNotNil checks if the target interface is not nil.
func IsNotNil(target interface{}) bool {
	return target != nil
}

// IsNil checks if the target interface is nil.
func IsNil(target interface{}) bool {
	return target == nil
}

// AddInteger adds two integers and returns the result.
func AddInteger(a int, b int) int {
	return a + b
}

// UntilInteger returns a slice of integers from 0 up to (but not including) n.
// It is commonly used to simulate a basic for-loop in templates.
//
// For example, UntilInteger(3) returns []int{0, 1, 2}, which is useful
// for iterating with index values inside a Go HTML template.
func UntilInteger(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

// SubtractInteger subtracts one integer from another and returns the result.
func SubtractInteger(a int, b int) int {
	return a - b
}

// FormatFloatComma formats a float64 with commas.
func FormatFloatComma(v float64) string {
	return humanize.Commaf(v)
}

// FormatFloatCommaDecimal formats a float64 with commas and a specified number of decimal places.
func FormatFloatCommaDecimal(f float64, decimals int) string {
	return humanize.CommafWithDigits(f, decimals)
}

// FormatFloatNoTrailingZeroes formats a float64 without trailing zeroes.
func FormatFloatNoTrailingZeroes(num float64) string {
	return humanize.Ftoa(num)
}

// FormatFloatDecimalNoTrailingZeroes formats a float64 without trailing zeroes and with a specified number of decimal places.
func FormatFloatDecimalNoTrailingZeroes(num float64, digits int) string {
	return humanize.FtoaWithDigits(num, digits)
}

// FormatIntegerOrdinal formats an integer with its ordinal representation.
func FormatIntegerOrdinal(x int) string {
	return humanize.Ordinal(x)
}

// ToUpperFirst capitalizes the first letter of a string.
func ToUpperFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

// Dict function to create a dictionary (map[string]interface{})
func Dict(values ...interface{}) map[string]interface{} {
	d := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("dict function keys must be strings") // Panic if wrong key type (optional)
		}
		d[key] = values[i+1]
	}
	return d
}

// JoinString, which works just like strings and is safe for use in Go HTML templates.
func JoinString(items []string, sep string) string {
	if items == nil || len(items) == 0 {
		return ""
	}
	return strings.Join(items, sep)
}

// ArrayContains checks if an array of interface{} contains the target value.
func ArrayContains(arr []interface{}, target interface{}) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// ArrayContainsString checks if an array of string contains the target value.
func ArrayContainsString(arr []string, target string) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// Title capitalizes the first letter of the input string using Unicode-aware handling.
func Title(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

// TruncateString conditionally truncates a string to a maximum character length
// and appends an optional suffix (e.g., "..."). It handles Unicode characters correctly
// by operating on runes. If the string length is already less than or equal to
// the maximum length, the original string is returned.
//
// Usage in template: {{ TruncateString .Notes 20 }} or {{ TruncateString .Notes 20 "..." }}
func TruncateString(s string, length int, suffix ...string) string {
	r := []rune(s)
	rLen := len(r)

	// 1. Check if truncation is necessary based on character count (runes).
	if rLen <= length {
		return s
	}

	// 2. Determine the suffix to use.
	// If a suffix is provided in the variadic argument list, use the first one.
	// Otherwise, use an empty string.
	var finalSuffix string
	if len(suffix) > 0 {
		finalSuffix = suffix[0]
	}

	// 3. Truncate the rune slice and convert back to string, appending the suffix.
	// We slice up to 'length' characters.
	return string(r[:length]) + finalSuffix
}

// SliceString returns a substring of s from a start index up to an end index.
// It handles Unicode characters correctly by operating on runes.
// If start or end indices are out of bounds, they are adjusted.
func SliceString(s string, start int, end int) string {
	r := []rune(s)
	rLen := len(r)

	// Adjust start index
	if start < 0 {
		start = 0
	}
	// Adjust end index
	if end > rLen {
		end = rLen
	}

	// Check if the slice is valid after adjustments
	if start >= end || start > rLen {
		return ""
	}

	return string(r[start:end])
}
