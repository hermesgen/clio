package ssg

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Processor struct {
	parser goldmark.Markdown
}

// NewMarkdownProcessor creates and configures a new Markdown processor.
func NewMarkdownProcessor() *Processor {
	// Use enhanced renderer with empty context by default
	return NewMarkdownProcessorWithImageContext(&ImageContext{
		Images: make(map[string]ImageMetadata),
	})
}

// NewMarkdownProcessorWithImageContext creates a processor with image context for enhanced rendering.
func NewMarkdownProcessorWithImageContext(imageContext *ImageContext) *Processor {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			// Add extensions here, e.g., syntax.New()
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

// ToHTMLWithImageContext converts markdown to HTML and post-processes images with context
func (p *Processor) ToHTMLWithImageContext(markdown []byte, imageContext *ImageContext) (string, error) {
	// First convert markdown to HTML normally
	html, err := p.ToHTML(markdown)
	if err != nil {
		return "", err
	}

	// Post-process HTML to enhance images
	if imageContext != nil {
		html = enhanceImagesInHTML(html, imageContext)
	}

	return html, nil
}

// enhanceImagesInHTML post-processes HTML to enhance images with captions and metadata
func enhanceImagesInHTML(html string, imageContext *ImageContext) string {
	// Regex to match img tags with alt text containing pipe separator
	imgRegex := regexp.MustCompile(`<img([^>]*?)alt="([^"]*?)"([^>]*?)>`)

	// Process each img tag
	result := imgRegex.ReplaceAllStringFunc(html, func(match string) string {
		srcRegex := regexp.MustCompile(`src="([^"]*)"`)
		altRegex := regexp.MustCompile(`alt="([^"]*)"`)

		srcMatch := srcRegex.FindStringSubmatch(match)
		altMatch := altRegex.FindStringSubmatch(match)

		if len(srcMatch) < 2 || len(altMatch) < 2 {
			return match
		}

		srcValue := srcMatch[1]
		altValue := altMatch[1]

		var altText, longDescription string
		if strings.Contains(altValue, "|||") {
			parts := strings.SplitN(altValue, "|||", 2)
			altText = strings.TrimSpace(parts[0])
			longDescription = strings.TrimSpace(parts[1])
			fmt.Printf("[DEBUG] Found long description in HTML alt: alt='%s', longDesc='%s'\n", altText, longDescription)
		} else {
			altText = altValue
		}

		enhancedImg := fmt.Sprintf(`<img src="%s" alt="%s" class="prose-img">`, srcValue, altText)

		if longDescription != "" {
			return fmt.Sprintf(`<figure class="prose-figure">%s<figcaption class="prose-figcaption">%s</figcaption></figure>`, enhancedImg, longDescription)
		}

		return enhancedImg
	})

	return result
}
