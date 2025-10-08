package ssg

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	hm "github.com/hermesgen/hm"
)

type Generator struct {
	hm.Core
}

func NewGenerator(params hm.XParams) *Generator {
	core := hm.NewCore("ssg-generator", params)
	g := &Generator{
		Core: core,
	}
	return g
}

func (g *Generator) Generate(contents []Content) error {
	g.Log().Info("Starting markdown generation")

	basePath := g.Cfg().StrValOrDef(hm.Key.SSGMarkdownPath, "_workspace/documents/markdown")

	for _, content := range contents {
		fileName := content.Slug() + ".md"
		filePath := filepath.Join(basePath, content.SectionPath, fileName)

		// --- Frontmatter Generation (Ordered) ---
		var frontMatter yaml.MapSlice

		frontMatter = append(frontMatter, yaml.MapItem{Key: "title", Value: content.Heading})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "slug", Value: content.Slug()})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "permalink", Value: ""}) // TODO: Construct full permalink

		// Taxonomy
		var tags []string
		for _, t := range content.Tags {
			tags = append(tags, t.Name)
		}
		if len(tags) > 0 {
			frontMatter = append(frontMatter, yaml.MapItem{Key: "tags", Value: tags})
		}
		frontMatter = append(frontMatter, yaml.MapItem{Key: "layout", Value: content.SectionName}) // Assuming layout is related to section

		// Status
		frontMatter = append(frontMatter, yaml.MapItem{Key: "draft", Value: content.Draft})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "featured", Value: content.Featured})

		// Content
		frontMatter = append(frontMatter, yaml.MapItem{Key: "excerpt", Value: content.Meta.Description}) // Using description as a stand-in
		frontMatter = append(frontMatter, yaml.MapItem{Key: "summary", Value: ""})                       // TODO: Add a summary field if needed
		frontMatter = append(frontMatter, yaml.MapItem{Key: "description", Value: content.Meta.Description})

		// Media
		frontMatter = append(frontMatter, yaml.MapItem{Key: "image", Value: ""})        // TODO: Add image field to a model
		frontMatter = append(frontMatter, yaml.MapItem{Key: "social-image", Value: ""}) // TODO: Add social-image field

		// Timestamps
		frontMatter = append(frontMatter, yaml.MapItem{Key: "published-at", Value: content.PublishedAt})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "created-at", Value: content.CreatedAt})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "updated-at", Value: content.UpdatedAt})

		// SEO
		frontMatter = append(frontMatter, yaml.MapItem{Key: "robots", Value: content.Meta.Robots})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "keywords", Value: content.Meta.Keywords})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "canonical-url", Value: content.Meta.CanonicalURL})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "sitemap", Value: content.Meta.Sitemap})

		// Page Config
		frontMatter = append(frontMatter, yaml.MapItem{Key: "table-of-contents", Value: content.Meta.TableOfContents})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "comments", Value: content.Meta.Comments})
		frontMatter = append(frontMatter, yaml.MapItem{Key: "share", Value: content.Meta.Share})

		// Localization
		frontMatter = append(frontMatter, yaml.MapItem{Key: "locale", Value: ""}) // TODO: Add locale field

		// --- End of Frontmatter ---

		yamlBytes, err := yaml.Marshal(frontMatter)
		if err != nil {
			g.Log().Error("Cannot marshal front matter", "error", err, "content_id", content.GetShortID())
			continue
		}

		fileContent := fmt.Sprintf("---\n%s---\n%s", string(yamlBytes), content.Body)

		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			g.Log().Error("Cannot create directory", "error", err, "path", dir)
			continue
		}

		if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
			g.Log().Error("Cannot write file", "error", err, "path", filePath)
			continue
		}

		g.Log().Debug("Generated file", "path", filePath)
	}

	g.Log().Info("Markdown generation finished")

	return nil
}
