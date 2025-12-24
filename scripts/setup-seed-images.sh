#!/bin/bash
# Setup seed images for development environment
# This script copies seed images to workspace directories

set -e

echo "Setting up seed images for development..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Directories
SEED_IMAGES_DIR="assets/seed/images"
WORKSPACE_DIR="_workspace/sites"

# Setup structured site images
echo -e "${BLUE}Copying images for structured site...${NC}"
STRUCTURED_SEED="${SEED_IMAGES_DIR}/structured"
STRUCTURED_DEST="${WORKSPACE_DIR}/structured/documents/assets/images"

if [ -d "$STRUCTURED_SEED" ]; then
    mkdir -p "$STRUCTURED_DEST"
    cp -v "${STRUCTURED_SEED}"/*.png "$STRUCTURED_DEST/"
    echo -e "${GREEN}✓ Copied $(ls -1 ${STRUCTURED_SEED}/*.png | wc -l) images for structured site${NC}"
else
    echo "Warning: No seed images found for structured site"
fi

# Setup blog site images (placeholder for future)
echo -e "${BLUE}Checking blog site images...${NC}"
BLOG_SEED="${SEED_IMAGES_DIR}/blog"
BLOG_DEST="${WORKSPACE_DIR}/blog/documents/assets/images"

if [ -d "$BLOG_SEED" ]; then
    mkdir -p "$BLOG_DEST"
    cp -v "${BLOG_SEED}"/*.png "$BLOG_DEST/"
    echo -e "${GREEN}✓ Copied $(ls -1 ${BLOG_SEED}/*.png | wc -l) images for blog site${NC}"
else
    echo "Note: No seed images for blog site (will be generated later)"
fi

echo -e "${GREEN}✓ Seed images setup complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Start the server: make run"
echo "  2. Seed image database records: make seed-images-structured"
echo "  3. Regenerate HTML: make regenerate-html SITE=structured"
