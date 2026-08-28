package publication

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildLibraryRendersBilingualCoreDocs(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "docs")
	out := filepath.Join(root, "public", "library")
	writeTestFile(t, filepath.Join(source, "core", "en", "README.md"), "# Core Overview\n\nSee [Combat](combat/attack.md).\n")
	writeTestFile(t, filepath.Join(source, "core", "en", "combat", "attack.md"), "# Attack\n\n- Roll R\n")
	writeTestFile(t, filepath.Join(source, "core", "PT-br", "README.md"), "# Visao Core\n\nVeja [Combate](combat/ataque.md).\n")
	writeTestFile(t, filepath.Join(source, "core", "PT-br", "combat", "ataque.md"), "# Ataque\n\n- Role R\n")

	if err := BuildLibrary(LibraryOptions{SourceDir: source, OutDir: out}); err != nil {
		t.Fatalf("BuildLibrary returned error: %v", err)
	}

	index := readTestFile(t, filepath.Join(out, "index.html"))
	if !strings.Contains(index, "<h1 id=\"visao-core\">Visao Core</h1>") {
		t.Fatalf("expected index to render the pt-BR root page by default, got:\n%s", index)
	}
	if !strings.Contains(index, `<a class="home" href="../">Home</a>`) {
		t.Fatalf("expected index to link back to the landing home, got:\n%s", index)
	}
	if !strings.Contains(index, `href="core/en/"`) {
		t.Fatalf("expected index to keep the English tree reachable, got:\n%s", index)
	}
	if !strings.Contains(index, `href="core/PT-br/combat/ataque/"`) {
		t.Fatalf("expected index's aliased body links to be re-anchored under core/PT-br/, got:\n%s", index)
	}

	ptPage := readTestFile(t, filepath.Join(out, "core", "PT-br", "index.html"))
	if !strings.Contains(ptPage, "<h1 id=\"visao-core\">Visao Core</h1>") {
		t.Fatalf("expected the core/PT-br tree to remain unchanged, got:\n%s", ptPage)
	}

	page := readTestFile(t, filepath.Join(out, "core", "en", "index.html"))
	if !strings.Contains(page, "<h1 id=\"core-overview\">Core Overview</h1>") {
		t.Fatalf("expected rendered heading, got:\n%s", page)
	}
	if !strings.Contains(page, "href=\"combat/attack/\"") {
		t.Fatalf("expected rewritten markdown link, got:\n%s", page)
	}
	if !strings.Contains(page, `<a class="home" href="../../../">Home</a>`) {
		t.Fatalf("expected content page to link back to the landing home, got:\n%s", page)
	}
}

// The original 14-unit Core-V2-pilot subset the search-index projection
// covered before .analysis/refined/20260827-html-library-search-refinement/
// widened it to the full master semantic index (34 units) — a regression
// test that this widening never silently drops one of them.
var originalFourteenSearchUnitIDs = []string{
	"core.vector.r", "core.vector.g", "core.vector.b",
	"core.resource.health", "core.resource.shield",
	"core.combat.attack-margin", "core.combat.evade",
	"core.damage.flow", "core.damage.impact-source", "core.damage.penetration",
	"core.damage.armor-reduction", "core.damage.shield-absorption",
	"core.ability.contract", "core.example.combat-turn",
}

func TestSearchIndexWideningPreservesOriginalFourteen(t *testing.T) {
	root := repoRoot(t)
	raw := readTestFile(t, filepath.Join(root, "generated", "search", "core-v2.index.json"))
	var projection struct {
		Units []struct {
			ID string `json:"id"`
		} `json:"units"`
	}
	if err := json.Unmarshal([]byte(raw), &projection); err != nil {
		t.Fatalf("invalid search-index projection JSON: %v", err)
	}
	if len(projection.Units) < 34 {
		t.Fatalf("expected the widened search-index projection to cover at least 34 units, got %d", len(projection.Units))
	}
	present := map[string]bool{}
	for _, unit := range projection.Units {
		present[unit.ID] = true
	}
	for _, id := range originalFourteenSearchUnitIDs {
		if !present[id] {
			t.Fatalf("widening dropped originally-covered unit %s", id)
		}
	}
}

