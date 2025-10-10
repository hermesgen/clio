package ssg

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Processor struct {
	parser goldmark.Markdown
}

// NewMarkdownProcessor creates and configures a new Markdown processor.
func NewMarkdownProcessor() *Processor {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			// Add extensions here, e.g., syntax.New()
		),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(NewTailwindRenderer(), 2000),
			),
		),
	)

	return &Processor{
		parser: md,
	}
}

// ToHTML converts a Markdown string to an HTML string.
func (p *Processor) ToHTML(markdown []byte) (string, error) {
	var buf bytes.Buffer
	if err := p.parser.Convert(markdown, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
