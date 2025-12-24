-- Res: ssg
-- Table: image
-- Create
INSERT INTO image (id, site_id, short_id, file_name, file_path, alt_text, title, width, height, created_by, updated_by, created_at, updated_at)
VALUES (:id, :site_id, :short_id, :file_name, :file_path, :alt_text, :title, :width, :height, :created_by, :updated_by, :created_at, :updated_at);

-- Res: ssg
-- Table: image
-- Get
SELECT id, site_id, short_id, file_name, file_path, alt_text, title, width, height, created_by, updated_by, created_at, updated_at
FROM image
WHERE id = ?;

-- Res: ssg
-- Table: image
-- GetImageByShortID
SELECT id, site_id, short_id, file_name, file_path, alt_text, title, width, height, created_by, updated_by, created_at, updated_at
FROM image
WHERE short_id = ?;

-- Res: ssg
-- Table: image
-- GetImageByContentHash
SELECT id, site_id, short_id, file_name, file_path, alt_text, title, width, height, created_by, updated_by, created_at, updated_at
FROM image
WHERE file_path = ?;

-- Res: ssg
-- Table: image
-- Update
UPDATE image
SET file_name = :file_name, file_path = :file_path, alt_text = :alt_text, title = :title, width = :width, height = :height, updated_by = :updated_by, updated_at = :updated_at
WHERE id = :id;

-- Res: ssg
-- Table: image
-- Delete
DELETE FROM image
WHERE id = ?;

-- Res: ssg
-- Table: image
-- List
SELECT id, site_id, short_id, file_name, file_path, alt_text, title, width, height, created_by, updated_by, created_at, updated_at
FROM image;
