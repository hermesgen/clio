# Clio Roadmap

Clio was created to address a personal need for a lightweight, straightforward static site generator.  
Many existing tools are excellent, but often too generic, introducing friction when adapting to a specific workflow, or too constrained and slow due to their underlying stacks.  
Development follows a pragmatic, minimal-core approach, aiming for a practical middle ground.  
The goal is to reach a stable and maintainable version focused on core usefulness before adding non-essential features.  
Once the core milestones are complete, the project will move into maintenance mode: improving stability, expanding tests, and evolving only in response to concrete publishing needs or broader shifts in content workflows.

---

### Fundamental Features

- [x] Content CRUD **(Status: Completed)**
  Listing, search, pagination, view, edit, and delete operations.

- [x] Sections CRUD **(Status: Completed)**
  Listing, search, view, edit, and delete operations.

- [x] Layouts CRUD **(Status: Completed)**
  Listing, search, view, edit, and delete operations.

- [x] Parameters CRUD **(Status: Completed)**
  Dynamic configuration management, including GitHub credentials, repository and branch settings for Markdown versioning and publication.

- [x] Markdown generation **(Status: Completed)**
  Export content as Markdown files.

- [x] HTML generation **(Status: Completed)**
  Render Markdown content into HTML layouts.

- [x] Publishing to GitHub Pages **(Status: Completed)**
  Generate and deploy the site directly to GitHub Pages.

- [ ] Asset editing (images) **(Status: In Progress)**
  Upload and manage images for headers and inline content, with description and caption metadata.
  - Current: updating asset metadata requires re-uploading the file; CRUD exists but is incomplete.
  - Goal: edit asset metadata without replacing the file; keep Markdown usage in sync when applicable.

- [ ] Complete code coverage **(Status: Backlog)**
  Expand and refine the automated test suite for full functional and regression coverage.
  - Current tests are partial and oriented toward core operations.
  - Full coverage will follow once the feature set is stable, ensuring maintainability and future extensibility.

- [ ] robots.txt generation **(Status: Backlog)**
  Generate a standard `robots.txt` file at the site root, defining crawler access rules.

- [ ] Sitemap generation **(Status: Backlog)**
  Generate `/sitemap.xml` and optional index files such as `/sitemap_index.xml` or compressed `.xml.gz` variants.
  - Reference the sitemap automatically in `robots.txt`.
  - Keep the sitemap synchronized with published URLs.

- [ ] Feed generation **(Status: Backlog)**
  Generate `/rss.xml`, `/feed.xml`, and `/atom.xml` for content syndication.
  - Update automatically when new articles or posts are published.

---

### Complementary Features

- [ ] Improved HTML generation **(Status: Backlog)**
  Extend the HTML rendering pipeline to support richer formatting.
  - Advanced syntax highlighting using libraries like Prism.js.

- [ ] HTML backup and versioning **(Status: Backlog)**
  Maintain independent versioning of generated HTML in the repository.
  - The database remains the single source of truth, storing only the latest snapshot.
  - Repository history acts as both backup and record of content evolution.

- [ ] Instance regeneration from Markdown **(Status: Backlog)**
  Allow creating a new Clio instance from a versioned Markdown repository.
  - Rebuild database and layouts using the metadata and frontmatter stored in Markdown files.
  - Ensure compatibility between exported structure and re-import process.

- [ ] Optimized HTML generation **(Status: Backlog)**
  Generate only content that has changed since the last build to reduce processing time and unnecessary writes.

- [ ] Image variant generation **(Status: Backlog)**
  Generate optimized and resized variants of uploaded images.
  - Automatically produce lighter versions from high-resolution or low-compression originals.
  - Create aspect ratio variants for thumbnails and social media previews.
  - Progressive implementation: start with automatic header variants, then global automation, and finally customizable profiles for selective generation.

- [ ] Scheduled autopublication **(Status: Backlog)**
  Publish content automatically based on scheduled publish dates.
  - Periodic check (configurable interval) for items with `publish_at ≤ now`.
  - Timezone-aware; integrates with optimized builds.

- [x] Blog mode **(Status: Completed)**
  Enable a simplified `blog mode` where all content is treated as blog posts under the root path (`/`).
  - Blog mode filters index to show only blog-type content associated with the root section.
  - Content lives at `/{slug}/` instead of `/{section-path}/{slug}/`.
  - Site can switch between normal mode (multi-section) and blog mode (single chronological feed).
  - Intended for users who only need a single, continuous blog without sections or mixed content types.

- [ ] Tag-based navigation and indexes **(Status: Backlog)**
  Extend the current tagging system to generate browsable indexes and filtered views.
  - Automatically create index pages per tag, including pagination and search.
  - Link tags in rendered documents to their respective indexes.
  - Ensure consistency between tag metadata and generated site structure.

- [ ] User management for authorship **(Status: Backlog)**
  Introduce lightweight user management to support proper content attribution.
  - Maintain the single-user nature of Clio (no authentication or access control).
  - Allow defining multiple authors for attribution, collaboration, or republication contexts.
  - Support associating content with one or more authors and displaying author metadata in generated pages.

---

### Desirable Features

- [ ] Content import **(Status: Concept Stage)**
  Allow importing Markdown files from an external directory for those who prefer editing in their own environment (e.g., Neovim).
  - Support manual or automatic import modes.
  - Optional removal or hiding of source files once imported.
  - Smart mode: auto-import and hide unless the file mtime is newer than the last recorded version.

- [ ] Local API specification (OpenAPI) **(Status: Backlog)**
  Refine and document the existing internal API using the OpenAPI format.
  - Stabilize endpoints and define a versioning policy.
  - Publish minimal client examples (e.g., for a Neovim plugin).
  - **Application token:** external applications must authenticate with a locally issued token before interacting with the API; tokens are scoped and revocable.

- [ ] Image variant editing **(Status: Backlog)**
  Extend the asset manager to allow manual editing or replacement of generated image variants.
  - Enable overriding specific variants when automatic cropping or composition is unsuitable.
  - Support custom replacements to adapt visuals for different aspect ratios or platforms.
  - Maintain synchronization with metadata and variant naming to avoid inconsistencies.

- [ ] Layout integration and customization **(Status: Backlog)**
  Integrate user-defined layouts into the generation workflow, allowing them to override default templates.
  - Provide documentation describing the required structure and variables for valid layouts.
  - Define a convention-based discovery mechanism for user layouts.
  - Explore the inclusion of hooks or extension points for user-defined styles (CSS) or minor behavioral overrides.

- [ ] Tag manager **(Status: Backlog)**
  Provide a tag management interface to review, merge, and clean up tag data.
  - Detect duplicates and normalize casing (e.g., “Tag” vs “tag” vs “TAG”).
  - Allow merging tags and automatically reassigning associated content.
  - Identify and remove orphaned tags with no linked content.

- [ ] humans.txt generation **(Status: Backlog)**
  Generate a [`/humans.txt`](https://humanstxt.org/) file containing author and contributor information.
  - Optional feature for transparency and curiosity rather than indexing.
  - Editable through parameters or metadata.

---

*This roadmap is subject to refinement as Clio matures.*
