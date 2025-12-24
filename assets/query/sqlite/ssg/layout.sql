-- Res: Layout
-- Table: layout

-- Create
INSERT INTO layout (
    id, site_id, short_id, name, description, code, created_by, updated_by, created_at, updated_at, header_image_id
) VALUES (
    :id, :site_id, :shortID, :name, :description, :code, :created_by, :updated_by, :created_at, :updated_at, :header_image_id
);

-- GetAll
SELECT * FROM layout;

-- Get
SELECT * FROM layout WHERE id = :id;

-- Update
UPDATE layout SET
    name = :name,
    description = :description,
    code = :code,
    updated_by = :updated_by,
    updated_at = :updated_at,
    header_image_id = :header_image_id
WHERE id = :id;

-- Delete
DELETE FROM layout WHERE id = :id;
