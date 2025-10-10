# Image Relationship System Redesign

## Current Problem
The existing system has inconsistent approaches for handling images:
- Header images stored as string paths in direct fields (`content.image`, `section.header`, etc.)
- Content images uploaded but not tracked in database
- No accessibility metadata (alt text, captions) stored
- No relationship tracking between entities and their images

## Proposed Architecture

### Unified Relationship Tables
Instead of mixing direct foreign keys for some images and loose files for others, we'll use intermediate relationship tables for all image associations.

### Database Schema

#### Relationship Tables
```sql
-- For content-image relationships
CREATE TABLE content_images (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL,
    image_id TEXT NOT NULL,
    purpose TEXT NOT NULL, -- 'header', 'content', 'thumbnail'
    position INTEGER, -- for ordering content images
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,

    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE(content_id, image_id, purpose)
);

-- For section-image relationships
CREATE TABLE section_images (
    id TEXT PRIMARY KEY,
    section_id TEXT NOT NULL,
    image_id TEXT NOT NULL,
    purpose TEXT NOT NULL, -- 'header', 'blog_header'
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,

    FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
    UNIQUE(section_id, image_id, purpose)
);
```

#### Updated Entity Tables
```sql
-- Remove direct image fields from content
ALTER TABLE content DROP COLUMN image; -- migrate data first

-- Remove direct image fields from section
ALTER TABLE section DROP COLUMN header;
ALTER TABLE section DROP COLUMN blog_header;
```

### Accessibility Metadata Storage Decision

**Option 1: Store in Image table**
- `alt_text` and `caption` live in the `images` table
- Simpler, one definition per image
- Most images have consistent meaning across contexts

**Option 2: Store in Relationship table**
- `alt_text` and `caption` in `content_images`/`layout_images`
- More flexible, context-specific descriptions
- More complex queries and management

### Benefits of This Approach

1. **Consistency**: All images handled the same way
2. **Accessibility**: Proper metadata storage for all images
3. **Flexibility**: Easy to add new image purposes
4. **History**: We can track image changes over time (Christmas headers, etc.)
5. **Scalability**: 1:N relationship naturally supports multiple images per content

### Image Purposes
- `header` - Main header image for content/section
- `blog_header` - Blog-specific header for sections
- `content` - Images used within markdown content
- `thumbnail` - Future: auto-generated thumbnails
- `seasonal` - Future: seasonal variations

### Migration Strategy

1. **Create new tables** with relationship structure
2. **Migrate existing data**: Convert current string paths to Image records + relationships
3. **Update models** to use relationships instead of direct fields
4. **Update ImageManager** to create proper Image records
5. **Update APIs** to work with new structure
6. **Update frontend** to handle new response formats
7. **Remove old fields** after migration complete
