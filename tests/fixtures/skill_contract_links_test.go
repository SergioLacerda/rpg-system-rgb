package fixtures

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var markdownLinkTarget = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

// mirrors: .analysis/pending/20260802-skill-contract-deterministic-validation-analysis.md
// TestSkillContractLinksResolve walks every skills/**/*.md file and checks
// that every relative Markdown link target actually exists on disk. This
// is the structural check the 20260802-refined-packages-closure-audit
// mission found missing: skills/maker/SKILL.md once referenced schema,
// template, prompt, and example files that did not exist yet, and nothing
// caught it. External links (http/https/mailto) are skipped.
func TestSkillContractLinksResolve(t *testing.T) {
	root := "../../skills"
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}

		body, readErr := os.ReadFile(path) //nolint:gosec // G304: path from WalkDir over the repo's own skills/ tree
		if readErr != nil {
			t.Fatal(readErr)
		}
		dir := filepath.Dir(path)

		for _, match := range markdownLinkTarget.FindAllStringSubmatch(string(body), -1) {
			target := match[1]
			if isExternalLink(target) {
				continue
			}
			resolved := filepath.Join(dir, filepath.FromSlash(target))
			if _, statErr := os.Stat(resolved); statErr != nil {
				t.Errorf("%s: link target %q does not resolve (looked for %s)", path, target, resolved)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func isExternalLink(target string) bool {
	for _, prefix := range []string{"http://", "https://", "mailto:"} {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return false
}
