package ascripts

import (
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/jpfluger/alibs-slim/ajson"
	"github.com/jpfluger/alibs-slim/areflect"
	"github.com/jpfluger/alibs-slim/autils"
)

// Scriptable represents a script with its type and body content.
type Scriptable struct {
	Key      ajson.JsonKey `json:"key,omitempty"`  // Unique identifier for the script
	Type     ScriptType    `json:"type,omitempty"` // Type of the script (e.g., Go, HTML)
	Body     string        `json:"body,omitempty"` // The script content
	cache    string        // Cached result of the script execution
	compiler ICompiler     // Compiler interface to compile the script
	mu       sync.RWMutex  // Mutex to protect concurrent access
}

// NewScriptableFromPath creates a new Scriptable from a file path.
func NewScriptableFromPath(filePath string) (*Scriptable, error) {
	return NewScriptableFromPathWithType(filePath, "")
}

// NewScriptableFromPathWithType creates a new Scriptable from a file path with a specified script type.
func NewScriptableFromPathWithType(filePath string, scriptType ScriptType) (*Scriptable, error) {
	return NewScriptableFromPathWithTypeByDirOption(filePath, scriptType, "")
}

// NewScriptableFromPathWithTypeByDirOption creates a new Scriptable from a file path with a specified script type and directory option.
func NewScriptableFromPathWithTypeByDirOption(filePath string, scriptType ScriptType, dirOption string) (*Scriptable, error) {
	b, err := autils.ReadByFilePathWithDirOption(filePath, dirOption)
	if err != nil {
		return nil, err
	}

	filePath = strings.TrimSpace(filePath)

	if scriptType.IsEmpty() {
		_, _, ext := autils.FileNameParts(filePath)
		scriptType = ExtToScriptType(ext)
	}

	return &Scriptable{
		Key:  ajson.JsonKey(path.Base(filePath)),
		Type: scriptType,
		Body: string(b),
	}, nil
}

// GetKey returns the script's key.
func (s *Scriptable) GetKey() ajson.JsonKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Key
}

// GetType returns the script's type.
func (s *Scriptable) GetType() ScriptType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Type
}

// SetType sets the script's type and resets the compiler.
func (s *Scriptable) SetType(scriptType ScriptType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Type = scriptType
	s.compiler = nil // Invalidate the existing compiler
}

// GetBody returns the script's body content.
func (s *Scriptable) GetBody() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Body
}

// CanCompile checks if the script has the necessary properties to be compiled.
func (s *Scriptable) CanCompile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s != nil && strings.TrimSpace(s.Body) != "" && !s.Type.IsEmpty()
}

// GetCache returns the cached result of the script execution.
func (s *Scriptable) GetCache() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cache
}

// SetCache sets the cached result of the script execution.
func (s *Scriptable) SetCache(cache string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = cache
}

// GetCompiler returns the current compiler.
func (s *Scriptable) GetCompiler() ICompiler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.compiler
}

// NewCompiler creates a new compiler based on the script's type.
func (s *Scriptable) NewCompiler() (ICompiler, error) {
	rtype, err := areflect.TypeManager().FindReflectType(TYPEMANAGER_SCRIPT_COMPILER, s.GetType().String())
	if err != nil {
		return nil, fmt.Errorf("cannot find script compiler for type '%s': %v", s.GetType().String(), err)
	}

	obj, ok := reflect.New(rtype).Interface().(ICompiler)
	if !ok {
		return nil, fmt.Errorf("created object is not of type ICompiler")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.compiler = obj // Set the new compiler

	return s.compiler, nil
}

// Compile compiles the script with the given parameters.
func (s *Scriptable) Compile(params interface{}) (interface{}, error) {
	if !s.CanCompile() {
		return nil, fmt.Errorf("scriptable is missing compilable properties")
	}

	if s.GetCompiler() == nil {
		if _, err := s.NewCompiler(); err != nil {
			return nil, fmt.Errorf("failed to create new compiler: %v", err)
		}
	}

	return s.GetCompiler().Run(s.GetBody(), params)
}

// Render renders the script with the given data.
func (s *Scriptable) Render(data interface{}) (string, error) {
	if !s.CanCompile() {
		return "", fmt.Errorf("scriptable is missing compilable properties")
	}

	if s.GetCompiler() == nil {
		if _, err := s.NewCompiler(); err != nil {
			return "", fmt.Errorf("failed to create new compiler: %v", err)
		}
	}

	return s.GetCompiler().Render(s.GetBody(), data)
}

// RenderWithCache renders the script with the given data, using the cache if available and not refreshing.
func (s *Scriptable) RenderWithCache(data interface{}, refreshCache bool) (string, error) {
	if !refreshCache {
		cache := s.GetCache()
		if strings.TrimSpace(cache) != "" {
			return cache, nil
		}
	}

	content, err := s.Render(data)
	if err != nil {
		return "", err
	}

	s.SetCache(content) // Update the cache with the new content
	return content, nil
}
