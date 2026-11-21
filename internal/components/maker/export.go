package maker

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportFileStat records size metrics for one exported chat-loading Markdown file.
type ExportFileStat struct {
	Name  string
	Lines int
	Words int
}

// ExportResult summarizes a generated compact chat-loading pack.
type ExportResult struct {
	OutDir string
	Files  []ExportFileStat
}

// ExportChatPack condenses a validated Maker BatchInput into four token-efficient
// Markdown files ready to attach to a chat or LLM project:
//   - 00_BOOT_PROMPT.md
//   - 01_CORE_AND_RULES.md
//   - 02_CAST_AND_FACTIONS.md
//   - 03_TECH_AND_LOCATIONS.md
func ExportChatPack(batch BatchInput, outDir string) (ExportResult, error) {
	if strings.TrimSpace(outDir) == "" {
		return ExportResult{}, fmt.Errorf("maker export: out dir must not be empty")
	}
	if len(batch.EntityPackages) == 0 {
		return ExportResult{}, fmt.Errorf("maker export: batch has no entity_packages")
	}
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return ExportResult{}, err
	}

	visMap := indexVisualsByEntity(batch.VisualObservations)
	byType := groupEntitiesByType(batch.EntityPackages)

	files := []struct {
		name    string
		content string
	}{
		{"00_BOOT_PROMPT.md", renderBootPrompt(batch.ReportID)},
		{"01_CORE_AND_RULES.md", renderCoreAndRules(byType["other"], batch.EntityPackages)},
		{"02_CAST_AND_FACTIONS.md", renderSectionPack("Cast & Factions", []entityGroup{
			{"Factions", byType["faction"]},
			{"Characters (NPCs)", byType["npc"]},
		}, visMap)},
		{"03_TECH_AND_LOCATIONS.md", renderSectionPack("Tech, Items & Locations", []entityGroup{
			{"Items & Technology", byType["item"]},
			{"Locations & Layers", byType["location"]},
		}, visMap)},
	}

	res := ExportResult{OutDir: outDir, Files: make([]ExportFileStat, 0, len(files))}
	for _, f := range files {
		data := []byte(f.content)
		if err := os.WriteFile(filepath.Join(outDir, f.name), data, 0o644); err != nil { //nolint:gosec // G306: rendered Markdown export artifact.
			return ExportResult{}, err
		}
		res.Files = append(res.Files, ExportFileStat{
			Name:  f.name,
			Lines: countTextLines(data),
			Words: len(bytes.Fields(data)),
		})
	}
	return res, nil
}

// RenderExportSummary formats ExportResult for CLI display.
func RenderExportSummary(res ExportResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Maker Chat Pack Exported: %s\n\n", res.OutDir)
	totalLines, totalWords := 0, 0
	for _, f := range res.Files {
		totalLines += f.Lines
		totalWords += f.Words
		fmt.Fprintf(&b, "- `%s`: %d lines · %d words (~%d tokens)\n", f.Name, f.Lines, f.Words, estimateTokens(f.Words))
	}
	fmt.Fprintf(&b, "\n**Total:** %d files · %d lines · %d words (~%d tokens)\n", len(res.Files), totalLines, totalWords, estimateTokens(totalWords))
	return b.String()
}

type entityGroup struct {
	title string
	items []EntityPackage
}

func estimateTokens(words int) int {
	return (words*4 + 2) / 3
}

func indexVisualsByEntity(obs []VisualObservation) map[string][]VisualObservation {
	out := make(map[string][]VisualObservation)
	for _, o := range obs {
		if o.LinkedEntity != "" {
			out[o.LinkedEntity] = append(out[o.LinkedEntity], o)
		}
	}
	return out
}

func groupEntitiesByType(pkgs []EntityPackage) map[string][]EntityPackage {
	out := make(map[string][]EntityPackage)
	for _, p := range pkgs {
		out[p.EntityType] = append(out[p.EntityType], p)
	}
	return out
}

