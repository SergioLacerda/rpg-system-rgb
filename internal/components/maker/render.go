package maker

import (
	"fmt"
	"strings"
)

var tableCellEscaper = strings.NewReplacer("|", `\|`, "\r\n", " ", "\n", " ")

func escapeTableCell(s string) string {
	return tableCellEscaper.Replace(s)
}

// RenderEntityPackage renders an EntityPackage as Markdown matching
// skills/maker/templates/entity-package.md. Canonical Dependencies is omitted
// when empty so standalone story packages carry no external rule metadata.
func RenderEntityPackage(p EntityPackage) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Entity Package: %s\n\n", p.ID)
	fmt.Fprintf(&b, "**Type:** %s · **Status:** %s\n", p.EntityType, p.Status)
	fmt.Fprintf(&b, "**Source(s):** %s\n\n", joinOrNone(p.SourceRefs))
	writeExplicitFactsSection(&b, p.ExplicitFacts)
	writeInferredSection(&b, p.InferredSuggestions)
	writeUnresolvedSection(&b, p.UnresolvedFields)
	writeBulletSection(&b, "Conflicts", p.Conflicts)
	if len(p.CanonicalDeps) > 0 {
		writeBulletSection(&b, "Canonical Dependencies", p.CanonicalDeps)
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func writeExplicitFactsSection(b *strings.Builder, facts []ExplicitFact) {
	b.WriteString("## Explicit Facts\n\nStated directly in the source.\n\n")
	if len(facts) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	b.WriteString("| Field | Value | Source |\n|---|---|---|\n")
	for _, f := range facts {
		fmt.Fprintf(b, "| %s | %s | %s |\n", escapeTableCell(f.Field), escapeTableCell(f.Value), escapeTableCell(f.SourceRef))
	}
	b.WriteByte('\n')
}

func writeInferredSection(b *strings.Builder, items []InferredSuggestion) {
	b.WriteString("## Inferred Suggestions\n\nPlausible, but not stated directly — mark confidence honestly.\n\n")
	if len(items) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	b.WriteString("| Field | Value | Confidence | Rationale |\n|---|---|---|---|\n")
	for _, s := range items {
		fmt.Fprintf(b, "| %s | %s | %s | %s |\n", escapeTableCell(s.Field), escapeTableCell(s.Value), escapeTableCell(s.Confidence), escapeTableCell(s.Rationale))
	}
	b.WriteByte('\n')
}

func writeUnresolvedSection(b *strings.Builder, items []UnresolvedField) {
	b.WriteString("## Unresolved Fields\n\nFields this entity type normally has, that the source didn't cover.\n\n")
	if len(items) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	for _, u := range items {
		fmt.Fprintf(b, "- %s — %s\n", u.Field, u.Reason)
	}
	b.WriteByte('\n')
}

func writeBulletSection(b *strings.Builder, title string, items []string) {
	fmt.Fprintf(b, "## %s\n\n", title)
	if len(items) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	for _, item := range items {
		fmt.Fprintf(b, "- %s\n", item)
	}
	b.WriteByte('\n')
}

// RenderVisualObservation renders a VisualObservation as Markdown matching
// skills/maker/templates/visual-observation.md.
func RenderVisualObservation(o VisualObservation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Visual Observation: %s\n\n", o.ID)
	fmt.Fprintf(&b, "**Status:** %s\n", o.Status)
	fmt.Fprintf(&b, "**Source image:** %s\n", o.SourceRef)
	fmt.Fprintf(&b, "**Linked entity:** %s\n\n", orNone(o.LinkedEntity))
	b.WriteString("## Appearance\n\nWhat the image shows — visual properties only.\n\n")
	b.WriteString("| Field | Value | Confidence |\n|---|---|---|\n")
	for _, f := range o.AppearanceFacts {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", escapeTableCell(f.Field), escapeTableCell(f.Value), escapeTableCell(f.Confidence))
	}
	fmt.Fprintf(&b, "\n---\n\n> %s\n", strings.TrimSpace(o.NonAuthoritativeNote))
	return b.String()
}

// RenderConflictReport renders a ConflictReport as Markdown matching
// skills/maker/templates/conflict-report.md.
func RenderConflictReport(r ConflictReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Conflict Report: %s\n\n", r.ID)
	fmt.Fprintf(&b, "**Status:** %s\n", r.Status)
	fmt.Fprintf(&b, "**Field in dispute:** %s\n", r.Field)
	fmt.Fprintf(&b, "**Linked entity:** %s\n\n", orNone(r.LinkedEntity))
	b.WriteString("## Competing Claims\n\n| Value | Source | Source Kind |\n|---|---|---|\n")
	for _, c := range r.Claims {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", escapeTableCell(c.Value), escapeTableCell(c.SourceRef), escapeTableCell(c.SourceKind))
	}
	fmt.Fprintf(&b, "\n## Resolution\n\n**Status:** %s\n", r.Resolution)
	if strings.TrimSpace(r.ResolutionNote) != "" {
		fmt.Fprintf(&b, "\n%s\n", strings.TrimSpace(r.ResolutionNote))
	}
	return b.String()
}

// RenderMakerReport renders a MakerReport as Markdown matching
// skills/maker/templates/maker-report.md.
func RenderMakerReport(r MakerReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Maker Report: %s\n\n", r.ID)
	fmt.Fprintf(&b, "**Status:** %s\n", r.Status)
	fmt.Fprintf(&b, "**Input batch:** %s\n\n", joinOrNone(r.InputBatch))
	writeEntityListSection(&b, r.EntityPackages)
	writeBulletSection(&b, "Conflict Reports", r.ConflictReports)
	writeBulletSection(&b, "Visual Observations", r.VisualObservations)
	writeBulletSection(&b, "Open Questions For The GM", r.OpenQuestions)
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func writeEntityListSection(b *strings.Builder, ids []string) {
	b.WriteString("## Entity Packages Produced\n\n")
	if len(ids) == 0 {
		b.WriteString("No extractable entities found in this batch.\n\n")
		return
	}
	for _, id := range ids {
		fmt.Fprintf(b, "- %s\n", id)
	}
	b.WriteByte('\n')
}

func joinOrNone(items []string) string {
	if len(items) == 0 {
		return "None."
	}
	return strings.Join(items, ", ")
}

func orNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return "None."
	}
	return value
}
