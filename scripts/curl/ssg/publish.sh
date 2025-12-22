#!/bin/bash
SITE_SLUG="${1:-default}"
API_URL="http://localhost:8081/api/v1/ssg/publish"

if [ -n "$COMMIT_MESSAGE" ]; then
  curl -X POST "${API_URL}" \
    -H "Content-Type: application/json" \
    -H "X-Site-Slug: $SITE_SLUG" \
    -d "{\"commit_message\": \"${COMMIT_MESSAGE}\"}"
else
  curl -X POST "${API_URL}" \
    -H "Content-Type: application/json" \
    -H "X-Site-Slug: $SITE_SLUG" \
    -d "{}"
fi