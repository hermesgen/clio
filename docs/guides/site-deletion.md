# Site Deletion Guide

When you delete a site from Clio, the site is removed from the list but **your content files are preserved as backup**.

## What Gets Deleted

- The site entry from the Clio database
- The site will no longer appear in your sites list

## What Gets Preserved

All your content files remain intact in:

- **Markdown files**: `{sites_path}/{site-slug}/documents/markdown/`
- **Generated HTML**: `{sites_path}/{site-slug}/documents/html/`
- **Images**: `{sites_path}/{site-slug}/documents/assets/images/`
- **Site database**: `{db_path}/{site-slug}/clio.db`

## Default Locations

**Development mode:**
- Sites: `_workspace/sites/{site-slug}/`
- Databases: `_workspace/db/{site-slug}/`

**Production mode:**
- Sites: `~/Documents/Clio/sites/{site-slug}/` (or configured path)
- Databases: `~/.local/share/clio/db/{site-slug}/` (or configured path)

## Manual Cleanup

If you want to permanently remove all files for a deleted site:

```bash
# Remove site files
rm -rf {sites_path}/{site-slug}

# Remove site database
rm -rf {db_path}/{site-slug}
```

**Warning:** This action is permanent and cannot be undone.

## Restoring a Deleted Site

If you accidentally deleted a site and the files still exist, you can restore it:

1. Go to **Sites** → **New**
2. Use the **exact same slug** as before
3. Choose the same mode (structured or blog)
4. Clio will reconnect to the existing files

---

*This guide will be expanded with more details in future versions.*
