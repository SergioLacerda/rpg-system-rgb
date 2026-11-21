package maker

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// WorkspaceStats summarizes raw inputs and filtered Maker outputs in a story
// workspace directory (sources/text, sources/images, packages/entities,
// packages/visual-observations, packages/conflicts).
type WorkspaceStats struct {
	RootDir                string
	RawTextFiles           int
	RawTextLines           int
	CatalogedImages        int
	ImagesByFolder         map[string]int
	EntityFiles            int
	EntityLines            int
	EntityWords            int
	VisualObservationFiles int
	VisualObservationLines int
	OpenConflictFiles      int
	ResolvedConflictFiles  int
}

// CollectWorkspaceStats inspects a Maker workspace directory and returns counts
// of cataloged images, raw text lines, filtered entity lines/words, visual
// observations, and open vs. resolved conflict reports.
func CollectWorkspaceStats(rootDir string) (WorkspaceStats, error) {
	if strings.TrimSpace(rootDir) == "" {
		return WorkspaceStats{}, fmt.Errorf("maker stats: dir must not be empty")
	}
	info, err := os.Stat(rootDir)
	if err != nil {
		return WorkspaceStats{}, err
	}
	if !info.IsDir() {
		return WorkspaceStats{}, fmt.Errorf("maker stats: %q is not a directory", rootDir)
	}

	stats := WorkspaceStats{
		RootDir:        rootDir,
		ImagesByFolder: map[string]int{},
	}
	stats.RawTextFiles, stats.RawTextLines, _ = countMarkdownDir(filepath.Join(rootDir, "sources", "text"))
	stats.CatalogedImages, stats.ImagesByFolder = countImagesTree(filepath.Join(rootDir, "sources", "images"))
	stats.EntityFiles, stats.EntityLines, stats.EntityWords = countMarkdownDir(filepath.Join(rootDir, "packages", "entities"))
	stats.VisualObservationFiles, stats.VisualObservationLines, _ = countMarkdownDir(filepath.Join(rootDir, "packages", "visual-observations"))
	stats.OpenConflictFiles, _, _ = countMarkdownDir(filepath.Join(rootDir, "packages", "conflicts"))
	stats.ResolvedConflictFiles, _, _ = countMarkdownDir(filepath.Join(rootDir, "packages", "conflicts", "resolved"))
	return stats, nil
}

// RenderWorkspaceStats formats WorkspaceStats as a concise, readable summary.
func RenderWorkspaceStats(s WorkspaceStats) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Maker Workspace Stats: %s\n\n", s.RootDir)
	fmt.Fprintf(&b, "- **Raw text sources (`sources/text/`):** %d files · %d lines\n", s.RawTextFiles, s.RawTextLines)
	fmt.Fprintf(&b, "- **Cataloged images (`sources/images/`):** %d files\n", s.CatalogedImages)
	fmt.Fprintf(&b, "- **Filtered entities (`packages/entities/`):** %d files · %d lines · %d words\n", s.EntityFiles, s.EntityLines, s.EntityWords)
	if s.RawTextLines > 0 {
		pct := float64(s.EntityLines) * 100.0 / float64(s.RawTextLines)
		fmt.Fprintf(&b, "- **Filtered / raw text ratio:** %.2f%% (%d filtered lines from %d raw lines)\n", pct, s.EntityLines, s.RawTextLines)
	}
	fmt.Fprintf(&b, "- **Visual observations (`packages/visual-observations/`):** %d files · %d lines\n", s.VisualObservationFiles, s.VisualObservationLines)
	fmt.Fprintf(&b, "- **Conflicts:** %d open (`packages/conflicts/`) · %d resolved (`packages/conflicts/resolved/`)\n", s.OpenConflictFiles, s.ResolvedConflictFiles)
	writeImageFolderBreakdown(&b, s.ImagesByFolder)
	return b.String()
}

func writeImageFolderBreakdown(b *strings.Builder, byFolder map[string]int) {
	if len(byFolder) == 0 {
		return
	}
	b.WriteString("\n## Cataloged Images by Folder\n\n")
	folders := make([]string, 0, len(byFolder))
	for folder := range byFolder {
		folders = append(folders, folder)
	}
	sort.Strings(folders)
	for _, folder := range folders {
		fmt.Fprintf(b, "- `%s`: %d\n", folder, byFolder[folder])
	}
}

func countMarkdownDir(dir string) (files, lines, words int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, 0
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(dir, entry.Name())) //nolint:gosec // G304: caller-supplied workspace directory.
		if readErr != nil {
			continue
		}
		files++
		lines += countTextLines(data)
		words += len(bytes.Fields(data))
	}
	return files, lines, words
}

func countImagesTree(root string) (int, map[string]int) {
	byFolder := map[string]int{}
	total := 0
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		relDir, relErr := filepath.Rel(root, filepath.Dir(path))
		if relErr != nil {
			relDir = "."
		}
		byFolder[filepath.ToSlash(relDir)]++
		total++
		return nil
	})
	return total, byFolder
}

func countTextLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	n := bytes.Count(data, []byte{'\n'})
	if data[len(data)-1] != '\n' {
		n++
	}
	return n
}
