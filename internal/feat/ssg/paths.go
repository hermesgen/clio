package ssg

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetContentPath returns the URL path for a content item based on site mode.
// In structured mode: /{section-path}/{slug}/
// In blog mode: /{slug}/
func GetContentPath(content Content, mode string) string {
	slug := content.Slug()

	if mode == "blog" {
		// Blog mode: all content at root
		return filepath.Join("/", slug, "/")
	}

	// Structured mode: respect section paths
	return filepath.Join("/", content.SectionPath, slug, "/")
}

// GetIndexPath returns the URL path for an index based on content type and site mode.
// In structured mode:
//   - blog posts: /{section-path}/blog/ or /blog/ for root
//   - other types: /{section-path}/
// In blog mode: all content goes to /
func GetIndexPath(sectionPath string, contentType string, mode string) string {
	if mode == "blog" {
		// Blog mode: everything at root
		return "/"
	}

	// Structured mode
	basePath := strings.TrimSuffix(sectionPath, "/")

	// Special handling for blog content type
	if contentType == "blog" {
		blogPath := basePath + "/blog/"
		if sectionPath == "/" {
			blogPath = "/blog/"
		}
		return blogPath
	}

	// Non-blog content: use section path directly
	if sectionPath == "/" {
		return "/"
	}
	return basePath + "/"
}

// GetPaginationPath returns the URL path for a paginated index page.
// In blog mode: /page/{n}/
// In structured mode: /{index-path}/page/{n}/
func GetPaginationPath(indexPath string, page int, mode string) string {
	// Clean the index path
	indexPath = strings.TrimSuffix(indexPath, "/")

	if page == 1 {
		// First page is the index itself
		if indexPath == "" {
			return "/"
		}
		return indexPath + "/"
	}

	// Paginated pages
	if indexPath == "" {
		return fmt.Sprintf("/page/%d/", page)
	}
	return fmt.Sprintf("%s/page/%d/", indexPath, page)
}

// GetContentFilePath returns the filesystem path for a content HTML file.
// This is used for HTML generation.
func GetContentFilePath(htmlPath string, content Content, mode string) string {
	if mode == "blog" {
		return filepath.Join(htmlPath, content.Slug(), "index.html")
	}
	return filepath.Join(htmlPath, content.SectionPath, content.Slug(), "index.html")
}

// GetIndexFilePath returns the filesystem path for an index HTML file.
func GetIndexFilePath(htmlPath string, indexPath string) string {
	// Normalize "/" to "" for proper filepath.Join behavior
	if indexPath == "/" {
		indexPath = ""
	}
	return filepath.Join(htmlPath, indexPath, "index.html")
}

// GetPaginationFilePath returns the filesystem path for a paginated index HTML file.
func GetPaginationFilePath(htmlPath string, indexPath string, page int) string {
	if page == 1 {
		return GetIndexFilePath(htmlPath, indexPath)
	}
	// Normalize "/" to "" for proper filepath.Join behavior
	if indexPath == "/" {
		indexPath = ""
	}
	return filepath.Join(htmlPath, indexPath, "page", fmt.Sprintf("%d", page), "index.html")
}

// Multi-Site Path Helpers
// These functions return filesystem paths for a specific site.

// GetSiteBasePath returns the base path for a specific site.
// e.g., _workspace/sites/my-blog or ~/Documents/Clio/sites/my-blog
func GetSiteBasePath(sitesBasePath, siteSlug string) string {
	return filepath.Join(sitesBasePath, siteSlug)
}

// GetSiteDBPath returns the database file path for a specific site.
// e.g., _workspace/db/my-blog/clio.db or ~/.local/share/clio/db/my-blog/clio.db
func GetSiteDBPath(dbBasePath, siteSlug string) string {
	return filepath.Join(dbBasePath, siteSlug, "clio.db")
}

// GetSiteDBDSN returns the DSN for a specific site's database.
func GetSiteDBDSN(dbBasePath, siteSlug string) string {
	dbPath := GetSiteDBPath(dbBasePath, siteSlug)
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath)
}

// GetSiteDocsPath returns the documents path for a specific site.
// e.g., _workspace/sites/my-blog/documents
func GetSiteDocsPath(sitesBasePath, siteSlug string) string {
	return filepath.Join(sitesBasePath, siteSlug, "documents")
}

// GetSiteMarkdownPath returns the markdown path for a specific site.
func GetSiteMarkdownPath(sitesBasePath, siteSlug string) string {
	return filepath.Join(GetSiteDocsPath(sitesBasePath, siteSlug), "markdown")
}

// GetSiteHTMLPath returns the HTML output path for a specific site.
func GetSiteHTMLPath(sitesBasePath, siteSlug string) string {
	return filepath.Join(GetSiteDocsPath(sitesBasePath, siteSlug), "html")
}

// GetSiteAssetsPath returns the assets path for a specific site.
func GetSiteAssetsPath(sitesBasePath, siteSlug string) string {
	return filepath.Join(GetSiteDocsPath(sitesBasePath, siteSlug), "assets")
}

// GetSiteImagesPath returns the images path for a specific site.
func GetSiteImagesPath(sitesBasePath, siteSlug string) string {
	return filepath.Join(GetSiteAssetsPath(sitesBasePath, siteSlug), "images")
}