func TestSearchIndexUnitsResolveToRealLibraryPages(t *testing.T) {
	root := repoRoot(t)
	out := t.TempDir()

	err := BuildLibrary(LibraryOptions{
		SourceDir:       filepath.Join(root, "docs"),
		OutDir:          out,
		SearchIndexFile: filepath.Join(root, "generated", "search", "core-v2.index.json"),
	})
	if err != nil {
		t.Fatalf("BuildLibrary returned error: %v", err)
	}

	indexPath := filepath.Join(out, "search-index.json")
	raw, err := os.ReadFile(indexPath) //nolint:gosec // G304: test-controlled temp path.
	if err != nil {
		t.Fatalf("expected search-index.json to be written: %v", err)
	}
	var entries []struct {
		ID      string `json:"id"`
		URLEn   string `json:"url_en"`
		URLPtBr string `json:"url_pt_br"`
	}
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatalf("invalid search-index.json: %v", err)
	}
	if len(entries) < 34 {
		t.Fatalf("expected at least 34 search-index entries, got %d", len(entries))
	}
	for _, entry := range entries {
		checked := 0
		for _, url := range []string{entry.URLEn, entry.URLPtBr} {
			if url == "" {
				continue
			}
			checked++
			pagePath := filepath.Join(out, filepath.FromSlash(strings.Trim(url, "/")), "index.html")
			if _, err := os.Stat(pagePath); err != nil {
				t.Errorf("%s: url %s does not resolve to a generated page (%s): %v", entry.ID, url, pagePath, err)
			}
		}
		if checked == 0 {
			t.Errorf("%s: no url_en or url_pt_br resolved from its projection_paths", entry.ID)
		}
	}
}

func TestSearchAssetsStayUnderSizeBudget(t *testing.T) {
	root := repoRoot(t)
	out := t.TempDir()

	err := BuildLibrary(LibraryOptions{
		SourceDir:       filepath.Join(root, "docs"),
		OutDir:          out,
		SearchIndexFile: filepath.Join(root, "generated", "search", "core-v2.index.json"),
	})
	if err != nil {
		t.Fatalf("BuildLibrary returned error: %v", err)
	}

	const budgetBytes = 50 * 1024 // item2's stated ceiling
	var total int64
	for _, name := range []string{"search-index.json", "search.js"} {
		info, err := os.Stat(filepath.Join(out, name))
		if err != nil {
			t.Fatalf("expected %s to exist: %v", name, err)
		}
		total += info.Size()
	}
	if total > budgetBytes {
		t.Fatalf("search assets total %d bytes, exceeds the %d byte budget", total, budgetBytes)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "docs", "core", "semantic", "core-v2.index.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repository root")
		}
		dir = parent
	}
}

func TestPublishPDFsPublishesLatestAndVersionedFiles(t *testing.T) {
	root := t.TempDir()
	public := filepath.Join(root, "downloads")
	en := filepath.Join(root, "en.pdf")
	pt := filepath.Join(root, "pt.pdf")
	writeTestFile(t, en, "%PDF-1.7\nEnglish\n")
	writeTestFile(t, pt, "%PDF-1.7\nPT-BR\n")

	err := PublishPDFs(PDFOptions{
		PublicDir: public,
		Basename:  "rgb-system-core-v2",
		Version:   "v0.2",
		SourceEN:  en,
		SourcePT:  pt,
	})
	if err != nil {
		t.Fatalf("PublishPDFs returned error: %v", err)
	}
	for _, name := range []string{
		"rgb-system-core-v2-latest-en.pdf",
		"rgb-system-core-v2-v0.2-en.pdf",
		"rgb-system-core-v2-latest-pt-br.pdf",
		"rgb-system-core-v2-v0.2-pt-br.pdf",
	} {
		if _, err := os.Stat(filepath.Join(public, name)); err != nil {
			t.Fatalf("expected %s: %v", name, err)
		}
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
