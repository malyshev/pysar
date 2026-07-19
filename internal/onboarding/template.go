package onboarding

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// templateWrapper is template-only metadata that never becomes part of
// Profile's own schema (dec-20260719-b12539fa invariant): a display name a
// template is saved under, addressed separately by its own stable slug so a
// rename never changes what the template is called on disk.
type templateWrapper struct {
	Name string `yaml:"name"`
}

// WrapTemplate prepends a small, separate frontmatter block carrying name
// ahead of profileContent, which must already be exactly what Render
// produces for a Profile. UnwrapTemplate reverses this exactly.
func WrapTemplate(name, profileContent string) (string, error) {
	fm, err := yaml.Marshal(templateWrapper{Name: name})
	if err != nil {
		return "", fmt.Errorf("render template wrapper: %w", err)
	}
	return "---\n" + string(fm) + "---\n\n" + profileContent, nil
}

// UnwrapTemplate splits a saved template file back into its display name and
// the underlying profile content. It requires the leading wrapper block;
// content saved without one returns an error naming what's missing, never a
// guessed name.
func UnwrapTemplate(content string) (name string, profileContent string, err error) {
	const openDelim = "---\n"
	const closeDelim = "\n---\n"

	if !strings.HasPrefix(content, openDelim) {
		return "", "", fmt.Errorf("template missing leading frontmatter wrapper")
	}
	rest := content[len(openDelim):]

	end := strings.Index(rest, closeDelim)
	if end == -1 {
		return "", "", fmt.Errorf("template wrapper frontmatter is not terminated")
	}

	var w templateWrapper
	if err := yaml.Unmarshal([]byte(rest[:end]), &w); err != nil {
		return "", "", fmt.Errorf("parse template wrapper: %w", err)
	}
	if strings.TrimSpace(w.Name) == "" {
		// A plain Render()-produced profile also starts with "---\n" --
		// its own frontmatter parses here as a wrapper with no "name" key,
		// which is exactly how a missing wrapper is distinguished from a
		// real one (save always requires a non-empty name).
		return "", "", fmt.Errorf("template missing leading frontmatter wrapper")
	}

	profileContent = strings.TrimPrefix(rest[end+len(closeDelim):], "\n")
	return w.Name, profileContent, nil
}
