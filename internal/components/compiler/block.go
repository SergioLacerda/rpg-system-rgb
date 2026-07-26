package compiler

// Block is a top-level Markdown block element produced by Parse.
type Block interface {
	isBlock()
}

// Heading is a Markdown heading (# through ######).
type Heading struct {
	Level  int
	Inline []InlineSpan
}

// Paragraph is a run of regular text.
type Paragraph struct {
	Inline []InlineSpan
}

// CodeFence is a fenced code block (``` lang ... ```). Lines are stored
// verbatim, with no inline parsing.
type CodeFence struct {
	Lang  string
	Lines []string
}

// List is a flat (non-nested) ordered or unordered list, per the UN-01
// construct census: no nested lists appear anywhere in docs/core/en/**.
type List struct {
	Ordered bool
	Items   [][]InlineSpan
}

// Table is a GFM-style pipe table: one header row and any number of
// pipe-delimited body rows. Column alignment markers are accepted but not
// preserved (minimal scope: no docs/core/en/** table uses alignment for
// anything beyond visual formatting).
type Table struct {
	Header []string
	Rows   [][]string
}

// HorizontalRule is a Markdown thematic break (---).
type HorizontalRule struct{}

func (Heading) isBlock()        {}
func (Paragraph) isBlock()      {}
func (CodeFence) isBlock()      {}
func (List) isBlock()           {}
func (Table) isBlock()          {}
func (HorizontalRule) isBlock() {}

// InlineSpan is one run of inline text with consistent formatting. Spans do
// not nest (e.g. bold text cannot itself contain a link) — no construct in
// docs/core/en/** requires nested inline formatting.
type InlineSpan struct {
	Text    string
	Bold    bool
	Italic  bool
	Code    bool
	LinkURL string
}
