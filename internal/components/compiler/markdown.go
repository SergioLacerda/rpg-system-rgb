package compiler

import (
	"regexp"
	"strings"
)

var (
	headingPattern        = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	orderedListPattern    = regexp.MustCompile(`^\d+\.\s+(.*)$`)
	tableSeparatorRow     = regexp.MustCompile(`^\|?[\s:|-]+\|?$`)
	horizontalRulePattern = regexp.MustCompile(`^-{3,}$`)
)

// Parse turns Markdown source into a flat sequence of Blocks. It covers
// exactly the construct set found across docs/core/en/** (UN-01 construct
// census): headings, paragraphs, code fences, flat lists, GFM pipe tables,
// horizontal rules, and inline bold/italic/code/links. It is a
// line-oriented scanner, not a general CommonMark parser.
//
//nolint:gocyclo // line-type dispatch switch: the branch count is the inherent shape of a line classifier, splitting it moves complexity rather than removing it.
func Parse(source string) []Block {
	lines := strings.Split(source, "\n")
	var blocks []Block
	i := 0
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])

		switch {
		case trimmed == "":
			i++
		case strings.HasPrefix(trimmed, "```"):
			var fence CodeFence
			fence, i = parseCodeFence(lines, i)
			blocks = append(blocks, fence)
		case headingPattern.MatchString(trimmed):
			match := headingPattern.FindStringSubmatch(trimmed)
			blocks = append(blocks, Heading{Level: len(match[1]), Inline: parseInline(match[2])})
			i++
		case isTableStart(lines, i):
			var table Table
			table, i = parseTable(lines, i)
			blocks = append(blocks, table)
		case horizontalRulePattern.MatchString(trimmed):
			blocks = append(blocks, HorizontalRule{})
			i++
		case isListItemLine(trimmed):
			var list List
			list, i = parseList(lines, i)
			blocks = append(blocks, list)
		default:
			var para Paragraph
			para, i = parseParagraph(lines, i)
			blocks = append(blocks, para)
		}
	}
	return blocks
}

func parseCodeFence(lines []string, start int) (CodeFence, int) {
	lang := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[start]), "```"))
	fence := CodeFence{Lang: lang, Lines: []string{}}
	i := start + 1
	for i < len(lines) {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
			return fence, i + 1
		}
		fence.Lines = append(fence.Lines, lines[i])
		i++
	}
	return fence, i
}

func isListItemLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "- ") || orderedListPattern.MatchString(trimmed)
}

func parseList(lines []string, start int) (List, int) {
	first := strings.TrimSpace(lines[start])
	ordered := orderedListPattern.MatchString(first)
	list := List{Ordered: ordered}
	i := start
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			break
		}
		if ordered {
			match := orderedListPattern.FindStringSubmatch(trimmed)
			if match == nil {
				break
			}
			list.Items = append(list.Items, parseInline(match[1]))
		} else {
			if !strings.HasPrefix(trimmed, "- ") {
				break
			}
			list.Items = append(list.Items, parseInline(strings.TrimPrefix(trimmed, "- ")))
		}
		i++
	}
	return list, i
}

// isTableStart reports whether lines[i] begins a GFM pipe table: the line
// contains a pipe and the following line is a separator row of only
// pipes, dashes, colons, and spaces.
func isTableStart(lines []string, i int) bool {
	if i+1 >= len(lines) {
		return false
	}
	current := strings.TrimSpace(lines[i])
	next := strings.TrimSpace(lines[i+1])
	if !strings.Contains(current, "|") {
		return false
	}
	return tableSeparatorRow.MatchString(next) && strings.Contains(next, "-")
}

func parseTable(lines []string, start int) (Table, int) {
	table := Table{Header: splitTableRow(lines[start])}
	i := start + 2 // header + separator
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || !strings.Contains(trimmed, "|") {
			break
		}
		table.Rows = append(table.Rows, splitTableRow(lines[i]))
		i++
	}
	return table, i
}

func splitTableRow(line string) []string {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	cells := strings.Split(trimmed, "|")
	result := make([]string, len(cells))
	for i, cell := range cells {
		result[i] = strings.TrimSpace(cell)
	}
	return result
}

