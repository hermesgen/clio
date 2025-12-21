package ssg

import (
	"sort"

	"github.com/google/uuid"
)

// GeneratedBlocks holds all the pre-processed content lists for the blocks.
type GeneratedBlocks struct {
	// For Articles
	ArticleTagRelatedSameSection []Content
	ArticleRecentSameSection     []Content
	ArticleTagRelatedAllSections []Content
	ArticleRecentAllSections     []Content

	// For Blog
	BlogTagRelated []Content
	BlogRecent     []Content

	// For Series
	SeriesNext          *Content
	SeriesPrev          *Content
	SeriesIndexForward  []Content
	SeriesIndexBackward []Content
}

// BuildBlocks takes the current content and a list of all other content,
// and returns a GeneratedBlocks struct with all potential blocks pre-calculated.
func BuildBlocks(current Content, allContent []Content, maxItems int) *GeneratedBlocks {
	blocks := &GeneratedBlocks{}

	switch current.Kind {
	case "article":
		buildArticleBlocks(blocks, current, allContent, maxItems)
	case "blog":
		buildBlogBlocks(blocks, current, allContent, maxItems)
	case "series":
		buildSeriesBlocks(blocks, current, allContent, maxItems)
	}

	return blocks
}

func limit(content []Content, max int) []Content {
	if len(content) > max {
		return content[:max]
	}
	return content
}

func buildArticleBlocks(blocks *GeneratedBlocks, current Content, allContent []Content, maxItems int) {
	added := make(map[uuid.UUID]bool)
	added[current.ID] = true

	// Tag-Related Content (Same Section)
	for _, c := range allContent {
		if c.Kind == "article" && c.SectionID == current.SectionID && hasCommonTags(current, c) && !added[c.ID] {
			blocks.ArticleTagRelatedSameSection = append(blocks.ArticleTagRelatedSameSection, c)
			added[c.ID] = true
		}
	}

	// Recent Content (Same Section)
	for _, c := range allContent {
		if c.Kind == "article" && c.SectionID == current.SectionID && !added[c.ID] {
			blocks.ArticleRecentSameSection = append(blocks.ArticleRecentSameSection, c)
			added[c.ID] = true
		}
	}

	// Tag-Related Content (All Sections)
	for _, c := range allContent {
		if c.Kind == "article" && c.SectionID != current.SectionID && hasCommonTags(current, c) && !added[c.ID] {
			blocks.ArticleTagRelatedAllSections = append(blocks.ArticleTagRelatedAllSections, c)
			added[c.ID] = true
		}
	}

	// Recent Content (All Sections)
	for _, c := range allContent {
		if c.Kind == "article" && c.SectionID != current.SectionID && !added[c.ID] {
			blocks.ArticleRecentAllSections = append(blocks.ArticleRecentAllSections, c)
			added[c.ID] = true
		}
	}

	// Sort recent blocks by date
	sort.Slice(blocks.ArticleRecentSameSection, func(i, j int) bool {
		if blocks.ArticleRecentSameSection[i].PublishedAt == nil {
			return false
		}
		if blocks.ArticleRecentSameSection[j].PublishedAt == nil {
			return true
		}
		return blocks.ArticleRecentSameSection[i].PublishedAt.After(*blocks.ArticleRecentSameSection[j].PublishedAt)
	})
	sort.Slice(blocks.ArticleRecentAllSections, func(i, j int) bool {
		if blocks.ArticleRecentAllSections[i].PublishedAt == nil {
			return false
		}
		if blocks.ArticleRecentAllSections[j].PublishedAt == nil {
			return true
		}
		return blocks.ArticleRecentAllSections[i].PublishedAt.After(*blocks.ArticleRecentAllSections[j].PublishedAt)
	})

	// Apply limits
	blocks.ArticleTagRelatedSameSection = limit(blocks.ArticleTagRelatedSameSection, maxItems)
	blocks.ArticleRecentSameSection = limit(blocks.ArticleRecentSameSection, maxItems)
	blocks.ArticleTagRelatedAllSections = limit(blocks.ArticleTagRelatedAllSections, maxItems)
	blocks.ArticleRecentAllSections = limit(blocks.ArticleRecentAllSections, maxItems)
}

