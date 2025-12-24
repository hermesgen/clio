-- Res: ssg
-- Table: param
-- Create
INSERT INTO param (id, site_id, name, description, value, ref_key, system, created_by, updated_by, created_at, updated_at)
VALUES (:id, :site_id, :name, :description, :value, :ref_key, :system, :created_by, :updated_by, :created_at, :updated_at);

-- Res: ssg
-- Table: param
-- Get
SELECT id, site_id, name, description, value, ref_key, system, created_by, updated_by, created_at, updated_at
FROM param
WHERE id = ?;

-- Res: ssg
-- Table: param
-- GetByName
SELECT id, site_id, name, description, value, ref_key, system, created_by, updated_by, created_at, updated_at
FROM param
WHERE name = ?;

-- Res: ssg
-- Table: param
-- GetByRefKey
SELECT id, site_id, name, description, value, ref_key, system, created_by, updated_by, created_at, updated_at
FROM param
WHERE ref_key = ?;

-- Res: ssg
-- Table: param
-- List
SELECT id, site_id, name, description, value, ref_key, system, created_by, updated_by, created_at, updated_at
FROM param;

-- Res: ssg
-- Table: param
-- Update
UPDATE param
SET site_id = :site_id, name = :name, description = :description, value = :value, ref_key = :ref_key, updated_by = :updated_by, updated_at = :updated_at
WHERE id = :id;

-- Res: ssg
-- Table: param
-- Delete
DELETE FROM param
WHERE id = ?;
