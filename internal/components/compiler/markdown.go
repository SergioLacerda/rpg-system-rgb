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
		switch {
		case matchAt(runes, i, "**"):
			end := indexOf(runes, i+2, "**")
			if end == -1 {
				plain.WriteRune(runes[i])
				i++
				continue
			}
			flushPlain()
			spans = append(spans, InlineSpan{Text: string(runes[i+2 : end]), Bold: true})
			i = end + 2
		case runes[i] == '`':
			end := indexOf(runes, i+1, "`")
			if end == -1 {
				plain.WriteRune(runes[i])
				i++
				continue
			}
			flushPlain()
			spans = append(spans, InlineSpan{Text: string(runes[i+1 : end]), Code: true})
			i = end + 1
		case runes[i] == '*':
			end := indexOf(runes, i+1, "*")
			if end == -1 {
				plain.WriteRune(runes[i])
				i++
				continue
			}
			flushPlain()
			spans = append(spans, InlineSpan{Text: string(runes[i+1 : end]), Italic: true})
			i = end + 1
		case runes[i] == '[':
			closeBracket := indexOf(runes, i+1, "]")
			if closeBracket != -1 && closeBracket+1 < len(runes) && runes[closeBracket+1] == '(' {
				closeParen := indexOf(runes, closeBracket+2, ")")
				if closeParen != -1 {
					flushPlain()
					spans = append(spans, InlineSpan{
						Text:    string(runes[i+1 : closeBracket]),
						LinkURL: string(runes[closeBracket+2 : closeParen]),
					})
					i = closeParen + 1
					continue
				}
			}
			plain.WriteRune(runes[i])
			i++
		default:
			plain.WriteRune(runes[i])
			i++
		}
	}
	flushPlain()
	return spans
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
