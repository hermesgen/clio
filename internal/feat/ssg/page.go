package ssg

import (
	"html/template"

	"github.com/hermesgen/hm"
)

// PageData holds all the data needed to render a complete HTML page.
type PageData struct {
	HeaderStyle     string
	AssetPath       string
	Menu            []Section
	IsIndex         bool
	ListPageContent []Content
	Content         PageContent
	Blocks          *GeneratedBlocks
	Pagination      *PaginationData
	Config          *hm.Config // Esto lo quitaremos después de refactorizar el service y el template
	Search          SearchData // Nueva estructura para la configuración de búsqueda
}

// SearchData holds the configuration for the search functionality.
type SearchData struct {
	Provider string
	ID       string
	Enabled  bool
}

// PageContent holds the specific content to be rendered in the template for a single page.
type PageContent struct {
	Heading     string
	HeaderImage string
	Body        template.HTML
	Kind        string
}

// PaginationData holds data for rendering pagination controls.
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	NextPageURL string
	PrevPageURL string
}
