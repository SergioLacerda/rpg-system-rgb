package publication

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// searchProjectionFile mirrors the subset of the generated search-index
// projection (generated/search/core-v2.index.json, produced by
// internal/components/tooling.Generate from docs/core/semantic/
// projection-manifest.v0.1.json) that the Library actually consumes, per
// consumer.library.v0_1 in docs/core/semantic/consumer-contracts.v0.1.json.
// Deliberately local and minimal (not imported from the tooling package) to
// keep internal/components/publication from importing a sibling component
// directly, per the project's component-boundary rule
// (tests/architecture/architecture_test.go: TestComponentsDoNotImportSiblings).
type searchProjectionFile struct {
	Units []searchProjectionUnit `json:"units"`
}

type searchProjectionUnit struct {
	ID               string              `json:"id"`
	Kind             string              `json:"kind"`
	Title            string              `json:"title"`
	RetrievalSummary string              `json:"retrieval_summary"`
	ProjectionPaths  map[string]string   `json:"projection_paths"`
	Relationships    map[string][]string `json:"relationships"`
	Tags             []string            `json:"tags"`
}

// searchEntry is one row of the Library's client-side search-index.json —
// one per semantic unit, carrying both locale URLs so search.js can pick
// the one matching the page the visitor is currently reading.
type searchEntry struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Kind    string   `json:"kind"`
	Vectors []string `json:"vectors"`
	Tags    []string `json:"tags"`
	URLEn   string   `json:"url_en,omitempty"`
	URLPtBr string   `json:"url_pt_br,omitempty"`
}

var coreVectorIDs = map[string]string{
	"core.vector.r": "R",
	"core.vector.g": "G",
	"core.vector.b": "B",
}

// buildSearchEntries reads the generated search-index projection and
// derives one searchEntry per unit: resolved Library page URLs (per
// locale) and an R/G/B vector facet.
//
// Per design.md § Vector Facet Derivation, no unit carries an explicit
// vector field — it is derived from relationships.depends_on intersected
// with {core.vector.r, core.vector.g, core.vector.b}, except a
// core.vector.{r,g,b} unit itself, which is its own vector. A unit may
// resolve to zero, one, or multiple vectors; callers must not assume
// exactly one.
func buildSearchEntries(searchIndexFile string) ([]searchEntry, error) {
	if searchIndexFile == "" {
		return nil, nil
	}
	raw, err := os.ReadFile(searchIndexFile) //nolint:gosec // G304: path is a caller-supplied CLI/Make target flag, by design.
	if err != nil {
		return nil, fmt.Errorf("read search index: %w", err)
	}
	var projection searchProjectionFile
	if err := json.Unmarshal(raw, &projection); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON: %w", searchIndexFile, err)
	}
	entries := make([]searchEntry, 0, len(projection.Units))
	for _, unit := range projection.Units {
		entries = append(entries, searchEntry{
			ID:      unit.ID,
			Title:   unit.Title,
			Summary: unit.RetrievalSummary,
			Kind:    unit.Kind,
			Vectors: deriveVectors(unit),
			Tags:    unit.Tags,
			URLEn:   projectionPathToLibraryURL(unit.ProjectionPaths["markdown_en"]),
			URLPtBr: projectionPathToLibraryURL(unit.ProjectionPaths["markdown_pt_br"]),
		})
	}
	return entries, nil
}

func deriveVectors(unit searchProjectionUnit) []string {
	if vector, ok := coreVectorIDs[unit.ID]; ok {
		return []string{vector}
	}
	var vectors []string
	for _, dep := range unit.Relationships["depends_on"] {
		if vector, ok := coreVectorIDs[dep]; ok {
			vectors = append(vectors, vector)
		}
	}
	return vectors
}

// projectionPathToLibraryURL converts a semantic-index projection path
// (e.g. "docs/core/en/core/attributes.md", rooted at the repo root) into
// the URL the Library serves that same page at (e.g.
// "/core/en/core/attributes/"), reusing the same convention collectPages
// uses for the pages it discovers directly from disk.
func projectionPathToLibraryURL(path string) string {
	if path == "" {
		return ""
	}
	rel := filepath.ToSlash(path)
	rel = strings.TrimPrefix(rel, "docs/")
	return outputURL(rel)
}

// writeSearchAssets writes search-index.json and search.js into the
// Library's output directory. Both are static assets — search.js is
// hand-written, dependency-free vanilla JavaScript (no bundled library),
// per ADR-005's zero-external-dependency posture and this mission's
// design.md.
func writeSearchAssets(outDir string, entries []searchEntry) error {
	if len(entries) == 0 {
		return nil
	}
	indexJSON, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "search-index.json"), indexJSON, 0o644); err != nil { //nolint:gosec // G306: public static asset.
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "search.js"), []byte(searchJS), 0o644) //nolint:gosec // G306: public static asset.
}

