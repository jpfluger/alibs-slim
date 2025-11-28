package amidware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SerializableCORSConfig is a JSON-serializable version of middleware.CORSConfig, omitting function fields.
type SerializableCORSConfig struct {
	AllowOrigins                             []string `json:"allowOrigins,omitempty"`
	AllowMethods                             []string `json:"allowMethods,omitempty"`
	AllowHeaders                             []string `json:"allowHeaders,omitempty"`
	AllowCredentials                         bool     `json:"allowCredentials,omitempty"`
	UnsafeWildcardOriginWithAllowCredentials bool     `json:"unsafeWildcardOriginWithAllowCredentials,omitempty"`
	ExposeHeaders                            []string `json:"exposeHeaders,omitempty"`
	MaxAge                                   int      `json:"maxAge,omitempty"`
}

// toSerializable converts middleware.CORSConfig to SerializableCORSConfig.
func toSerializable(cc middleware.CORSConfig) SerializableCORSConfig {
	return SerializableCORSConfig{
		AllowOrigins:                             cc.AllowOrigins,
		AllowMethods:                             cc.AllowMethods,
		AllowHeaders:                             cc.AllowHeaders,
		AllowCredentials:                         cc.AllowCredentials,
		UnsafeWildcardOriginWithAllowCredentials: cc.UnsafeWildcardOriginWithAllowCredentials,
		ExposeHeaders:                            cc.ExposeHeaders,
		MaxAge:                                   cc.MaxAge,
	}
}

// fromSerializable converts SerializableCORSConfig back to middleware.CORSConfig (function fields remain default).
func fromSerializable(s SerializableCORSConfig) middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:                             s.AllowOrigins,
		AllowMethods:                             s.AllowMethods,
		AllowHeaders:                             s.AllowHeaders,
		AllowCredentials:                         s.AllowCredentials,
		UnsafeWildcardOriginWithAllowCredentials: s.UnsafeWildcardOriginWithAllowCredentials,
		ExposeHeaders:                            s.ExposeHeaders,
		MaxAge:                                   s.MaxAge,
		// Skipper and AllowOriginFunc are omitted and will use defaults (e.g., nil or DefaultSkipper).
	}
}

// CORSConfig holds CORS-related global config.
type CORSConfig struct {
	// Embedded CORSConfig for direct configuration.
	middleware.CORSConfig `json:"cors,omitempty"`

	// Per-origin overrides (key: origin URL, value: custom CORSConfig).
	PerOriginConfigs map[string]middleware.CORSConfig `json:"perOriginConfigs,omitempty"`

	IsPermissive bool `json:"isPermissive"`

	mu sync.RWMutex // Instance-level mutex
}

// MarshalJSON customizes JSON encoding for CORSConfig, omitting non-serializable fields.
func (g *CORSConfig) MarshalJSON() ([]byte, error) {
	aux := struct {
		CORS             SerializableCORSConfig            `json:"cors,omitempty"`
		PerOriginConfigs map[string]SerializableCORSConfig `json:"perOriginConfigs,omitempty"`
		IsPermissive     bool                              `json:"isPermissive"`
	}{
		CORS:         toSerializable(g.CORSConfig),
		IsPermissive: g.IsPermissive,
	}
	aux.PerOriginConfigs = make(map[string]SerializableCORSConfig, len(g.PerOriginConfigs))
	for k, v := range g.PerOriginConfigs {
		aux.PerOriginConfigs[k] = toSerializable(v)
	}
	return json.Marshal(aux)
}

