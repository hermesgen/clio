#!/bin/bash
SITE_SLUG="${1:-default}"
curl -i -X POST http://localhost:8081/api/v1/ssg/generate-markdown -H "X-Site-Slug: $SITE_SLUG"
