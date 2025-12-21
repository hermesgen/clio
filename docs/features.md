# Features

This document provides a detailed overview of Clio's current and planned features.

## Content Creation & Management

- **Markdown editing**: Write content in Markdown and automatically generate HTML
- **Draft importing**: Write in your own editor and import the text when ready
- **Database-backed**: Content, layouts, and sections are managed and catalogued in a database
- **Frontmatter support**: Metadata in Markdown is used for indexing in the database on import, and preserved on export so Clio instances can be regenerated from Markdown files alone

## Site Structure

- **Section support**: Each section has its own path (`/section-name`), with the root path as the main section
- **Customizable layouts**: Layout per section, with sensible defaults for those who prefer not to provide templates
- **Content tagging**: Organize and categorize content with tags

## Site Modes

Clio supports two operational modes:

- **Normal mode**: Multi-section site structure where content lives at `/{section-path}/{slug}/`
- **Blog mode**: Single chronological feed where all blog posts live at `/{slug}/`, filtering only blog-type content associated with the root section

## Content Types

Clio supports multiple content types with contextual features:

- **Pages**: Content units with an open purpose, more oriented toward structure or general communication. They don't necessarily follow a narrative or editorial logic, but instead provide a flexible container for different kinds of information.

- **Articles**: Content units with a discursive intent, meant to be read as complete pieces, whether essay, commentary, narrative, or reflection.

- **Blog**: Chronologically organized content; both root and sections can have their own blog (`/blog`, `/section-name/blog`). In blog mode, the site operates as a single chronological feed of blog posts.

- **Series**: Ordered posts that form a sequence (guides, tutorials, multi-part essays), with automatic navigation links backward/forward within the series.

## Contextual Features

- **Contextual blocks**: Each content type allows attaching context-aware blocks to reference related content, filtered by type, section, or multi-section configuration

## Preview & Publishing

- **Local preview**: Preview the generated site before publishing
- **Version control**: Both source Markdown and generated content are version controlled
- **GitHub Pages support**: Publish generated content directly to GitHub Pages (first supported target)

---

For the project roadmap and planned features, see [roadmap.md](roadmap.md).
