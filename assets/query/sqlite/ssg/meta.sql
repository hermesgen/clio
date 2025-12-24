-- Res: Meta
-- Table: meta

-- Create
INSERT INTO meta (
    id, site_id, short_id, content_id, summary, excerpt, description, keywords, robots, canonical_url, sitemap, table_of_contents, share, comments, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :site_id, :short_id, :content_id, :summary, :excerpt, :description, :keywords, :robots, :canonical_url, :sitemap, :table_of_contents, :share, :comments, :created_by, :updated_by, :created_at, :updated_at
);

-- GetByContentID
SELECT * FROM meta WHERE content_id = :content_id;

-- Update
UPDATE meta SET
    description = :description,
    keywords = :keywords,
    robots = :robots,
    canonical_url = :canonical_url,
    sitemap = :sitemap,
    table_of_contents = :table_of_contents,
    share = :share,
    comments = :comments,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE content_id = :content_id;
