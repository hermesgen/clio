#!/bin/bash

# Description: Tests the Content Search and Pagination API endpoints.
# Usage: ./content-search.sh

# --- Configuration ---
source "$(dirname "$0")/_config.sh"
RESOURCE="contents"

# --- Main Functions ---

search_content() {
    local search_query=$1
    local page=${2:-1}
    echo "--- GET /$RESOURCE/search (Search Content: '$search_query', Page: $page) ---"
    local url="$BASE_URL/$RESOURCE/search?page=$page"
    if [ -n "$search_query" ]; then
        url="${url}&search=${search_query}"
    fi
    curl -s -X GET "$url" | jq .
}

paginated_content() {
    local page=${1:-1}
    echo "--- GET /$RESOURCE/search (Paginated Content - Page: $page) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/search?page=$page" | jq .
}

# --- Main Execution ---

echo "--- Running Content Search and Pagination Tests ---"

echo ""
echo "=== Test 1: Get first page of all content ==="
paginated_content 1
sleep 1

echo ""
echo "=== Test 2: Get second page of all content ==="
paginated_content 2
sleep 1

echo ""
echo "=== Test 3: Search for 'Indian' content ==="
search_content "Indian" 1
sleep 1

echo ""
echo "=== Test 4: Search for 'cuisine' content ==="
search_content "cuisine" 1
sleep 1

echo ""
echo "=== Test 5: Search for 'vegetarian' content ==="
search_content "vegetarian" 1
sleep 1

echo ""
echo "=== Test 6: Search with no query (should return all, paginated) ==="
search_content "" 1
sleep 1

echo ""
echo "=== Test 7: Search for non-existent content ==="
search_content "nonexistent" 1
sleep 1

echo ""
echo "=== Test 8: Test pagination with search results ==="
echo "First page of 'cuisine' results:"
search_content "cuisine" 1
sleep 1
echo "Second page of 'cuisine' results:"
search_content "cuisine" 2
sleep 1

echo "--- Content Search and Pagination Tests Finished ---"