// UnmarshalJSON customizes JSON decoding for CORSConfig.
func (g *CORSConfig) UnmarshalJSON(data []byte) error {
	aux := struct {
		CORS             SerializableCORSConfig            `json:"cors,omitempty"`
		PerOriginConfigs map[string]SerializableCORSConfig `json:"perOriginConfigs,omitempty"`
		IsPermissive     bool                              `json:"isPermissive"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	g.IsPermissive = aux.IsPermissive
	g.CORSConfig = fromSerializable(aux.CORS)
	g.PerOriginConfigs = make(map[string]middleware.CORSConfig, len(aux.PerOriginConfigs))
	for k, v := range aux.PerOriginConfigs {
		g.PerOriginConfigs[k] = fromSerializable(v)
	}
	return nil
}

// NewCORSConfig creates and initializes a CORSConfig instance with optional initial config.
func NewCORSConfig(initialConfig *middleware.CORSConfig) (*CORSConfig, error) {
	gc := &CORSConfig{}
	if initialConfig != nil {
		gc.CORSConfig = *initialConfig
	}
	// No immediate validation; defer to Validate()
	return gc, nil
}

// Validate sets defaults if permissive or using defaultPublicUrl, and validates the CORS config.
func (gc *CORSConfig) Validate(defaultPublicUrl anetwork.NetURL) error {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	// Apply default origin if not permissive and defaultPublicUrl is valid
	if defaultPublicUrl.IsUrl() {
		if !gc.IsPermissive {
			if len(gc.AllowOrigins) == 0 {
				gc.AllowOrigins = []string{defaultPublicUrl.String()}
			}
			if len(gc.AllowMethods) == 0 {
				gc.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} // Adjust as needed.
			}
		}
	}

	// Apply to main config
	if err := gc.validateSingleConfig(&gc.CORSConfig); err != nil {
		return err
	}

	// Apply to per-origin configs
	for origin, config := range gc.PerOriginConfigs {
		if err := gc.validateSingleConfig(&config); err != nil {
			return fmt.Errorf("per-origin config for %s: %v", origin, err)
		}
		gc.PerOriginConfigs[origin] = config // Update if defaults applied
	}

	return nil
}

// validateSingleConfig validates and applies defaults to a single CORSConfig.
func (gc *CORSConfig) validateSingleConfig(config *middleware.CORSConfig) error {
	if gc.IsPermissive {
		// Apply permissive defaults if fields are empty (skip AllowHeaders as per note)
		if len(config.AllowOrigins) == 0 {
			config.AllowOrigins = []string{"*"}
		}
		if len(config.AllowMethods) == 0 {
			config.AllowMethods = []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}
		}
		config.AllowCredentials = false // Default
	}

	// Basic validation
	if len(config.AllowOrigins) == 0 || len(config.AllowMethods) == 0 {
		return fmt.Errorf("CORS config must specify at least AllowOrigins and AllowMethods (set IsPermissive=true for defaults)")
	}

	return nil
}

// GetIsEnabled returns true if the CORS config is meaningfully set.
func (gc *CORSConfig) GetIsEnabled() bool {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return len(gc.AllowOrigins) > 0 || len(gc.PerOriginConfigs) > 0
}

// GetCustomMiddleware returns a dynamic CORS middleware func for per-origin logic.
func (gc *CORSConfig) GetCustomMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			origin := req.Header.Get(echo.HeaderOrigin)

			res.Header().Add(echo.HeaderVary, echo.HeaderOrigin)

			// Preflight check
			preflight := req.Method == http.MethodOptions

			// Skip if no origin or skipper allows
			defaultConfig := gc.CORSConfig
			if defaultConfig.Skipper != nil && defaultConfig.Skipper(c) {
				return next(c)
			}
			if origin == "" {
				if !preflight {
					return next(c)
				}
				return c.NoContent(http.StatusNoContent)
			}

			// Select config based on origin
			selectedConfig := defaultConfig
			if perConfig, ok := gc.PerOriginConfigs[origin]; ok {
				selectedConfig = perConfig
			}

			// Prepare patterns for origin matching
			allowOriginPatterns := make([]*regexp.Regexp, 0, len(selectedConfig.AllowOrigins))
			for _, o := range selectedConfig.AllowOrigins {
				if o == "*" {
					continue
				}
				pattern := regexp.QuoteMeta(o)
				pattern = strings.ReplaceAll(pattern, "\\*", ".*")
				pattern = strings.ReplaceAll(pattern, "\\?", ".")
				pattern = "^" + pattern + "$"
				re, err := regexp.Compile(pattern)
				if err == nil {
					allowOriginPatterns = append(allowOriginPatterns, re)
				}
			}

			allowMethods := strings.Join(selectedConfig.AllowMethods, ",")
			allowHeaders := strings.Join(selectedConfig.AllowHeaders, ",")
			exposeHeaders := strings.Join(selectedConfig.ExposeHeaders, ",")
			maxAge := "0"
			if selectedConfig.MaxAge > 0 {
				maxAge = strconv.Itoa(selectedConfig.MaxAge)
			}

			allowOrigin := ""
			if selectedConfig.AllowOriginFunc != nil {
				allowed, err := selectedConfig.AllowOriginFunc(origin)
				if err != nil {
					return err
				}
				if allowed {
					allowOrigin = origin
				}
			} else {
				for _, o := range selectedConfig.AllowOrigins {
					if o == "*" && selectedConfig.AllowCredentials && selectedConfig.UnsafeWildcardOriginWithAllowCredentials {
						allowOrigin = origin
						break
					}
					if o == "*" || o == origin {
						allowOrigin = o
						break
					}
					if matchSubdomain(origin, o) {
						allowOrigin = origin
						break
					}
				}
				if allowOrigin == "" {
					for _, re := range allowOriginPatterns {
						if re.MatchString(origin) {
							allowOrigin = origin
							break
						}
					}
				}
			}

			// Origin not allowed
			if allowOrigin == "" {
				if !preflight {
					return next(c)
				}
				return c.NoContent(http.StatusNoContent)
			}

			res.Header().Set(echo.HeaderAccessControlAllowOrigin, allowOrigin)
			if selectedConfig.AllowCredentials {
				res.Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
			}

			// Simple request
			if !preflight {
				if exposeHeaders != "" {
					res.Header().Set(echo.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				return next(c)
			}

			// Preflight request
			res.Header().Add(echo.HeaderVary, echo.HeaderAccessControlRequestMethod)
			res.Header().Add(echo.HeaderVary, echo.HeaderAccessControlRequestHeaders)

			// Handle router-specific Allow for OPTIONS
			routerAllowMethods := ""
			tmpAllowMethods, ok := c.Get(echo.ContextKeyHeaderAllow).(string)
			if ok && tmpAllowMethods != "" {
				routerAllowMethods = tmpAllowMethods
				res.Header().Set(echo.HeaderAllow, routerAllowMethods)
			}

			hasCustomAllowMethods := len(selectedConfig.AllowMethods) > 0
			if !hasCustomAllowMethods && routerAllowMethods != "" {
				res.Header().Set(echo.HeaderAccessControlAllowMethods, routerAllowMethods)
			} else {
				res.Header().Set(echo.HeaderAccessControlAllowMethods, allowMethods)
			}

			if allowHeaders != "" {
				res.Header().Set(echo.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := req.Header.Get(echo.HeaderAccessControlRequestHeaders)
				if h != "" {
					res.Header().Set(echo.HeaderAccessControlAllowHeaders, h)
				}
			}
			if selectedConfig.MaxAge > 0 {
				res.Header().Set(echo.HeaderAccessControlMaxAge, maxAge)
			}

			return c.NoContent(http.StatusNoContent)
		}
	}
}

// matchSubdomain checks if origin matches a subdomain pattern (e.g., *.example.com).
func matchSubdomain(origin, pattern string) bool {
	if !strings.HasPrefix(pattern, "*.") {
		return false
	}
	pattern = pattern[2:] // Remove *.
	return strings.HasSuffix(origin, pattern)
}