// searchJS is a hand-written, dependency-free client-side search engine.
// It loads search-index.json once per page load, matches the query against
// title/summary/tags, supports R/G/B vector filter chips (a unit may match
// zero, one, or multiple vectors), and provides Ctrl+K/"/" shortcuts with
// keyboard navigation and match highlighting.
const searchJS = `(function () {
  "use strict";
  var root = document.querySelector("[data-search-index]");
  if (!root) return;
  var lang = document.documentElement.lang === "pt-br" ? "pt-br" : "en";
  var input = root.querySelector("[data-search-input]");
  var results = root.querySelector("[data-search-results]");
  var chips = root.querySelectorAll("[data-vector]");
  var activeVectors = {};
  var entries = [];
  var activeIndex = -1;

  fetch(root.getAttribute("data-search-index"))
    .then(function (res) { return res.json(); })
    .then(function (data) { entries = data || []; })
    .catch(function () { entries = []; });

  function entryURL(entry) {
    return (lang === "pt-br" && entry.url_pt_br) ? entry.url_pt_br : (entry.url_en || entry.url_pt_br || "");
  }

  function matchesVectorFilter(entry) {
    var selected = Object.keys(activeVectors).filter(function (k) { return activeVectors[k]; });
    if (selected.length === 0) return true;
    var vectors = entry.vectors || [];
    return selected.some(function (v) { return vectors.indexOf(v) !== -1; });
  }

  function highlight(text, query) {
    if (!query) return escapeHTML(text);
    var idx = text.toLowerCase().indexOf(query.toLowerCase());
    if (idx === -1) return escapeHTML(text);
    return escapeHTML(text.slice(0, idx)) + "<mark>" + escapeHTML(text.slice(idx, idx + query.length)) + "</mark>" + escapeHTML(text.slice(idx + query.length));
  }

  function escapeHTML(s) {
    return String(s).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }

  function render(query) {
    var q = (query || "").trim().toLowerCase();
    var matches = entries.filter(matchesVectorFilter).filter(function (entry) {
      if (!q) return false;
      var haystack = (entry.title + " " + entry.summary + " " + (entry.tags || []).join(" ")).toLowerCase();
      return haystack.indexOf(q) !== -1;
    }).slice(0, 20);
    activeIndex = -1;
    if (!q) {
      results.hidden = true;
      results.innerHTML = "";
      return;
    }
    if (matches.length === 0) {
      results.hidden = false;
      results.innerHTML = '<p class="search-empty">No results.</p>';
      return;
    }
    results.hidden = false;
    results.innerHTML = matches.map(function (entry, i) {
      var vectors = (entry.vectors || []).join("/");
      return '<a href="' + escapeHTML(entryURL(entry)) + '" data-result-index="' + i + '">' +
        '<span class="search-result-title">' + highlight(entry.title, query) + '</span>' +
        (vectors ? '<span class="search-result-vector">' + escapeHTML(vectors) + '</span>' : '') +
        '<span class="search-result-summary">' + highlight(entry.summary, query) + '</span>' +
        '</a>';
    }).join("");
  }

  input.addEventListener("input", function () { render(input.value); });

  chips.forEach(function (chip) {
    chip.addEventListener("click", function () {
      var vector = chip.getAttribute("data-vector");
      activeVectors[vector] = !activeVectors[vector];
      chip.classList.toggle("active", !!activeVectors[vector]);
      render(input.value);
    });
  });

  input.addEventListener("keydown", function (event) {
    var items = results.querySelectorAll("a[data-result-index]");
    if (event.key === "ArrowDown" && items.length) {
      event.preventDefault();
      activeIndex = Math.min(activeIndex + 1, items.length - 1);
      items[activeIndex].focus();
    } else if (event.key === "ArrowUp" && items.length) {
      event.preventDefault();
      activeIndex = Math.max(activeIndex - 1, 0);
      items[activeIndex].focus();
    } else if (event.key === "Escape") {
      input.value = "";
      render("");
      input.blur();
    }
  });

  document.addEventListener("keydown", function (event) {
    var isShortcut = (event.key === "/" && document.activeElement !== input) ||
      ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k");
    if (isShortcut) {
      event.preventDefault();
      input.focus();
    }
  });

  document.addEventListener("click", function (event) {
    if (!root.contains(event.target)) {
      results.hidden = true;
    }
  });
})();
`
