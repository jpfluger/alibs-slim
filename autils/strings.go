package autils

import "strings"

// ToStringTrimLower returns the input string in lowercase after trimming whitespace.
func ToStringTrimLower(target string) string {
	return strings.ToLower(strings.TrimSpace(target))
}

// ToStringTrimUpper returns the input string in uppercase after trimming whitespace.
func ToStringTrimUpper(target string) string {
	return strings.ToUpper(strings.TrimSpace(target))
}

// HasPrefixPath checks if the target string starts with any of the provided prefixes.
func HasPrefixPath(target string, prefixPaths []string) bool {
	if len(prefixPaths) == 0 || target == "" {
		return false
	}
	for _, prefixPath := range prefixPaths {
		if strings.HasPrefix(target, prefixPath) {
			return true
		}
	}
	return false
}

// ExtractPrefixBrackets extracts content within brackets and the remaining string.
func ExtractPrefixBrackets(target string) (inner string, outer string) {
	return ExtractPlaceholderLeftRight(target, "[", "]")
}

// ExtractPrefixParenthesis extracts content within parentheses and the remaining string.
func ExtractPrefixParenthesis(target string) (inner string, outer string) {
	return ExtractPlaceholderLeftRight(target, "(", ")")
}

// ExtractPrefixBraces extracts content within braces and the remaining string.
func ExtractPrefixBraces(target string) (inner string, outer string) {
	return ExtractPlaceholderLeftRight(target, "{", "}")
}

// ExtractPlaceholderLeftRight extracts content between left and right delimiters and the remaining string.
func ExtractPlaceholderLeftRight(target string, left string, right string) (inner string, outer string) {
	outer = target
	if strings.HasPrefix(target, left) {
		outer = strings.TrimPrefix(outer, left)
		ss := strings.SplitN(outer, right, 2)
		if len(ss) == 2 {
			inner = ss[0]
			outer = ss[1]
		}
	}
	return
}

// MergeMaps merges multiple map[string]string maps.
func MergeMaps(maps ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, m := range maps {
		if m == nil || len(m) == 0 {
			continue
		}
		for k, v := range m {
			if strings.TrimSpace(k) == "" {
				continue
			}
			merged[k] = v
		}
	}
	return merged
}