func parseParagraph(lines []string, start int) (Paragraph, int) {
	var text []string
	i := start
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || isBlockStart(trimmed, lines, i) {
			break
		}
		text = append(text, trimmed)
		i++
	}
	return Paragraph{Inline: parseInline(strings.Join(text, " "))}, i
}

// isBlockStart reports whether trimmed begins a non-paragraph block, used
// by parseParagraph to know where a paragraph run ends.
func isBlockStart(trimmed string, lines []string, i int) bool {
	if i == 0 {
		return false
	}
	return strings.HasPrefix(trimmed, "```") ||
		headingPattern.MatchString(trimmed) ||
		horizontalRulePattern.MatchString(trimmed) ||
		isListItemLine(trimmed) ||
		isTableStart(lines, i)
}

// inlineMatcher tries to match one inline construct starting at runes[i]. It
// returns the parsed span, the index just past the match, and whether it
// matched at all.
type inlineMatcher func(runes []rune, i int) (InlineSpan, int, bool)

var inlineMatchers = []inlineMatcher{tryBold, tryCode, tryItalic, tryLink}

func tryBold(runes []rune, i int) (InlineSpan, int, bool) {
	if !matchAt(runes, i, "**") {
		return InlineSpan{}, i, false
	}
	end := indexOf(runes, i+2, "**")
	if end == -1 {
		return InlineSpan{}, i, false
	}
	return InlineSpan{Text: string(runes[i+2 : end]), Bold: true}, end + 2, true
}

func tryCode(runes []rune, i int) (InlineSpan, int, bool) {
	if runes[i] != '`' {
		return InlineSpan{}, i, false
	}
	end := indexOf(runes, i+1, "`")
	if end == -1 {
		return InlineSpan{}, i, false
	}
	return InlineSpan{Text: string(runes[i+1 : end]), Code: true}, end + 1, true
}

func tryItalic(runes []rune, i int) (InlineSpan, int, bool) {
	if runes[i] != '*' {
		return InlineSpan{}, i, false
	}
	end := indexOf(runes, i+1, "*")
	if end == -1 {
		return InlineSpan{}, i, false
	}
	return InlineSpan{Text: string(runes[i+1 : end]), Italic: true}, end + 1, true
}

func tryLink(runes []rune, i int) (InlineSpan, int, bool) {
	if runes[i] != '[' {
		return InlineSpan{}, i, false
	}
	closeBracket := indexOf(runes, i+1, "]")
	if closeBracket == -1 || closeBracket+1 >= len(runes) || runes[closeBracket+1] != '(' {
		return InlineSpan{}, i, false
	}
	closeParen := indexOf(runes, closeBracket+2, ")")
	if closeParen == -1 {
		return InlineSpan{}, i, false
	}
	return InlineSpan{
		Text:    string(runes[i+1 : closeBracket]),
		LinkURL: string(runes[closeBracket+2 : closeParen]),
	}, closeParen + 1, true
}

// parseInline scans a line of text into a sequence of InlineSpans, per the
// construct set found in docs/core/en/**: **bold**, *italic*, `code`, and
// [text](url). Spans do not nest.
func parseInline(text string) []InlineSpan {
	var spans []InlineSpan
	var plain strings.Builder
	flushPlain := func() {
		if plain.Len() > 0 {
			spans = append(spans, InlineSpan{Text: plain.String()})
			plain.Reset()
		}
	}

	runes := []rune(text)
	i := 0
	for i < len(runes) {
		span, next, ok := matchInline(runes, i)
		if !ok {
			plain.WriteRune(runes[i])
			i++
			continue
		}
		flushPlain()
		spans = append(spans, span)
		i = next
	}
	flushPlain()
	return spans
}

func matchInline(runes []rune, i int) (InlineSpan, int, bool) {
	for _, matcher := range inlineMatchers {
		if span, next, ok := matcher(runes, i); ok {
			return span, next, true
		}
	}
	return InlineSpan{}, i, false
}

func matchAt(runes []rune, i int, s string) bool {
	target := []rune(s)
	if i+len(target) > len(runes) {
		return false
	}
	for j, r := range target {
		if runes[i+j] != r {
			return false
		}
	}
	return true
}

func indexOf(runes []rune, from int, s string) int {
	target := []rune(s)
	for i := from; i+len(target) <= len(runes); i++ {
		match := true
		for j, r := range target {
			if runes[i+j] != r {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
