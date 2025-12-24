#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Seeding Verification ===${NC}\n"

# Check sites in database
ADMIN_DB="_workspace/config/clioadmin.db"
echo -e "${YELLOW}Sites in Database:${NC}"
if [ -f "$ADMIN_DB" ]; then
    SITE_COUNT=$(sqlite3 "$ADMIN_DB" "SELECT COUNT(*) FROM site" 2>/dev/null || echo "0")
    echo -e "  Total sites registered: ${GREEN}$SITE_COUNT${NC}\n"
    sqlite3 "$ADMIN_DB" "SELECT '  ' || CASE WHEN active=1 THEN '✓' ELSE '✗' END || ' ' || slug || ' - ' || name || ' (' || mode || ')' FROM site ORDER BY slug" 2>/dev/null
    echo ""
else
    echo -e "${RED}✗${NC} Admin database not found: $ADMIN_DB\n"
fi

# Check structured site
STRUCTURED_DB="_workspace/sites/structured/db/clio.db"
echo -e "${YELLOW}Structured Site:${NC}"
if [ -f "$STRUCTURED_DB" ]; then
    echo -e "${GREEN}✓${NC} Database exists: $STRUCTURED_DB"

    CONTENT_COUNT=$(sqlite3 "$STRUCTURED_DB" "SELECT COUNT(*) FROM content" 2>/dev/null || echo "0")
    SECTION_COUNT=$(sqlite3 "$STRUCTURED_DB" "SELECT COUNT(*) FROM section" 2>/dev/null || echo "0")

    echo -e "  Content items: ${GREEN}$CONTENT_COUNT${NC} (expected: 43)"
    echo -e "  Sections: ${GREEN}$SECTION_COUNT${NC} (expected: 4)"

    if [ "$CONTENT_COUNT" -eq 43 ] && [ "$SECTION_COUNT" -eq 4 ]; then
        echo -e "  ${GREEN}✓ SEEDED CORRECTLY${NC}"
    else
        echo -e "  ${RED}✗ NOT SEEDED OR INCOMPLETE${NC}"
    fi

    echo -e "\n  Sections:"
    sqlite3 "$STRUCTURED_DB" "SELECT '    - ' || name || ' (' || path || ')' FROM section ORDER BY path" 2>/dev/null

    echo -e "\n  Content by type:"
    sqlite3 "$STRUCTURED_DB" "SELECT '    ' || kind || ': ' || COUNT(*) FROM content GROUP BY kind ORDER BY kind" 2>/dev/null
else
    echo -e "${RED}✗${NC} Database not found: $STRUCTURED_DB"
    echo -e "  ${YELLOW}→${NC} Run: make seed-structured && make run && navigate to structured site"
fi

echo ""

# Check blog site
BLOG_DB="_workspace/sites/blog/db/clio.db"
echo -e "${YELLOW}Blog Site:${NC}"
if [ -f "$BLOG_DB" ]; then
    echo -e "${GREEN}✓${NC} Database exists: $BLOG_DB"

    CONTENT_COUNT=$(sqlite3 "$BLOG_DB" "SELECT COUNT(*) FROM content" 2>/dev/null || echo "0")
    SECTION_COUNT=$(sqlite3 "$BLOG_DB" "SELECT COUNT(*) FROM section" 2>/dev/null || echo "0")

    echo -e "  Content items: ${GREEN}$CONTENT_COUNT${NC} (expected: 9)"
    echo -e "  Sections: ${GREEN}$SECTION_COUNT${NC} (expected: 1)"

    if [ "$CONTENT_COUNT" -eq 9 ] && [ "$SECTION_COUNT" -eq 1 ]; then
        echo -e "  ${GREEN}✓ SEEDED CORRECTLY${NC}"
    else
        echo -e "  ${RED}✗ NOT SEEDED OR INCOMPLETE${NC}"
    fi

    echo -e "\n  Sections:"
    sqlite3 "$BLOG_DB" "SELECT '    - ' || name || ' (' || path || ')' FROM section ORDER BY path" 2>/dev/null

    echo -e "\n  Content by type:"
    sqlite3 "$BLOG_DB" "SELECT '    ' || kind || ': ' || COUNT(*) FROM content GROUP BY kind ORDER BY kind" 2>/dev/null

    echo -e "\n  Blog posts:"
    sqlite3 "$BLOG_DB" "SELECT '    ' || ROW_NUMBER() OVER (ORDER BY created_at) || '. ' || heading FROM content ORDER BY created_at LIMIT 9" 2>/dev/null
else
    echo -e "${RED}✗${NC} Database not found: $BLOG_DB"
    echo -e "  ${YELLOW}→${NC} Run: make seed-blog && make run && navigate to blog site"
fi

echo -e "\n${BLUE}=== Summary ===${NC}"
STRUCTURED_EXISTS=$([ -f "$STRUCTURED_DB" ] && echo "yes" || echo "no")
BLOG_EXISTS=$([ -f "$BLOG_DB" ] && echo "yes" || echo "no")

if [ "$STRUCTURED_EXISTS" = "yes" ] && [ "$BLOG_EXISTS" = "yes" ]; then
    echo -e "${GREEN}✓ Both sites have databases${NC}"
else
    echo -e "${YELLOW}⚠ Some sites are not seeded yet${NC}"
    [ "$STRUCTURED_EXISTS" = "no" ] && echo -e "  ${RED}✗${NC} structured not seeded"
    [ "$BLOG_EXISTS" = "no" ] && echo -e "  ${RED}✗${NC} blog not seeded"
fi
