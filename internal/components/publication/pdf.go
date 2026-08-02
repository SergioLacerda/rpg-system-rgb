package publication

import (
	"fmt"
	"os"
	"path/filepath"
)

// PDFOptions configures PDF publication without the retired renderer runtime.
type PDFOptions struct {
	PublicDir string
	Basename  string
	Version   string
	SourceEN  string
	SourcePT  string
}

// PublishPDFs publishes curated PDFs to the stable latest and versioned names.
func PublishPDFs(options PDFOptions) error {
	if options.PublicDir == "" || options.Basename == "" || options.Version == "" {
		return fmt.Errorf("public-dir, basename, and version are required")
	}
	if err := os.MkdirAll(options.PublicDir, 0o750); err != nil {
		return err
	}
	if err := publishLocalePDF(options, "en", options.SourceEN); err != nil {
		return err
	}
	return publishLocalePDF(options, "pt-br", options.SourcePT)
}

func publishLocalePDF(options PDFOptions, locale, source string) error {
	latest := filepath.Join(options.PublicDir, fmt.Sprintf("%s-latest-%s.pdf", options.Basename, locale))
	versioned := filepath.Join(options.PublicDir, fmt.Sprintf("%s-%s-%s.pdf", options.Basename, options.Version, locale))
	if source == "" {
		if fileExists(versioned) {
			source = versioned
		} else {
			source = latest
		}
	}
	if !fileExists(source) {
		return fmt.Errorf("missing %s PDF source: %s", locale, source)
	}
	if err := copyFile(source, latest); err != nil {
		return err
	}
	return copyFile(latest, versioned)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func copyFile(source, dest string) error {
	content, err := os.ReadFile(source) //nolint:gosec // G304: publication source path is provided by the Make target.
	if err != nil {
		return err
	}
	if len(content) < 5 || string(content[:5]) != "%PDF-" {
		return fmt.Errorf("%s is not a PDF", source)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0o644) //nolint:gosec // G306: public release artifact.
}
