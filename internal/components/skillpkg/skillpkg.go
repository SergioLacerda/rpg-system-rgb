// Package skillpkg zips a skill directory into a downloadable .zip for the
// existing Go-owned publication pipeline (ADR-013). It is a distinct
// component from internal/components/tooling's PDF release manifest/check:
// a skill archive has no PDF-style editorial validation to perform, so this
// package never imports tooling (see manifest.go).
package skillpkg

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/SergioLacerda/rpg-system-rgb/internal/components"
)

// Descriptor identifies this component to the application layer.
func Descriptor() components.Component {
	return components.Component{
		ID:          "skillpkg",
		Name:        "RGB Skill Packaging",
		Description: "Zips a skill directory and writes its release manifest for the existing publication pipeline (ADR-013).",
	}
}

// Options configures packaging a skill directory into a downloadable .zip.
type Options struct {
	SourceDir string
	OutDir    string
	Name      string
	Version   string
}

// Package zips options.SourceDir into <name>-<version>.zip and
// <name>-latest.zip under options.OutDir. The archive's top-level entry is
// the skill's own directory name (options.Name), never a repo-root-relative
// or GitHub-auto-generated wrapper folder.
func Package(options Options) error {
	if err := validatePackageOptions(options); err != nil {
		return err
	}
	if err := os.MkdirAll(options.OutDir, 0o750); err != nil {
		return err
	}

	versioned := filepath.Join(options.OutDir, fmt.Sprintf("%s-%s.zip", options.Name, options.Version))
	if err := writeZip(options.SourceDir, options.Name, versioned); err != nil {
		return err
	}
	latest := filepath.Join(options.OutDir, fmt.Sprintf("%s-latest.zip", options.Name))
	return copyFile(versioned, latest)
}

func validatePackageOptions(options Options) error {
	if options.SourceDir == "" || options.OutDir == "" || options.Name == "" || options.Version == "" {
		return fmt.Errorf("source, out, name, and version are required")
	}
	info, err := os.Stat(options.SourceDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("missing skill source directory: %s", options.SourceDir)
	}
	return nil
}

func writeZip(sourceDir, entryRoot, dest string) error {
	file, err := os.Create(dest) //nolint:gosec // G304: destination path is provided by the Make target.
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := zip.NewWriter(file)
	walkErr := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(sourceDir, path)
		if relErr != nil {
			return relErr
		}
		return addZipFile(writer, path, filepath.ToSlash(filepath.Join(entryRoot, rel)))
	})
	if walkErr != nil {
		_ = writer.Close()
		return walkErr
	}
	return writer.Close()
}

func addZipFile(writer *zip.Writer, sourcePath, zipPath string) error {
	content, err := os.ReadFile(sourcePath) //nolint:gosec // G304: source path is a file discovered under a Make-provided skill directory.
	if err != nil {
		return err
	}
	content = sanitizePackagedSkillContent(zipPath, content)
	entry, err := writer.Create(zipPath)
	if err != nil {
		return err
	}
	_, err = entry.Write(content)
	return err
}

var packagedSkillEngineeringHeading = regexp.MustCompile(`(?i)^(#{1,6}\s*)?(architecture decisions|decis(?:õ|o)es de arquitetura|relationship to prior adrs)\s*$`)

func sanitizePackagedSkillContent(zipPath string, content []byte) []byte {
	switch strings.ToLower(filepath.Ext(zipPath)) {
	case ".md", ".yaml", ".yml", ".txt":
	default:
		return content
	}
	var out []string
	skipSection := false
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "#") {
			skipSection = packagedSkillEngineeringHeading.MatchString(line)
		}
		if skipSection || containsADRReference(line) {
			continue
		}
		out = append(out, raw)
	}
	return []byte(strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n")
}

func containsADRReference(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "docs/adr/") ||
		strings.Contains(lower, "../adr/") ||
		strings.Contains(lower, "../../adr/") ||
		strings.Contains(lower, "/adr/") ||
		strings.Contains(lower, "adr-") ||
		strings.Contains(lower, "adrs")
}

func copyFile(source, dest string) error {
	content, err := os.ReadFile(source) //nolint:gosec // G304: source path is produced by this package.
	if err != nil {
		return err
	}
	return os.WriteFile(dest, content, 0o644) //nolint:gosec // G306: public release artifact.
}
