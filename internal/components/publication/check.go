package publication

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckOptions configures publication artifact validation.
type CheckOptions struct {
	LibraryDir string
	PublicDir  string
	Basename   string
	Version    string
}

// Check validates the Go-owned publication output at a smoke-test level.
func Check(options CheckOptions) error {
	if options.LibraryDir == "" || options.PublicDir == "" || options.Basename == "" || options.Version == "" {
		return fmt.Errorf("library, public-dir, basename, and version are required")
	}
	for _, path := range []string{
		filepath.Join(options.LibraryDir, "index.html"),
		filepath.Join(options.LibraryDir, "core", "en", "index.html"),
		filepath.Join(options.LibraryDir, "core", "PT-br", "index.html"),
		filepath.Join(options.PublicDir, fmt.Sprintf("%s-latest-en.pdf", options.Basename)),
		filepath.Join(options.PublicDir, fmt.Sprintf("%s-%s-en.pdf", options.Basename, options.Version)),
		filepath.Join(options.PublicDir, fmt.Sprintf("%s-latest-pt-br.pdf", options.Basename)),
		filepath.Join(options.PublicDir, fmt.Sprintf("%s-%s-pt-br.pdf", options.Basename, options.Version)),
	} {
		if err := requireFile(path); err != nil {
			return err
		}
	}
	return nil
}

func requireFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("missing publication artifact: %s", path)
	}
	if info.IsDir() || info.Size() == 0 {
		return fmt.Errorf("invalid publication artifact: %s", path)
	}
	return nil
}
