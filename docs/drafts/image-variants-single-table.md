# Image Variants: Single Table Approach

## Current State

The system uses a single `images` table with purpose-based classification (`header`, `content`) and simple fallback logic in templates. The separate `image_variants` table exists but adds complexity without significant current benefit.

## Proposed Evolution

Replace the image_variants table with an enhanced single-table approach that supports automatic variant generation while maintaining simplicity.

### Enhanced Images Table Structure

```sql
ALTER TABLE images ADD COLUMN parent_id UUID REFERENCES images(id);
ALTER TABLE images ADD COLUMN original BOOLEAN DEFAULT true;
ALTER TABLE images ADD COLUMN width INTEGER;
ALTER TABLE images ADD COLUMN height INTEGER;
ALTER TABLE images ADD COLUMN aspect_ratio VARCHAR(10); -- "16:9", "1:1", etc
ALTER TABLE images ADD COLUMN meta JSONB; -- {target: "facebook", compression: 85, etc}
```

### Relationship Model

```
images:
├── header.jpg (purpose: header, original: true)
├── header_thumb.jpg (purpose: thumbnail, original: false, parent_id: header_id)  
├── header_social.jpg (purpose: social, meta: {target: "facebook"}, parent_id: header_id)
```

## Implementation Phases

### Phase 1: Current (Working)
- Single images table with purpose classification
- Template fallback logic (header → thumbnail → placeholder)
- No automatic generation

### Phase 2: Automatic Generation
- Optional automatic variant creation on image upload
- Thumbnail generation for header images
- Batch processing for existing images
- Maintains backward compatibility

### Phase 3: Manual Override
- Revive abandoned assets manager
- UI for overriding auto-generated variants
- Fine-grained control: "For Facebook, use this specific crop instead"

## Benefits

- **Simplicity**: One table instead of two
- **Flexibility**: Can generate variants without forcing them
- **Scalability**: Natural evolution path without breaking changes
- **Rich Metadata**: Supports aspect ratios, dimensions, target platforms
- **User Control**: Eventually allows manual override of auto-generated content

## Technical Advantages

- Fewer JOINs in queries
- Clear parent-child relationships via parent_id
- No sync issues between master and variants
- Optional generation means no forced complexity
- Template logic remains simple with fallbacks

This approach provides the benefits of image variants without the overhead of a separate table, supporting future optimization needs without current complexity.