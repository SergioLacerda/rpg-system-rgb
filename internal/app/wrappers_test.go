package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFacadeMethodsSurfaceComponentErrors(t *testing.T) {
	root := t.TempDir()
	if err := BuildBundle(root); err == nil {
		t.Fatal("expected BuildBundle to fail without semantic index")
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
	if err := CheckPDFRegression(PDFRegressionPaths{BaselineDir: filepath.Join(root, "downloads"), CandidateDir: filepath.Join(root, "downloads"), Basename: "rgb", Version: "v1"}); err == nil {
		t.Fatal("expected CheckPDFRegression to fail without release PDFs or PDF tools")
	}
	if err := PackageSkill(SkillOptions{SourceDir: filepath.Join(root, "missing-skill"), OutDir: filepath.Join(root, "downloads"), Name: "rgb-specialist", Version: "v1"}); err == nil {
		t.Fatal("expected PackageSkill to fail without a skill source directory")
	}
	if err := WriteSkillManifest(SkillManifestPaths{PublicDir: filepath.Join(root, "downloads"), Name: "rgb-specialist", Version: "v1"}); err == nil {
		t.Fatal("expected WriteSkillManifest to fail without skill artifacts")
	}
	if err := CheckSkillManifest(SkillManifestPaths{PublicDir: filepath.Join(root, "downloads"), Name: "rgb-specialist", Version: "v1"}); err == nil {
		t.Fatal("expected CheckSkillManifest to fail without skill release artifacts")
	}
	if err := ValidateDocs(root); err == nil {
		t.Fatal("expected ValidateDocs to fail without semantic docs")
	}
	if err := GenerateProjections(root); err == nil {
		t.Fatal("expected GenerateProjections to fail without semantic docs")
	}
	if _, err := ProcessMakerBatch(filepath.Join(root, "missing.json"), "", ""); err == nil {
		t.Fatal("expected ProcessMakerBatch to fail when batch file is missing")
	}
}

func TestProcessMakerBatchSuccessAndIndexErrors(t *testing.T) {
	dir := t.TempDir()
	batchFile := filepath.Join(dir, "batch.json")
	batch := `{
		"report_id": "report.app-test",
		"input_batch": ["note.txt"],
		"entity_packages": [{
			"id": "npc.guard",
			"entity_type": "npc",
			"status": "draft",
			"source_refs": ["note.txt"],
			"explicit_facts": [{"field": "role", "value": "watchman", "source_ref": "note.txt"}],
			"canonical_dependencies": ["core.ability.contract"]
		}]
	}`
	if err := os.WriteFile(batchFile, []byte(batch), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := ProcessMakerBatch(batchFile, filepath.Join(dir, "missing-index.json"), ""); err == nil {
		t.Fatal("expected missing index file to fail")
	}

	indexFile := filepath.Join(dir, "index.json")
	if err := os.WriteFile(indexFile, []byte(`{"units":[{"id":"core.ability.contract","authority_type":"canonical_semantic"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(dir, "out")
	md, err := ProcessMakerBatch(batchFile, indexFile, outDir)
	if err != nil {
		t.Fatalf("ProcessMakerBatch failed: %v", err)
	}
	if !strings.Contains(md, "# Entity Package: npc.guard") {
		t.Fatalf("unexpected markdown:\n%s", md)
	}
	if _, err := os.Stat(filepath.Join(outDir, "entities", "npc.guard.md")); err != nil {
		t.Fatalf("expected rendered entity file: %v", err)
	}

	if _, err := ProcessMakerBatch(batchFile, indexFile, batchFile); err == nil {
		t.Fatal("expected unwritable outDir to fail")
	}

	if err := os.WriteFile(batchFile, []byte(`{"report_id":""}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ProcessMakerBatch(batchFile, "", ""); err == nil {
		t.Fatal("expected invalid batch to fail")
	}

	statsMD, err := MakerWorkspaceStats(dir)
	if err != nil || !strings.Contains(statsMD, "# Maker Workspace Stats:") {
		t.Fatalf("MakerWorkspaceStats failed: err=%v out=%s", err, statsMD)
	}
	if _, err := MakerWorkspaceStats(filepath.Join(dir, "missing")); err == nil {
		t.Fatal("expected MakerWorkspaceStats on missing dir to fail")
	}

	if _, err := ExportMakerChatPack(filepath.Join(dir, "missing"), filepath.Join(dir, "out")); err == nil {
		t.Fatal("expected ExportMakerChatPack on missing workspace to fail")
	}
	emptyWS := t.TempDir()
	if err := os.MkdirAll(filepath.Join(emptyWS, "packages"), 0o750); err != nil {
		t.Fatal(err)
	}
	if _, err := ExportMakerChatPack(emptyWS, filepath.Join(emptyWS, "out")); err == nil {
		t.Fatal("expected ExportMakerChatPack with no batch json to fail")
	}
	badJSONWS := t.TempDir()
	if err := os.MkdirAll(filepath.Join(badJSONWS, "packages"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(badJSONWS, "packages", "bad.json"), []byte(`{`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ExportMakerChatPack(badJSONWS, filepath.Join(badJSONWS, "out")); err == nil {
		t.Fatal("expected ExportMakerChatPack with invalid json to fail")
	}
	if err := os.WriteFile(filepath.Join(badJSONWS, "packages", "bad.json"), []byte(`{"report_id":"r","entity_packages":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ExportMakerChatPack(badJSONWS, filepath.Join(badJSONWS, "out")); err == nil {
		t.Fatal("expected ExportMakerChatPack with empty entities to fail")
	}
	goodWS := t.TempDir()
	if err := os.MkdirAll(filepath.Join(goodWS, "packages"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goodWS, "packages", "good.json"), []byte(`{
		"report_id": "report.export",
		"input_batch": ["n.txt"],
		"entity_packages": [{"id":"npc.a","entity_type":"npc","status":"draft","source_refs":["n.txt"],"explicit_facts":[{"field":"role","value":"scout","source_ref":"n.txt"}]}]
	}`), 0o600); err != nil {
		t.Fatal(err)
	}
	summary, err := ExportMakerChatPack(goodWS, filepath.Join(goodWS, "chat"))
	if err != nil || !strings.Contains(summary, "00_BOOT_PROMPT.md") {
		t.Fatalf("ExportMakerChatPack failed: err=%v out=%s", err, summary)
	}
}
