
-- Create
INSERT INTO section (id, site_id, short_id, name, description, path, layout_id, layout_name, created_by, updated_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Update
UPDATE section SET
    name = :name,
    description = :description,
    path = :path,
    layout_id = :layout_id,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Get
SELECT s.*, l.name as layout_name FROM section s LEFT JOIN layout l ON s.layout_id = l.id WHERE s.id = ?;

-- GetAll
SELECT s.*, l.name as layout_name FROM section s LEFT JOIN layout l ON s.layout_id = l.id;

-- Delete
DELETE FROM section WHERE id = ?;
