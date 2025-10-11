package ssg

import (
	"fmt"
	"strings"

	gmast "github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// ImageContext contains metadata about images for enhanced rendering
type ImageContext struct {
	Images map[string]ImageMetadata // key is the image path relative to /static/images/
}

// ImageMetadata holds accessibility and semantic information for an image
type ImageMetadata struct {
	AltText         string
	Caption         string
	LongDescription string
	Title           string
	Decorative      bool
}

// TailwindRenderer is a custom renderer for goldmark that adds Tailwind CSS classes.
type TailwindRenderer struct {
	html.Config
	ImageContext *ImageContext
}

// ImageRenderer is a simple renderer that only handles image nodes
type ImageRenderer struct {
	ImageContext *ImageContext
}

// NewTailwindRenderer creates a new TailwindRenderer with optional image context.
func NewTailwindRenderer(imageContext *ImageContext, opts ...html.Option) renderer.NodeRenderer {
	// imageCount := 0
	// if imageContext != nil && imageContext.Images != nil {
	//	imageCount = len(imageContext.Images)
	// }
	r := &TailwindRenderer{
		Config:       html.NewConfig(),
		ImageContext: imageContext,
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs registers the render functions for the nodes.
func (r *TailwindRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gmast.KindHeading, r.renderHeading)
	reg.Register(gmast.KindParagraph, r.renderParagraph)
	reg.Register(gmast.KindList, r.renderList)
	reg.Register(gmast.KindListItem, r.renderListItem)
	reg.Register(gmast.KindBlockquote, r.renderBlockquote)
	reg.Register(gmast.KindThematicBreak, r.renderHorizontalRule)
	reg.Register(gmast.KindEmphasis, r.renderEmphasis)
	reg.Register(extast.KindStrikethrough, r.renderDel)
	reg.Register(gmast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(gmast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(gmast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(extast.KindTable, r.renderTable)
	reg.Register(extast.KindTableHeader, r.renderTableHeader)
	reg.Register(extast.KindTableRow, r.renderTableRow)
	reg.Register(extast.KindTableCell, r.renderTableCell)
	reg.Register(gmast.KindLink, r.renderLink)
	reg.Register(gmast.KindImage, r.renderImage)
	reg.Register(gmast.KindText, r.renderText)
}

func (r *TailwindRenderer) renderHeading(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Heading)
	if entering {
		level := n.Level
		var class string
		switch level {
		case 1:
			class = "prose-h1"
		case 2:
			class = "prose-h2"
		case 3:
			class = "prose-h3"
		case 4:
			class = "prose-h4"
		case 5:
			class = "prose-h5"
		case 6:
			class = "prose-h6"
		}
		_, _ = w.WriteString(fmt.Sprintf("<h%d class=\"%s\">", level, class))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("</h%d>", n.Level))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderParagraph(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<p class=\"prose-p\">")
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderList(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.List)
	tag := "ul"
	if n.IsOrdered() {
		tag = "ol"
	}
	if entering {
		_, _ = w.WriteString(fmt.Sprintf("<%s class=\"prose-ul\">", tag))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("</%s>", tag))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderListItem(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<li class=\"prose-li\">")
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderBlockquote(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<blockquote class=\"prose-blockquote\">")
	} else {
		_, _ = w.WriteString("</blockquote>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderHorizontalRule(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<hr class=\"prose-hr\">\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderEmphasis(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Emphasis)
	if entering {
		if n.Level == 2 {
			_, _ = w.WriteString("<strong>")
		} else {
			_, _ = w.WriteString("<em>")
		}
	} else {
		if n.Level == 2 {
			_, _ = w.WriteString("</strong>")
		} else {
			_, _ = w.WriteString("</em>")
		}
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderDel(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<del>")
	} else {
		_, _ = w.WriteString("</del>")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderCodeSpan(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<code class=\"prose-code\">")
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*gmast.Text).Segment
			_, _ = w.Write(segment.Value(source))
		}
		_, _ = w.WriteString("</code>")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderCodeBlock(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre class=\"prose-pre\"><code>")
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}
		_, _ = w.WriteString("</code></pre>\n")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<pre class=\"prose-pre\"><code>")
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}
		_, _ = w.WriteString("</code></pre>\n")
	}
	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderTable(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table class=\"prose-table\">")
	} else {
		_, _ = w.WriteString("</table>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableHeader(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<thead>")
	} else {
		_, _ = w.WriteString("</thead>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableRow(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<tr>")
	} else {
		_, _ = w.WriteString("</tr>\n")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderTableCell(w util.BufWriter, source []byte, n gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<td>")
	} else {
		_, _ = w.WriteString("</td>")
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderLink(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Link)
	if entering {
		_, _ = w.Write([]byte(fmt.Sprintf("<a href=\"%s\" class=\"prose-a\">", n.Destination)))
	} else {
		_, _ = w.Write([]byte("</a>"))
	}
	return gmast.WalkContinue, nil
}

func (r *TailwindRenderer) renderImage(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Image)
	imgSrc := string(n.Destination)

	var altText, figCaption, figCaptionHTML string

	textAlt := ""
	if n.FirstChild() != nil {
		if text, ok := n.FirstChild().(*gmast.Text); ok {
			textAlt = string(text.Segment.Value(source))
		}
	}

	var markdownAltText, markdownLongDescription string
	if strings.Contains(textAlt, "|||") {
		parts := strings.SplitN(textAlt, "|||", 2)
		markdownAltText = strings.TrimSpace(parts[0])
		markdownLongDescription = strings.TrimSpace(parts[1])
	} else {
		markdownAltText = textAlt
	}

	hasMetadata := false
	if r.ImageContext != nil && r.ImageContext.Images != nil {
		imgPath := imgSrc
		imgPath = strings.TrimPrefix(imgPath, "/static/images/")
		imgPath = strings.TrimPrefix(imgPath, "/static/images")
		imgPath = strings.ReplaceAll(imgPath, "//", "/")
		imgPath = strings.TrimPrefix(imgPath, "/")

		if metadata, found := r.ImageContext.Images[imgPath]; found {
			hasMetadata = true
			altText = metadata.AltText
			figCaption = metadata.LongDescription

			if altText == "" {
				altText = textAlt
			}

			if figCaption != "" {
				figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
			}
		}
	}

	if !hasMetadata {
		altText = markdownAltText
		figCaption = markdownLongDescription

		if figCaption != "" {
			figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
		}
	} else {
		if markdownLongDescription != "" {
			figCaption = markdownLongDescription
			figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
		}
	}

	if entering {
		if figCaption != "" {
			_, _ = w.WriteString("<figure class=\"prose-figure\">")
		}

		_, _ = w.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"%s\" class=\"prose-img\">",
			imgSrc, altText))

		if figCaption != "" {
			_, _ = w.WriteString(figCaptionHTML + "</figure>")
		}
	}

	return gmast.WalkSkipChildren, nil
}

