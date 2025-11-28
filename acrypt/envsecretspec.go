package acrypt

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// EnvSecretSpec defines a secret to load from env vars.
type EnvSecretSpec struct {
	Key     SecretsKey // The target key (e.g., "conn:ldap:bind:pass").
	ValVar  string     // Env var for inline value (e.g., "LDAP_BIND_PASS").
	FileVar string     // Env var for explicit file path (e.g., "LDAP_BIND_PASS_FILE").
}

// GetInitialSecretsFromEnv loads secrets based on provided specs.
// Customize specs for your app's needs.
func GetInitialSecretsFromEnv(specs []EnvSecretSpec) map[SecretsKey]string {
	secrets := make(map[SecretsKey]string)
	for _, s := range specs {
		if v, ok := ResolveFromEnv(s.ValVar, s.FileVar); ok && v != "" {
			secrets[s.Key] = v
		}
	}
	return secrets
}

// ResolveFromEnv resolves a value from env vars with file support.
// Precedence: explicit *_FILE > inline var (which may itself reference a file with @/file:/~/).
// Returns (value, ok); ok is false if no value found.
// Uses fmt.Printf for warnings since logging may not be initialized yet.
func ResolveFromEnv(valVar, fileVar string) (string, bool) {
	// 1) *_FILE has highest priority.
	if p := strings.TrimSpace(os.Getenv(fileVar)); p != "" {
		if v, err := readSecretFile(p); err == nil {
			return v, true
		} else {
			fmt.Printf("WARNING: failed to read secret file %q for %s: %v\n", p, valVar, err)
			// Fall through to valVar (fail-soft).
		}
	}

	// 2) Inline value or inline file reference.
	raw := strings.TrimSpace(os.Getenv(valVar))
	if raw == "" {
		return "", false
	}

	// If inline looks like a file reference, try reading it.
	if maybePath, isRef := asFileReference(raw); isRef {
		if v, err := readSecretFile(maybePath); err == nil {
			return v, true
		} else if !errors.Is(err, os.ErrNotExist) {
			// Warn on non-ENOENT errors.
			fmt.Printf("WARNING: failed to read inline secret file %q for %s: %v; using literal value\n", maybePath, valVar, err)
		} // ENOENT: treat as literal.
	}

	// Treat as a literal password.
	return raw, true
}

// asFileReference checks if s encodes a path reference rather than a literal secret.
// Supported forms:
//   - "@/path/to/secret"
//   - "file:/path/to/secret", "file:///abs/path"
//   - "~/secret" or "~"
//   - Absolute/relative that *looks* like a path (heuristic: contains path separator).
//
// Returns (cleanedPath, isReference).
func asFileReference(s string) (string, bool) {
	// @/path
	if strings.HasPrefix(s, "@") {
		return cleanPath(s[1:]), true
	}

	// file: URL
	if strings.HasPrefix(s, "file:") {
		u, err := url.Parse(s)
		if err == nil && u.Scheme == "file" {
			// url.Parse("file:/x") yields Path "/x"; "file:///x" also "/x".
			return cleanPath(u.Path), true
		}
		// Malformed: treat as literal (no warn here; fallback).
		return "", false
	}

	// ~ expansion
	if strings.HasPrefix(s, "~/") || s == "~" {
		return expandUser(cleanPath(s)), true
	}

	// Heuristic: if it *looks* like a path (contains a path separator) we’ll try it.
	// If the read fails with ENOENT, we’ll use the literal.
	if strings.Contains(s, "/") || strings.Contains(s, string(os.PathSeparator)) {
		return cleanPath(s), true
	}

	// Otherwise, treat as literal.
	return "", false
}

// readSecretFile reads a small text secret file and trims trailing whitespace/newlines.
// Returns error if file doesn't exist, too large, or read fails.
func readSecretFile(p string) (string, error) {
	path := expandUser(cleanPath(p))
	// Avoid a separate os.Stat; just open & read. If it doesn't exist, we’ll get ENOENT.
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Limit size to protect against accidental huge reads (e.g., 1 MiB).
	const max = 1 << 20
	var sb strings.Builder
	sb.Grow(256) // Pre-allocate for typical secrets.
	lr := io.LimitedReader{R: f, N: max + 1}
	if _, err := io.Copy(&sb, &lr); err != nil {
		return "", err
	}
	if lr.N <= 0 {
		return "", fmt.Errorf("secret file %q exceeds %d bytes", path, max)
	}

	secret := strings.TrimSpace(sb.String())
	return secret, nil
}

// cleanPath trims and cleans a path (removes . and ..).
func cleanPath(p string) string {
	return filepath.Clean(strings.TrimSpace(p))
}

// expandUser expands ~ to user home dir if present.
func expandUser(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p // Fallback: no expansion.
		}
		if p == "~" {
			return home
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
