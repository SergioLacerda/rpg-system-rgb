package maker

import (
	"fmt"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

// InputType classifies the kind of raw input the Maker receives.
type InputType string

const (
	// InputTypeText covers plain text, notes, PDFs (text-extracted), and
	// any input whose primary content is natural language.
	InputTypeText InputType = "text"
	// InputTypeImage covers bitmap images (.png, .jpg, .jpeg, .webp, .gif)
	// whose contribution goes through a VisualObservation, not directly into
	// ExplicitFacts.
	InputTypeImage InputType = "image"
)

// InputRef is the normalized, typed descriptor of one Maker input.
// Slug is a stable, lowercase identifier derived from the raw input name or
// content — it becomes the value used in source_refs across all emitted
// packages, so it must be consistent for a given source.
type InputRef struct {
	Raw  string    // original string as provided by the caller
	Slug string    // normalized slug for use in source_refs
	Type InputType // text | image
}

var slugReplacer = regexp.MustCompile(`[^a-z0-9]+`)

// slugify converts a raw string into a lowercase, hyphen-separated slug
// suitable for use in source_refs.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugReplacer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "input"
	}
	return s
}

// detectType classifies an input string as text or image based on its file
// extension (when it looks like a filename) using the standard MIME type
// registry. Inputs without a recognizable extension default to text.
func detectType(raw string) InputType {
	ext := strings.ToLower(filepath.Ext(raw))
	if ext == "" {
		return InputTypeText
	}
	mimeType := mime.TypeByExtension(ext)
	if strings.HasPrefix(mimeType, "image/") {
		return InputTypeImage
	}
	return InputTypeText
}

// RouteInput normalizes a single raw input string into an InputRef with a
// stable slug and detected type. It is deterministic and requires no LLM.
//
// Rules:
//   - Inputs with image extensions (.png, .jpg, .jpeg, .webp, .gif, etc.)
//     are classified as InputTypeImage.
//   - Everything else is classified as InputTypeText.
//   - The slug is derived from the raw string by lowercasing and replacing
//     non-alphanumeric runs with hyphens.
func RouteInput(raw string) (InputRef, error) {
	if strings.TrimSpace(raw) == "" {
		return InputRef{}, fmt.Errorf("router: raw input must not be empty")
	}

	// For filenames, use just the base name (without extension) for the slug
	// so that "harbor-npc-sketch.png" slugifies to "harbor-npc-sketch".
	base := filepath.Base(raw)
	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))
	slug := slugify(nameWithoutExt)

	return InputRef{
		Raw:  raw,
		Slug: slug,
		Type: detectType(raw),
	}, nil
}

// RouteInputs normalizes a slice of raw input strings, returning one
// InputRef per entry. Returns all errors encountered (one per invalid entry).
func RouteInputs(raws []string) ([]InputRef, []error) {
	refs := make([]InputRef, 0, len(raws))
	var errs []error
	for _, raw := range raws {
		ref, err := RouteInput(raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		refs = append(refs, ref)
	}
	return refs, errs
}
