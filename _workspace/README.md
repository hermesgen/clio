# Development Workspace

This directory is used for **development mode only**.

## Production vs Development

**Development mode** (current directory structure):
- Database: `_workspace/db/clio.db`
- Sites: `_workspace/sites/{slug}/`
- Everything in project root

**Production mode** (user's machine):
- Database: `~/.clio/clio.db`
- Configuration: `~/.config/clio/` (if needed)
- Sites: `~/Documents/Clio/sites/{slug}/`
- Workspace: `~/Documents/Clio/`

## Directory Structure

```
_workspace/
├── db/
│   └── clio.db              # Single unified database for all sites
└── sites/
    ├── structured/
    │   └── documents/
    │       ├── markdown/     # Source markdown files with frontmatter
    │       ├── assets/       # Static assets (images, etc.)
    │       └── html/         # Generated HTML (not versioned)
    └── blog/
        └── documents/
            ├── markdown/
            ├── assets/
            └── html/

```

## What's Versioned

This workspace is selectively version controlled for development purposes:

**Versioned** (source of truth):
- Database file (`db/clio.db`) - contains all content, sections, metadata
- Original seed images (`../assets/seed/images/`) - source images outside workspace
- Seeding scripts (`../scripts/seeding/`) - for workspace reconstruction

**NOT Versioned** (generated, can be reconstructed):
- Generated HTML files (`sites/*/documents/html/`)
- Markdown export files (`sites/*/documents/markdown/`)
- Workspace images (`sites/*/documents/assets/images/`)

## Setup After Fresh Clone

When you clone this repository fresh from the remote, you need to reconstruct the workspace. The database (`_workspace/db/clio.db`) is **already seeded** with all content and comes from version control, but the following need to be generated:

**What's missing after clone:**
- ❌ Workspace images (`sites/*/documents/assets/images/`)
- ❌ Markdown exports (`sites/*/documents/markdown/`)
- ❌ Generated HTML (`sites/*/documents/html/`)

**What's already there:**
- ✅ Database with all content (`_workspace/db/clio.db`)
- ✅ Source seed images (`assets/seed/images/`)

### Step-by-Step Setup

Follow these steps **in order** after cloning:

#### 1. Start the Application

```bash
make run
```

**What happens:**
- Compiles the Go application
- Runs database migrations (idempotent, safe to re-run)
- Starts three servers:
  - `8080` - Admin interface (editing)
  - `8081` - REST API
  - `8082` - Preview server (generated HTML)

**Important:** Keep this running in a terminal. The next steps need the API server (8081) running.

#### 2. Seed Images to Workspace

The database has image records but the physical files need to be copied from `assets/seed/images/` to the workspace:

```bash
# For structured site (required)
SITE_SLUG=structured DB_FILE=_workspace/db/clio.db go run scripts/seeding/seed-images.go

# For blog site (optional)
SITE_SLUG=blog DB_FILE=_workspace/db/clio.db go run scripts/seeding/seed-images.go
```

**What happens:**
- Reads image metadata from database
- Copies images from `assets/seed/images/{site}/` → `_workspace/sites/{site}/documents/assets/images/`
- Creates proper directory structure matching content slugs
- Links images to content/sections via database relationships

**Expected output:**
```
✓ Imported: image-name.png (1936x608, 1.5 MB)
⚠ Error inserting image: UNIQUE constraint failed  # ← OK if re-running
```

**Note:** UNIQUE constraint errors are normal if you re-run this script - it means images are already seeded.

#### 3. Generate HTML and Markdown

With the API server running (from step 1), generate the exports:

```bash
# For structured site
curl -X POST "http://localhost:8081/api/v1/ssg/generate-html" \
  -H "Content-Type: application/json" \
  -H "X-Site-Slug: structured"

# For blog site
curl -X POST "http://localhost:8081/api/v1/ssg/generate-html" \
  -H "Content-Type: application/json" \
  -H "X-Site-Slug: blog"
```

**What happens:**
- Exports markdown with frontmatter → `sites/{slug}/documents/markdown/`
- Generates static HTML → `sites/{slug}/documents/html/`
- Copies images to HTML static directory → `sites/{slug}/documents/html/static/images/`

**Expected output:**
```json
{"status":"success","message":"HTML generation process completed successfully"}
```

#### 4. Verify Everything Works

Open your browser:

- **Structured site preview:** http://structured.localhost:8082
- **Blog site preview:** http://blog.localhost:8082
- **Admin interface:** http://localhost:8080/ssg/

**Check that:**
- ✅ Images display correctly (no placeholders or 404s)
- ✅ All pages load without errors
- ✅ Navigation works between sections

### Quick Setup Script

For convenience, you can use the automated setup script:

```bash
# Sets up everything for a site in one command
./scripts/setup/setup-site.sh structured
./scripts/setup/setup-site.sh blog
```

This script combines all the steps above.

### Troubleshooting

**Images showing as placeholders or 404:**
- Run the image seeding script again (step 2)
- Verify files exist in `_workspace/sites/{slug}/documents/assets/images/`
- Check database has image records: `sqlite3 _workspace/db/clio.db "SELECT COUNT(*) FROM image;"`

**HTML generation fails:**
- Ensure API server is running on port 8081
- Check for errors in the terminal running `make run`
- Verify database exists: `ls -lh _workspace/db/clio.db`

**Port conflicts:**
- Check nothing is using ports 8080, 8081, 8082
- Use `ss -tlnp | grep "808"` to find conflicting processes
