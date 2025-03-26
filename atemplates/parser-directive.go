package atemplates

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// StringParserDirectiveFunction defines a function type that takes a key and value,
// and returns a processed key and value.
type StringParserDirectiveFunction func(pkey ParseDirectiveType, pvalue string) (key ParseDirectiveType, value string)

// PARSEDIRECTIVETYPE_FILE is a constant that represents the 'file' directive type.
const PARSEDIRECTIVETYPE_FILE = ParseDirectiveType("file")

// p is a package-level variable that holds the singleton instance of pds.
var p *pds

// pds is a struct that holds a map of parser directive instances and a mutex for concurrent access.
type pds struct {
	instances StringParserDirectiveMap
	mu        sync.RWMutex
}

// init initializes the pds singleton with the 'file' directive function.
func init() {
	p = &pds{
		instances: StringParserDirectiveMap{},
	}
	p.instances[PARSEDIRECTIVETYPE_FILE] = readFile_StringParserDirectiveFunction
}

// StringParserDirectives returns the singleton instance of pds.
func StringParserDirectives() *pds {
	return p
}

// RegisterStringParserDirective registers a new string parser directive function for a given key.
func (p *pds) RegisterStringParserDirective(key ParseDirectiveType, fn StringParserDirectiveFunction) error {
	if key.IsEmpty() {
		return fmt.Errorf("key parameter missing")
	}
	if fn == nil {
		return fmt.Errorf("parser directive function parameter missing")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.instances[key] = fn
	return nil
}

// Get retrieves a registered string parser directive function by key.
func (p *pds) Get(key ParseDirectiveType) StringParserDirectiveFunction {
	if key.IsEmpty() {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.instances[key]
}

// ParseDirective parses a directive from a target string and returns the processed value.
func ParseDirective(target string) string {
	_, value := ParseDirectiveKeyValueWithAutoload(target, true)
	return value
}

// ParseDirectiveKeyValue parses a directive from a target string and returns the key and value.
func ParseDirectiveKeyValue(target string) (key ParseDirectiveType, value string) {
	return ParseDirectiveKeyValueWithAutoload(target, false)
}

// ParseDirectiveKeyValueWithAutoload parses a directive with an option to autoload the directive function.
func ParseDirectiveKeyValueWithAutoload(target string, doAutoLoad bool) (key ParseDirectiveType, value string) {
	tmp := strings.TrimSpace(target)
	if tmp == "" {
		return key, target
	}

	if strings.HasPrefix(tmp, "{{") && strings.HasSuffix(tmp, "}}") {
		tmp = strings.TrimSuffix(strings.TrimPrefix(tmp, "{{"), "}}")
		split := strings.Split(tmp, ":")

		if len(split) < 2 {
			return key, ""
		} else if len(split) > 2 {
			merged := strings.Join(split[1:], ":")
			split = []string{split[0], merged}
		}

		if !doAutoLoad {
			return ParseDirectiveType(split[0]), split[1]
		}

		fn := StringParserDirectives().Get(ParseDirectiveType(split[0]))
		if fn == nil {
			return ParseDirectiveType(split[0]), split[1]
		}

		return fn(ParseDirectiveType(split[0]), split[1])
	}
	return key, target
}

// HasDirective checks if a target string contains a directive.
func HasDirective(target string) bool {
	key, _ := ParseDirectiveKeyValueWithAutoload(target, false)
	return !key.IsEmpty()
}

// readFile_StringParserDirectiveFunction is a directive function that reads a file and returns its content.
func readFile_StringParserDirectiveFunction(pkey ParseDirectiveType, pvalue string) (key ParseDirectiveType, value string) {
	data, err := os.ReadFile(pvalue)
	if err != nil {
		return pkey, ""
	}

	return pkey, strings.TrimRight(string(data), "\n")
}
