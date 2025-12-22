package ssg_test

import (
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

func TestBuildIndexes(t *testing.T) {
	sections, content := setupIndexTestData(t)

	type expectedIndex struct {
		path            string
		contentCount    int
		orderedHeadings []string
	}

	type testCase struct {
		name            string
		content         []ssg.Content
		sections        []ssg.Section
		expectedIndexes []expectedIndex
	}

	testCases := []testCase{
		{
			name:     "Indexes for articles, blogs, and series",
			content:  content,
			sections: sections,
			expectedIndexes: []expectedIndex{
				{
					path:            "/",
					contentCount:    7,
					orderedHeadings: []string{"Tech Blog 2", "Article 2", "Go Series Part 2", "Go Series Part 1", "Tech Blog 1", "Go Series Part 3", "Article 1"},
				},
				{
					path:            "/news/",
					contentCount:    2,
					orderedHeadings: []string{"Article 2", "Article 1"},
				},
				{
					path:            "/tech/",
					contentCount:    5,
					orderedHeadings: []string{"Tech Blog 2", "Go Series Part 2", "Go Series Part 1", "Tech Blog 1", "Go Series Part 3"},
				},
				{
					path:            "/tech/blog/",
					contentCount:    2,
					orderedHeadings: []string{"Tech Blog 2", "Tech Blog 1"},
				},
				{
					path:            "/tech/go-series/",
					contentCount:    3,
					orderedHeadings: []string{"Go Series Part 1", "Go Series Part 2", "Go Series Part 3"},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexes := ssg.BuildIndexes(tc.content, tc.sections, "structured")

			if len(indexes) != len(tc.expectedIndexes) {
				t.Fatalf("Expected %d indexes, but got %d", len(tc.expectedIndexes), len(indexes))
			}

			actualIndexes := make(map[string]*ssg.Index)
			for _, idx := range indexes {
				actualIndexes[idx.Path] = idx
			}

			for _, expected := range tc.expectedIndexes {
				actual, ok := actualIndexes[expected.path]
				if !ok {
					t.Errorf("Expected index with path '%s' was not generated", expected.path)
					continue
				}

				if len(actual.Content) != expected.contentCount {
					t.Errorf("Index '%s': expected %d content items, got %d", expected.path, expected.contentCount, len(actual.Content))
				}

				if len(actual.Content) != len(expected.orderedHeadings) {
					t.Errorf("Index '%s': assertion error, content count (%d) does not match expected headings count (%d)", expected.path, len(actual.Content), len(expected.orderedHeadings))
					continue
				}

				for i, expectedHeading := range expected.orderedHeadings {
					if actual.Content[i].Heading != expectedHeading {
						t.Errorf("Index '%s' item %d: expected heading '%s', got '%s'", expected.path, i, expectedHeading, actual.Content[i].Heading)
					}
				}
			}
		})
	}
}

func setupIndexTestData(t *testing.T) (sections []ssg.Section, content []ssg.Content) {
	secRootID := uuid.New()
	secNewsID := uuid.New()
	secTechID := uuid.New()

	sections = []ssg.Section{
		{ID: secRootID, Name: "root", Path: "/"},
		{ID: secNewsID, Name: "news", Path: "/news/"},
		{ID: secTechID, Name: "tech", Path: "/tech/"},
	}

	now := time.Now()
	tArticle1 := now.Add(-1 * time.Hour) // Oldest
	tSeries3 := now.Add(-40 * time.Minute)
	tBlog1 := now.Add(-30 * time.Minute)
	tSeries1 := now.Add(-20 * time.Minute)
	tSeries2 := now.Add(-10 * time.Minute)
	tArticle2 := now.Add(-5 * time.Minute)
	tBlog2 := now.Add(-1 * time.Minute) // Newest

	content = []ssg.Content{
		{ID: uuid.New(), SectionID: secNewsID, Kind: "Article", Heading: "Article 1", PublishedAt: &tArticle1, SectionPath: "/news/"},
		{ID: uuid.New(), SectionID: secNewsID, Kind: "Article", Heading: "Article 2", PublishedAt: &tArticle2, SectionPath: "/news/"},
		{ID: uuid.New(), SectionID: secNewsID, Kind: "Page", Heading: "About News", SectionPath: "/news/"},

		// Tech Section
		{ID: uuid.New(), SectionID: secTechID, Kind: "Blog", Heading: "Tech Blog 1", PublishedAt: &tBlog1, SectionPath: "/tech/"},
		{ID: uuid.New(), SectionID: secTechID, Kind: "Blog", Heading: "Tech Blog 2", PublishedAt: &tBlog2, SectionPath: "/tech/"},
		{ID: uuid.New(), SectionID: secTechID, Kind: "Series", Series: "go-series", SeriesOrder: 1, Heading: "Go Series Part 1", PublishedAt: &tSeries1, SectionPath: "/tech/"},
		{ID: uuid.New(), SectionID: secTechID, Kind: "Series", Series: "go-series", SeriesOrder: 2, Heading: "Go Series Part 2", PublishedAt: &tSeries2, SectionPath: "/tech/"},
		{ID: uuid.New(), SectionID: secTechID, Kind: "Series", Series: "go-series", SeriesOrder: 3, Heading: "Go Series Part 3", PublishedAt: &tSeries3, SectionPath: "/tech/"},
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Path < sections[j].Path
	})

	return sections, content
}
