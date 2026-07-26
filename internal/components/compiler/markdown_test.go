package compiler

import (
	"reflect"
	"testing"
)

func TestParseHeading(t *testing.T) {
	blocks := Parse("# Title\n")
	want := []Block{Heading{Level: 1, Inline: []InlineSpan{{Text: "Title"}}}}
	if !reflect.DeepEqual(blocks, want) {
		t.Fatalf("got %#v want %#v", blocks, want)
	}
}

func TestParseHeadingLevels(t *testing.T) {
	blocks := Parse("### Sub Title\n")
	heading, ok := blocks[0].(Heading)
	if !ok || heading.Level != 3 {
		t.Fatalf("got %#v want level 3 heading", blocks[0])
	}
}

func TestParseParagraph(t *testing.T) {
	blocks := Parse("Line one\nLine two continues.\n\nNext paragraph.\n")
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks want 2: %#v", len(blocks), blocks)
	}
	first, ok := blocks[0].(Paragraph)
	if !ok {
		t.Fatalf("blocks[0] is %T, want Paragraph", blocks[0])
	}
	if first.Inline[0].Text != "Line one Line two continues." {
		t.Fatalf("got %q", first.Inline[0].Text)
	}
}

func TestParseBoldItalicCode(t *testing.T) {
	blocks := Parse("Use **bold**, *italic*, and `code`.\n")
	para := blocks[0].(Paragraph)
	var gotBold, gotItalic, gotCode bool
	for _, span := range para.Inline {
		if span.Bold && span.Text == "bold" {
			gotBold = true
		}
		if span.Italic && span.Text == "italic" {
			gotItalic = true
		}
		if span.Code && span.Text == "code" {
			gotCode = true
		}
	}
	if !gotBold || !gotItalic || !gotCode {
		t.Fatalf("missing formatted spans: %#v", para.Inline)
	}
}

func TestParseLink(t *testing.T) {
	blocks := Parse("See [Armor](armor.md) for details.\n")
	para := blocks[0].(Paragraph)
	var found bool
	for _, span := range para.Inline {
		if span.LinkURL == "armor.md" && span.Text == "Armor" {
			found = true
		}
	}
	if !found {
		t.Fatalf("link span not found: %#v", para.Inline)
	}
}

func TestParseCodeFence(t *testing.T) {
	blocks := Parse("```text\nline one\nline two\n```\n")
	fence, ok := blocks[0].(CodeFence)
	if !ok {
		t.Fatalf("blocks[0] is %T, want CodeFence", blocks[0])
	}
	if fence.Lang != "text" {
		t.Fatalf("lang got %q want text", fence.Lang)
	}
	want := []string{"line one", "line two"}
	if !reflect.DeepEqual(fence.Lines, want) {
		t.Fatalf("lines got %#v want %#v", fence.Lines, want)
	}
}

func TestParseUnorderedList(t *testing.T) {
	blocks := Parse("- first\n- second\n- third\n")
	list, ok := blocks[0].(List)
	if !ok {
		t.Fatalf("blocks[0] is %T, want List", blocks[0])
	}
	if list.Ordered {
		t.Fatal("expected unordered list")
	}
	if len(list.Items) != 3 {
		t.Fatalf("got %d items want 3", len(list.Items))
	}
	if list.Items[0][0].Text != "first" {
		t.Fatalf("first item got %q", list.Items[0][0].Text)
	}
}

func TestParseOrderedList(t *testing.T) {
	blocks := Parse("1. first\n2. second\n")
	list, ok := blocks[0].(List)
	if !ok {
		t.Fatalf("blocks[0] is %T, want List", blocks[0])
	}
	if !list.Ordered {
		t.Fatal("expected ordered list")
	}
	if len(list.Items) != 2 {
		t.Fatalf("got %d items want 2", len(list.Items))
	}
}

func TestParseTable(t *testing.T) {
	source := "| Margin | Result |\n| ---: | --- |\n| 3 or more | strong hit |\n| 1 to 2 | hit |\n"
	blocks := Parse(source)
	table, ok := blocks[0].(Table)
	if !ok {
		t.Fatalf("blocks[0] is %T, want Table", blocks[0])
	}
	wantHeader := []string{"Margin", "Result"}
	if !reflect.DeepEqual(table.Header, wantHeader) {
		t.Fatalf("header got %#v want %#v", table.Header, wantHeader)
	}
	if len(table.Rows) != 2 {
		t.Fatalf("got %d rows want 2", len(table.Rows))
	}
	if table.Rows[0][0] != "3 or more" || table.Rows[0][1] != "strong hit" {
		t.Fatalf("row 0 got %#v", table.Rows[0])
	}
}

func TestParseHorizontalRule(t *testing.T) {
	blocks := Parse("Above.\n\n---\n\nBelow.\n")
	if len(blocks) != 3 {
		t.Fatalf("got %d blocks want 3: %#v", len(blocks), blocks)
	}
	if _, ok := blocks[1].(HorizontalRule); !ok {
		t.Fatalf("blocks[1] is %T, want HorizontalRule", blocks[1])
	}
}

func TestParseTableSeparatorIsNotHorizontalRule(t *testing.T) {
	source := "| A | B |\n| --- | --- |\n| 1 | 2 |\n"
	blocks := Parse(source)
	if len(blocks) != 1 {
		t.Fatalf("got %d blocks want 1 (a single Table): %#v", len(blocks), blocks)
	}
	if _, ok := blocks[0].(Table); !ok {
		t.Fatalf("blocks[0] is %T, want Table", blocks[0])
	}
}
