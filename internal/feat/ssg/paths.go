package ssg

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GetContentPath returns the URL path for a content item based on site mode.
// In normal mode: /{section-path}/{slug}/
// In blog mode: /{slug}/
func GetContentPath(content Content, mode string) string {
	slug := content.Slug()

	if mode == "blog" {
		// Blog mode: all content at root
		return filepath.Join("/", slug, "/")
	}

	// Normal mode: respect section paths
	return filepath.Join("/", content.SectionPath, slug, "/")
}

// GetIndexPath returns the URL path for an index based on content type and site mode.
// In normal mode:
//   - blog posts: /{section-path}/blog/ or /blog/ for root
//   - other types: /{section-path}/
// In blog mode: all content goes to /
func GetIndexPath(sectionPath string, contentType string, mode string) string {
	if mode == "blog" {
		// Blog mode: everything at root
		return "/"
	}

	// Normal mode
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
// In normal mode: /{index-path}/page/{n}/
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
