# Backlog

## Todo

- [ ] **Remove ImageVariant functionality**
  - Clean up handlers, templates, and database schema
  - Simplify image management
  - Redesign variants using same Image table/struct with additional metadata columns
  - Add resolution, context fields (e.g., purpose:social, context:x.com)
  - Implement automatic variant generation on image upload (default set of variations)
  - Add manual variant selection for users (e.g., "LinkedIn optimized but skip Facebook")

- [ ] **Integrate Prism.js for syntax highlighting**
  - Add to markdown code blocks in SSG templates
  - Configure for Go, JS, CSS, Shell

- [ ] **Fix image thumbnails URL resolution**
  - Images in list views don't load correctly
  - Need proper URL/path handling from variants to main image

- [ ] **Optimize Image/Assets CRUD functionality**
  - Fix thumbnail display issues in list views
  - Add proper edit form for image metadata (name, description, alt text, etc.)
  - Currently images can only be uploaded via content forms, need standalone editing
  - Add pagination and search to list-images (similar to list-content functionality)
  - Analyze feasibility of image file replacement/updating

- [ ] **Auth templates consistency**
  - Apply declarative menu pattern to Auth interfaces if needed
