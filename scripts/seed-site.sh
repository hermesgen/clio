#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SITE_SLUG=$1

if [ -z "$SITE_SLUG" ]; then
    echo -e "${RED}ERROR: Site slug required${NC}"
    echo "Usage: $0 <site-slug>"
    echo "Example: $0 structured"
    exit 1
fi

WORKSPACE_DIR="_workspace/sites/$SITE_SLUG"
DB_FILE="$WORKSPACE_DIR/db/clio.db"
SEED_SOURCE="assets/seed/sqlite/$SITE_SLUG"
SEED_TARGET="assets/seed/sqlite"

# Check if site directory exists
if [ ! -d "$WORKSPACE_DIR" ]; then
    echo -e "${RED}ERROR: Site directory not found: $WORKSPACE_DIR${NC}"
    echo "Available sites:"
    ls -1 _workspace/sites/ 2>/dev/null || echo "  (none)"
    exit 1
fi

# Check if seed file exists
if [ ! -d "$SEED_SOURCE" ]; then
    echo -e "${RED}ERROR: Seed directory not found: $SEED_SOURCE${NC}"
    echo "Available seed directories:"
    ls -1 assets/seed/sqlite/ | grep -v ".json" 2>/dev/null || echo "  (none)"
    exit 1
fi

SEED_COUNT=$(find "$SEED_SOURCE" -name "*-ssg-*.json" | wc -l)
if [ "$SEED_COUNT" -eq 0 ]; then
    echo -e "${RED}ERROR: No seed files found in $SEED_SOURCE${NC}"
    exit 1
fi

echo -e "${BLUE}=== Seeding $SITE_SLUG site ===${NC}"
echo ""

# Check for active ssg seeds in target directory
ACTIVE_SEEDS=$(find "$SEED_TARGET" -maxdepth 1 -name "*-ssg-*.json" 2>/dev/null | wc -l)
if [ "$ACTIVE_SEEDS" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found active SSG seed files in $SEED_TARGET${NC}"
    echo "These will be moved to backup before copying $SITE_SLUG seed:"
    find "$SEED_TARGET" -maxdepth 1 -name "*-ssg-*.json" -exec basename {} \;
    echo ""

    # Backup existing seeds
    BACKUP_DIR=".seed-backup/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    find "$SEED_TARGET" -maxdepth 1 -name "*-ssg-*.json" -exec mv {} "$BACKUP_DIR/" \;
    echo -e "${GREEN}✓${NC} Backed up existing seeds to: $BACKUP_DIR"
fi

# Copy site-specific seed to target
cp "$SEED_SOURCE"/*-ssg-*.json "$SEED_TARGET/"
echo -e "${GREEN}✓${NC} Copied $SITE_SLUG seed to $SEED_TARGET"

# Ensure db directory exists
mkdir -p "$WORKSPACE_DIR/db"
echo -e "${GREEN}✓${NC} Ensured db directory exists"

# Backup existing DB if present
if [ -f "$DB_FILE" ]; then
    BACKUP_FILE="${DB_FILE}.backup.$(date +%Y%m%d_%H%M%S)"
    cp "$DB_FILE" "$BACKUP_FILE"
    echo -e "${GREEN}✓${NC} Backed up existing DB to: $BACKUP_FILE"
fi

# Delete DB to trigger fresh seeding
rm -f "$DB_FILE"
echo -e "${GREEN}✓${NC} Deleted DB for site '${SITE_SLUG}'"
echo ""

# Print instructions
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Start the server: ${GREEN}make run${NC}"
echo "2. Navigate to: ${GREEN}http://localhost:8080/ssg/sites/switch?slug=${SITE_SLUG}${NC}"
echo "3. The seed data will be applied automatically on first access"
echo ""
echo -e "${BLUE}Info:${NC}"
echo "  Site: ${SITE_SLUG}"
echo "  Seed file: $(basename $(find "$SEED_TARGET" -maxdepth 1 -name "*-ssg-*.json"))"
echo "  Target DB: $DB_FILE"
echo ""
echo -e "${YELLOW}Note:${NC} After seeding is complete, you may want to clean up:"
echo "  ${GREEN}make clean-seed${NC}  # Removes active seed file"
echo "  ${GREEN}make restore-seed${NC}  # Restores backed up seeds"
