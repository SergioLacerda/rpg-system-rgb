package skillpkg

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDescriptor(t *testing.T) {
	descriptor := Descriptor()
	if descriptor.ID == "" || descriptor.Name == "" || descriptor.Description == "" {
		t.Fatalf("descriptor has empty fields: %#v", descriptor)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("creating %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func TestPackageWritesVersionedAndLatestZips(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "specialist")
	writeTestFile(t, filepath.Join(source, "SKILL.md"), "# Specialist")
	writeTestFile(t, filepath.Join(source, "procedures", "explain.md"), "explain procedure")

	out := filepath.Join(root, "downloads")
	options := Options{SourceDir: source, OutDir: out, Name: "rgb-specialist", Version: "v0.1"}
	if err := Package(options); err != nil {
		t.Fatalf("Package returned error: %v", err)
	}

	versioned := filepath.Join(out, "rgb-specialist-v0.1.zip")
	latest := filepath.Join(out, "rgb-specialist-latest.zip")
	for _, path := range []string{versioned, latest} {
		if info, err := os.Stat(path); err != nil || info.Size() == 0 {
			t.Fatalf("expected non-empty zip at %s: %v", path, err)
		}
	}

	entries := zipEntryNames(t, versioned)
	want := map[string]bool{
		"rgb-specialist/SKILL.md":              true,
		"rgb-specialist/procedures/explain.md": true,
	}
	if len(entries) != len(want) {
		t.Fatalf("unexpected zip entries: %v", entries)
	}
	for _, name := range entries {
		if !want[name] {
			t.Fatalf("unexpected zip entry %q (top-level entry must be the skill directory name, not a wrapper folder)", name)
		}
	}
}

func TestPackageRemovesADRReferencesFromGeneratedZip(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "specialist")
	writeTestFile(t, filepath.Join(source, "SKILL.md"), strings.Join([]string{
		"# Specialist",
		"",
		"Public skill guidance.",
		"",
		"See [ADR-002](../../docs/adr/adr-002-rgb-core-v2-design-control.md).",
		"",
		"## Relationship To Prior ADRs",
		"",
		"Internal publication history.",
		"",
		"## Usage",
		"",
		"Use the rules.",
	}, "\n"))
	writeTestFile(t, filepath.Join(source, "config.yaml"), strings.Join([]string{
		"# Bundled fallback, see ADR-014.",
		"reference: references/rgb-system.md",
	}, "\n"))

	out := filepath.Join(root, "downloads")
	options := Options{SourceDir: source, OutDir: out, Name: "rgb-specialist", Version: "v0.1"}
	if err := Package(options); err != nil {
		t.Fatalf("Package returned error: %v", err)
	}

	for _, name := range []string{"rgb-specialist/SKILL.md", "rgb-specialist/config.yaml"} {
		got := zipEntryContent(t, filepath.Join(out, "rgb-specialist-v0.1.zip"), name)
		if strings.Contains(got, "ADR-") || strings.Contains(got, "docs/adr") || strings.Contains(strings.ToLower(got), "adrs") {
			t.Fatalf("expected %s to omit ADR content, got:\n%s", name, got)
		}
	}
}

func zipEntryNames(t *testing.T, path string) []string {
	t.Helper()
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("opening zip %s: %v", path, err)
	}
	defer func() {
		_ = reader.Close()
	}()
	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, file.Name)
	}
	return names
}

func zipEntryContent(t *testing.T, path, name string) string {
	t.Helper()
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("opening zip %s: %v", path, err)
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		if file.Name != name {
			continue
		}
		handle, err := file.Open()
		if err != nil {
			t.Fatalf("opening zip entry %s: %v", name, err)
		}
		defer func() {
			_ = handle.Close()
		}()
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, handle); err != nil {
			t.Fatalf("reading zip entry %s: %v", name, err)
		}
		return buf.String()
	}
	t.Fatalf("zip entry %s not found in %s", name, path)
	return ""
}

func TestPackageRequiresAllOptions(t *testing.T) {
	root := t.TempDir()
	cases := []Options{
		{OutDir: root, Name: "n", Version: "v1"},
		{SourceDir: root, Name: "n", Version: "v1"},
		{SourceDir: root, OutDir: root, Version: "v1"},
		{SourceDir: root, OutDir: root, Name: "n"},
	}
	for _, options := range cases {
		if err := Package(options); err == nil {
			t.Fatalf("expected error for incomplete options: %#v", options)
		}
	}
}

func TestWriteZipFailsWhenDestinationIsUnwritable(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "specialist")
	writeTestFile(t, filepath.Join(source, "SKILL.md"), "# Specialist")

	if err := writeZip(source, "rgb-specialist", filepath.Join(root, "missing-dir", "out.zip")); err == nil {
		t.Fatal("expected writeZip to fail when the destination directory does not exist")
	}
}

func TestAddZipFileFailsForMissingSource(t *testing.T) {
	root := t.TempDir()
	file, err := os.Create(filepath.Join(root, "out.zip"))
	if err != nil {
		t.Fatalf("creating zip file: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()
	writer := zip.NewWriter(file)
	defer func() {
		_ = writer.Close()
	}()

	if err := addZipFile(writer, filepath.Join(root, "missing.md"), "rgb-specialist/missing.md"); err == nil {
		t.Fatal("expected addZipFile to fail for a missing source file")
	}
}

func TestCopyFileFailsForMissingSource(t *testing.T) {
	root := t.TempDir()
	if err := copyFile(filepath.Join(root, "missing.zip"), filepath.Join(root, "latest.zip")); err == nil {
		t.Fatal("expected copyFile to fail for a missing source file")
	}
}

func TestPackageRejectsMissingSourceDirectory(t *testing.T) {
	root := t.TempDir()
	options := Options{
		SourceDir: filepath.Join(root, "missing"),
		OutDir:    filepath.Join(root, "out"),
		Name:      "rgb-specialist",
		Version:   "v0.1",
	}
	if err := Package(options); err == nil {
		t.Fatal("expected error for missing source directory")
	}
}
