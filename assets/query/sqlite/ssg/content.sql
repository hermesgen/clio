-- Res: Content
-- Table: content

-- Create
INSERT INTO content (
    id, short_id, user_id, section_id, heading, body, draft, featured, published_at, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :short_id, :user_id, :section_id, :heading, :body, :draft, :featured, :published_at, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT id, user_id, section_id, heading, body, draft, featured, published_at, short_id, created_by, updated_by, created_at, updated_at FROM content;

-- Get
SELECT id, user_id, section_id, heading, body, draft, featured, published_at, short_id, created_by, updated_by, created_at, updated_at FROM content WHERE id = :id;

-- Update
UPDATE content SET
    user_id = :user_id,
    section_id = :section_id,
    heading = :heading,
    body = :body,
    draft = :draft,
    featured = :featured,
    published_at = :published_at,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Delete
DELETE FROM content WHERE id = :id;

-- GetAllContentWithMeta
SELECT
    c.id, c.user_id, c.section_id, c.kind, c.heading, c.body, c.draft, c.featured, c.published_at, c.short_id,
    c.created_by, c.updated_by, c.created_at, c.updated_at,
    s.path AS section_path, s.name AS section_name,
    m.id AS meta_id, m.description, m.keywords, m.robots, m.canonical_url, m.sitemap, m.table_of_contents, m.share, m.comments,
    t.id AS tag_id, t.short_id AS tag_short_id, t.name AS tag_name, t.slug AS tag_slug,
    ci.image_id AS content_image_id, ci.purpose AS image_purpose, 
    i.file_path AS image_file_path
FROM
    content c
LEFT JOIN
    section s ON c.section_id = s.id
LEFT JOIN
    meta m ON c.id = m.content_id
LEFT JOIN
    content_tag ct ON c.id = ct.content_id
LEFT JOIN
    tag t ON ct.tag_id = t.id
LEFT JOIN
    content_images ci ON c.id = ci.content_id AND ci.is_active = TRUE AND ci.purpose IN ('header', 'thumbnail', 'content')
LEFT JOIN
    images i ON ci.image_id = i.id
ORDER BY
    c.published_at DESC, c.created_at DESC, c.id ASC;

-- GetContentWithPaginationAndSearch
SELECT
    c.id, c.user_id, c.section_id, c.kind, c.heading, c.body, c.draft, c.featured, c.published_at, c.short_id,
    c.created_by, c.updated_by, c.created_at, c.updated_at,
    COALESCE(s.path, '') AS section_path, COALESCE(s.name, '') AS section_name,
    COALESCE(m.id, '') AS meta_id, COALESCE(m.description, '') AS description, COALESCE(m.keywords, '') AS keywords, 
    COALESCE(m.robots, '') AS robots, COALESCE(m.canonical_url, '') AS canonical_url, COALESCE(m.sitemap, '') AS sitemap, 
    COALESCE(m.table_of_contents, 0) AS table_of_contents, COALESCE(m.share, 0) AS share, COALESCE(m.comments, 0) AS comments,
    COALESCE(t.id, '') AS tag_id, COALESCE(t.short_id, '') AS tag_short_id, COALESCE(t.name, '') AS tag_name, COALESCE(t.slug, '') AS tag_slug,
    COALESCE(ci.image_id, '') AS content_image_id, COALESCE(ci.purpose, '') AS image_purpose, 
    COALESCE(i.file_path, '') AS image_file_path
FROM
    content c
LEFT JOIN
    section s ON c.section_id = s.id
LEFT JOIN
    meta m ON c.id = m.content_id
LEFT JOIN
    content_tag ct ON c.id = ct.content_id
LEFT JOIN
    tag t ON ct.tag_id = t.id
LEFT JOIN
    content_images ci ON c.id = ci.content_id AND ci.is_active = TRUE AND ci.purpose IN ('header', 'thumbnail', 'content')
LEFT JOIN
    images i ON ci.image_id = i.id
WHERE
    (? = '' OR c.heading LIKE '%' || ? || '%')
ORDER BY
    c.published_at DESC, c.created_at DESC, c.id ASC
LIMIT ? OFFSET ?;

-- GetContentCountWithSearch
SELECT COUNT(DISTINCT c.id)
FROM content c
WHERE
    (? = '' OR c.heading LIKE '%' || ? || '%');