func hasCommonTags(c1, c2 Content) bool {
	for _, t1 := range c1.Tags {
		for _, t2 := range c2.Tags {
			if t1.ID == t2.ID {
				return true
			}
		}
	}
	return false
}

func buildBlogBlocks(blocks *GeneratedBlocks, current Content, allContent []Content, maxItems int) {
	added := make(map[uuid.UUID]bool)
	added[current.ID] = true

	for _, c := range allContent {
		if c.Kind == "blog" && c.SectionID == current.SectionID && hasCommonTags(current, c) && !added[c.ID] {
			blocks.BlogTagRelated = append(blocks.BlogTagRelated, c)
			added[c.ID] = true
		}
	}

	for _, c := range allContent {
		if c.Kind == "blog" && c.SectionID == current.SectionID && !added[c.ID] {
			blocks.BlogRecent = append(blocks.BlogRecent, c)
			added[c.ID] = true
		}
	}

	// Sort recent block by date
	sort.Slice(blocks.BlogRecent, func(i, j int) bool {
		if blocks.BlogRecent[i].PublishedAt == nil {
			return false
		}
		if blocks.BlogRecent[j].PublishedAt == nil {
			return true
		}
		return blocks.BlogRecent[i].PublishedAt.After(*blocks.BlogRecent[j].PublishedAt)
	})

	// Apply limits
	blocks.BlogTagRelated = limit(blocks.BlogTagRelated, maxItems)
	blocks.BlogRecent = limit(blocks.BlogRecent, maxItems)
}

func buildSeriesBlocks(blocks *GeneratedBlocks, current Content, allContent []Content, maxItems int) {
	if current.Kind != "series" || current.Series == "" {
		return // Not part of a series
	}

	var seriesPosts []*Content
	for i := range allContent {
		if allContent[i].Series == current.Series {
			seriesPosts = append(seriesPosts, &allContent[i])
		}
	}

	// Sort posts by series order
	sort.Slice(seriesPosts, func(i, j int) bool {
		return seriesPosts[i].SeriesOrder < seriesPosts[j].SeriesOrder
	})

	currentIndex := -1
	for i, p := range seriesPosts {
		if p.ID == current.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return // Should not happen if data is consistent
	}

	// Block 1: Simple Next/Previous
	if currentIndex > 0 {
		blocks.SeriesPrev = seriesPosts[currentIndex-1]
	}
	if currentIndex < len(seriesPosts)-1 {
		blocks.SeriesNext = seriesPosts[currentIndex+1]
	}

	// Block 2: Full Series Index
	if currentIndex < len(seriesPosts)-1 {
		// Convert []*Content to []Content for the result
		forwardContent := make([]Content, len(seriesPosts[currentIndex+1:]))
		for i, p := range seriesPosts[currentIndex+1:] {
			forwardContent[i] = *p
		}
		blocks.SeriesIndexForward = forwardContent
	}

	if currentIndex > 0 {
		previousPostsPtrs := seriesPosts[:currentIndex]
		// Reverse the slice for "closest first" order
		for i, j := 0, len(previousPostsPtrs)-1; i < j; i, j = i+1, j-1 {
			previousPostsPtrs[i], previousPostsPtrs[j] = previousPostsPtrs[j], previousPostsPtrs[i]
		}

		// Convert []*Content to []Content for the result
		backwardContent := make([]Content, len(previousPostsPtrs))
		for i, p := range previousPostsPtrs {
			backwardContent[i] = *p
		}
		blocks.SeriesIndexBackward = backwardContent
	}

	// Apply limits
	blocks.SeriesIndexForward = limit(blocks.SeriesIndexForward, maxItems)
	blocks.SeriesIndexBackward = limit(blocks.SeriesIndexBackward, maxItems)
}