func renderBootPrompt(reportID string) string {
	var b strings.Builder
	b.WriteString("# Boot Prompt — Chat Loading Instructions\n\n")
	fmt.Fprintf(&b, "Pack ID: `%s`\n\n", reportID)
	b.WriteString("## Context Files\n\n")
	b.WriteString("1. `01_CORE_AND_RULES.md` — cosmology, phenomenology, timeline, rules, and open mysteries.\n")
	b.WriteString("2. `02_CAST_AND_FACTIONS.md` — factions, characters, capabilities, limits, and visual notes.\n")
	b.WriteString("3. `03_TECH_AND_LOCATIONS.md` — items, gear, technologies, and locations.\n\n")
	b.WriteString("## Operating Constraints\n\n")
	b.WriteString("- Treat these attached files as canonical truth; do not contradict stated facts.\n")
	b.WriteString("- Do not invent answers for items listed under **Open / Unresolved Mysteries**.\n")
	b.WriteString("- Keep visual notes strictly descriptive (appearance only; text governs facts).\n")
	b.WriteString("- Never assume decisive protagonist actions; pause and ask when a decision is pending.\n")
	return b.String()
}

func renderCoreAndRules(others, all []EntityPackage) string {
	var b strings.Builder
	b.WriteString("# Core Canon, Timeline & Rules\n\n")
	writeCompactTOC(&b, others)
	for _, p := range others {
		writeCompactEntity(&b, p, nil)
	}
	writeOpenMysteries(&b, all)
	return b.String()
}

func renderSectionPack(title string, groups []entityGroup, visMap map[string][]VisualObservation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", title)
	b.WriteString("## Index\n\n")
	for _, g := range groups {
		for _, p := range g.items {
			fmt.Fprintf(&b, "- `%s` (%s)\n", p.ID, g.title)
		}
	}
	b.WriteString("\n")
	for _, g := range groups {
		if len(g.items) == 0 {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n", g.title)
		for _, p := range g.items {
			writeCompactEntity(&b, p, visMap[p.ID])
		}
	}
	return b.String()
}

func writeCompactTOC(b *strings.Builder, pkgs []EntityPackage) {
	if len(pkgs) == 0 {
		return
	}
	b.WriteString("## Index\n\n")
	for _, p := range pkgs {
		fmt.Fprintf(b, "- `%s`\n", p.ID)
	}
	b.WriteString("\n")
}

func writeCompactEntity(b *strings.Builder, p EntityPackage, visuals []VisualObservation) {
	fmt.Fprintf(b, "### %s\n\n", p.ID)
	for _, f := range p.ExplicitFacts {
		fmt.Fprintf(b, "- **%s:** %s\n", f.Field, f.Value)
	}
	for _, inf := range p.InferredSuggestions {
		fmt.Fprintf(b, "- **%s** *(inferred, %s)*: %s\n", inf.Field, inf.Confidence, inf.Value)
	}
	writeVisualLine(b, visuals)
	b.WriteString("\n")
}

func writeVisualLine(b *strings.Builder, visuals []VisualObservation) {
	if len(visuals) == 0 {
		return
	}
	parts := make([]string, 0, len(visuals))
	for _, v := range visuals {
		if len(v.AppearanceFacts) > 0 {
			parts = append(parts, v.AppearanceFacts[0].Value)
		}
	}
	if len(parts) > 0 {
		fmt.Fprintf(b, "- **visual_reference:** %s\n", strings.Join(parts, " | "))
	}
}

func writeOpenMysteries(b *strings.Builder, all []EntityPackage) {
	var lines []string
	for _, p := range all {
		for _, u := range p.UnresolvedFields {
			lines = append(lines, fmt.Sprintf("- `%s` — **%s:** %s", p.ID, u.Field, u.Reason))
		}
	}
	if len(lines) == 0 {
		return
	}
	b.WriteString("## Open / Unresolved Mysteries (Do Not Invent)\n\n")
	for _, line := range lines {
		b.WriteString(line + "\n")
	}
}