func (r *TailwindRenderer) renderText(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	if entering {
		n := node.(*gmast.Text)
		_, _ = w.Write(n.Segment.Value(source))
	}
	return gmast.WalkContinue, nil
}

// RegisterFuncs for ImageRenderer
func (r *ImageRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gmast.KindImage, r.renderImage)
}

func (r *ImageRenderer) renderImage(w util.BufWriter, source []byte, node gmast.Node, entering bool) (gmast.WalkStatus, error) {
	n := node.(*gmast.Image)
	imgSrc := string(n.Destination)

	var altText, figCaption, figCaptionHTML string

	textAlt := ""
	if n.FirstChild() != nil {
		if text, ok := n.FirstChild().(*gmast.Text); ok {
			textAlt = string(text.Segment.Value(source))
		}
	}

	var markdownAltText, markdownLongDescription string
	if strings.Contains(textAlt, "|||") {
		parts := strings.SplitN(textAlt, "|||", 2)
		markdownAltText = strings.TrimSpace(parts[0])
		markdownLongDescription = strings.TrimSpace(parts[1])
		fmt.Printf("[DEBUG] Found long description in alt text: alt='%s', longDesc='%s'\n", markdownAltText, markdownLongDescription)
	} else {
		markdownAltText = textAlt
	}

	hasMetadata := false
	if r.ImageContext != nil && r.ImageContext.Images != nil {
		imgPath := imgSrc
		imgPath = strings.TrimPrefix(imgPath, "/static/images/")
		imgPath = strings.TrimPrefix(imgPath, "/static/images")
		imgPath = strings.ReplaceAll(imgPath, "//", "/")
		imgPath = strings.TrimPrefix(imgPath, "/")

		if metadata, found := r.ImageContext.Images[imgPath]; found {
			hasMetadata = true
			altText = metadata.AltText
			figCaption = metadata.LongDescription

			if altText == "" {
				altText = textAlt
			}

			if figCaption != "" {
				figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
			}
		}
	}

	if !hasMetadata {
		altText = markdownAltText
		figCaption = markdownLongDescription

		if figCaption != "" {
			figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
		}
	} else {
		if markdownLongDescription != "" {
			figCaption = markdownLongDescription
			figCaptionHTML = fmt.Sprintf("<figcaption class=\"prose-figcaption\">%s</figcaption>", figCaption)
		}
	}

	if entering {
		if figCaption != "" {
			_, _ = w.WriteString("<figure class=\"prose-figure\">")
		}

		_, _ = w.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"%s\" class=\"prose-img\">",
			imgSrc, altText))

		if figCaption != "" {
			_, _ = w.WriteString(figCaptionHTML + "</figure>")
		}
	}

	return gmast.WalkSkipChildren, nil
}
