package app

import (
	"path/filepath"
	"testing"
)

func TestFacadeMethodsSurfaceComponentErrors(t *testing.T) {
	root := t.TempDir()
	if err := BuildBundle(root); err == nil {
		t.Fatal("expected BuildBundle to fail without semantic index")
	}
	if err := BuildLibrary(LibraryOptions{SourceDir: filepath.Join(root, "missing"), OutDir: filepath.Join(root, "out")}); err == nil {
		t.Fatal("expected BuildLibrary to fail without docs/core")
	}
	if err := PublishPDFs(PDFOptions{PublicDir: filepath.Join(root, "downloads"), Basename: "rgb", Version: "v1"}); err == nil {
		t.Fatal("expected PublishPDFs to fail without PDF sources")
	}
	if err := CheckPublication(PublicationCheckOptions{LibraryDir: filepath.Join(root, "library"), PublicDir: filepath.Join(root, "downloads"), Basename: "rgb", Version: "v1"}); err == nil {
		t.Fatal("expected CheckPublication to fail without artifacts")
	}
	if err := WriteReleaseArtifactManifest(ReleaseArtifactPaths{PublicDir: filepath.Join(root, "downloads"), Basename: "rgb", Version: "v1"}); err == nil {
		t.Fatal("expected WriteReleaseArtifactManifest to fail without PDF artifacts")
	}
	if err := CheckReleaseArtifacts(ReleaseArtifactPaths{PublicDir: filepath.Join(root, "downloads"), Basename: "rgb", Version: "v1"}); err == nil {
		t.Fatal("expected CheckReleaseArtifacts to fail without release artifacts")
	}
	if err := ValidateDocs(root); err == nil {
		t.Fatal("expected ValidateDocs to fail without semantic docs")
	}
	if err := GenerateProjections(root); err == nil {
		t.Fatal("expected GenerateProjections to fail without semantic docs")
	}
}
