#!/bin/bash
set -e

SITE_SLUG=$1
DB_FILE=${2:-"_workspace/db/clio.db"}

if [ -z "$SITE_SLUG" ]; then
    echo "ERROR: SITE_SLUG not provided"
    echo "Usage: ./scripts/setup/setup-site.sh <site_slug> [db_file]"
    exit 1
fi

echo "Setting up site: $SITE_SLUG"
echo "Database: $DB_FILE"

# Clean assets and HTML
echo "Cleaning assets and HTML directories..."
rm -rf "_workspace/sites/$SITE_SLUG/documents/assets/images/"*
rm -rf "_workspace/sites/$SITE_SLUG/documents/html/"*

# Clean site-specific images from database
echo "Cleaning site images from database..."
sqlite3 "$DB_FILE" <<SQL
DELETE FROM content_images
WHERE content_id IN (
    SELECT id FROM content
    WHERE site_id = (SELECT id FROM site WHERE slug = '$SITE_SLUG')
);

DELETE FROM section_images
WHERE section_id IN (
    SELECT id FROM section
    WHERE site_id = (SELECT id FROM site WHERE slug = '$SITE_SLUG')
);

DELETE FROM image
WHERE site_id = (SELECT id FROM site WHERE slug = '$SITE_SLUG');
SQL

# Seed images
echo "Seeding images..."
SITE_SLUG="$SITE_SLUG" DB_FILE="$DB_FILE" go run scripts/seeding/seed-images.go

# Generate markdown
echo "Generating markdown..."
curl -s -X POST "http://localhost:8081/api/v1/ssg/generate-markdown" \
    -H "X-Site-Slug: $SITE_SLUG" > /dev/null

# Wait for markdown generation
sleep 3

# Regenerate HTML
echo "Regenerating HTML..."
curl -s -X POST "http://localhost:8081/api/v1/ssg/generate-html" \
    -H "X-Site-Slug: $SITE_SLUG" > /dev/null

echo "Site $SITE_SLUG setup complete!"
