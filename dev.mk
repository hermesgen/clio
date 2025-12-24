# Development-only targets
# These are used by the developer during development and testing

# Unified database location
DB_FILE = _workspace/db/clio.db

# Backup database in current directory
backup-db:
	$(call backup_db,.)

# Reset database by moving it to backup directory
reset-db:
	@mkdir -p $(DB_BACKUP_DIR)
	$(call backup_db,$(DB_BACKUP_DIR))
	@echo "A fresh database will be created on next application start"

# Seed image relations for a specific site (run with server running)
# Usage: make seed-images-structured or make seed-images-blog
seed-images-structured:
	@echo "Seeding image relations for structured site..."
	@SITE_SLUG=structured $(MAKE) seed-images-site

seed-images-blog:
	@echo "Seeding image relations for blog site..."
	@SITE_SLUG=blog $(MAKE) seed-images-site

seed-images-site:
	@if [ -z "$(SITE_SLUG)" ]; then \
		echo "ERROR: SITE_SLUG not set"; \
		exit 1; \
	fi; \
	if [ ! -f "$(DB_FILE)" ]; then \
		echo "ERROR: Database not found: $(DB_FILE)"; \
		echo "Start the server first"; \
		exit 1; \
	fi; \
	if [ -f "scripts/seeding/seed-images.go" ]; then \
		echo "Running image seeding script for $(SITE_SLUG)..."; \
		SITE_SLUG=$(SITE_SLUG) DB_FILE=$(DB_FILE) go run scripts/seeding/seed-images.go; \
		echo "Images seeded successfully for $(SITE_SLUG)"; \
	else \
		echo "ERROR: scripts/seeding/seed-images.go not found"; \
		exit 1; \
	fi

# Reset and seed structured site (multi-section, multiple content types)
seed-structured:
	@./scripts/seed-site.sh structured

# Reset and seed blog site (single section, blog posts only)
seed-blog:
	@./scripts/seed-site.sh blog

# Complete setup for a site (clean, seed images, generate markdown & HTML)
setup-site:
	@if [ -z "$(SITE_SLUG)" ]; then \
		echo "ERROR: SITE_SLUG not set"; \
		echo "Usage: make setup-site SITE_SLUG=structured"; \
		exit 1; \
	fi
	@./scripts/setup/setup-site.sh $(SITE_SLUG) $(DB_FILE)

# Verify seeding status for both sites
verify-seeding:
	@./scripts/verify-seeding.sh

# Clean up active seed files after seeding
clean-seed:
	@echo "Removing active SSG seed files..."
	@rm -f assets/seed/sqlite/*-ssg-*.json 2>/dev/null || true
	@echo "✓ Cleaned up seed files"

# Restore backed up seed files
restore-seed:
	@if [ -d ".seed-backup" ]; then \
		LATEST=$$(ls -t .seed-backup/ | head -1); \
		if [ -n "$$LATEST" ]; then \
			echo "Restoring seeds from: .seed-backup/$$LATEST"; \
			cp .seed-backup/$$LATEST/*-ssg-*.json assets/seed/sqlite/ 2>/dev/null && \
			echo "✓ Seeds restored" || \
			echo "No seed files found in backup"; \
		else \
			echo "No seed backups found"; \
		fi; \
	else \
		echo "No .seed-backup directory found"; \
	fi

# Snapshot current images and their database relations
# Used to recreate image state for gallery captures
snapshot-images:
	@TIMESTAMP=$$(date +%Y%m%d_%H%M%S); \
	echo "Creating image snapshot: $$TIMESTAMP"; \
	mkdir -p .snapshots/$$TIMESTAMP; \
	if [ -d "_workspace/sites/structured/documents/assets/images" ]; then \
		cp -r _workspace/sites/structured/documents/assets/images .snapshots/$$TIMESTAMP/; \
		echo "Images copied to .snapshots/$$TIMESTAMP/images"; \
	else \
		echo "No images directory found"; \
	fi; \
	if [ -f "$(DB_FILE)" ]; then \
		sqlite3 $(DB_FILE) ".dump images" > .snapshots/$$TIMESTAMP/images_table.sql; \
		sqlite3 $(DB_FILE) ".dump content_images" > .snapshots/$$TIMESTAMP/content_images_table.sql; \
		echo "Database relations exported to .snapshots/$$TIMESTAMP/"; \
	else \
		echo "No database file found"; \
	fi; \
	echo "Snapshot created: $$TIMESTAMP"

# Restore images from a specific snapshot
# Used to recreate image state for gallery captures
restore-images:
	@echo "Available snapshots:"; \
	ls -1 .snapshots/ 2>/dev/null || echo "No snapshots found"; \
	read -p "Enter snapshot name: " snapshot; \
	if [ ! -d ".snapshots/$$snapshot" ]; then \
		echo "Snapshot $$snapshot not found"; \
		exit 1; \
	fi; \
	echo "Restoring images from snapshot: $$snapshot"; \
	if [ -d ".snapshots/$$snapshot/images" ]; then \
		rm -rf _workspace/sites/structured/documents/assets/images; \
		cp -r .snapshots/$$snapshot/images _workspace/sites/structured/documents/assets/; \
		echo "Images restored from snapshot"; \
	else \
		echo "No images found in snapshot"; \
	fi; \
	if [ -f ".snapshots/$$snapshot/images_table.sql" ] && [ -f ".snapshots/$$snapshot/content_images_table.sql" ]; then \
		echo "Restoring database relations..."; \
		sqlite3 $(DB_FILE) "DELETE FROM content_images; DELETE FROM images;"; \
		grep "^INSERT INTO" .snapshots/$$snapshot/images_table.sql | sqlite3 $(DB_FILE) 2>/dev/null || true; \
		grep "^INSERT INTO" .snapshots/$$snapshot/content_images_table.sql | sqlite3 $(DB_FILE) 2>/dev/null || true; \
		echo "Database relations restored"; \
	else \
		echo "No database relations found in snapshot"; \
	fi; \
	echo "Restore complete"

# List available image snapshots
list-snapshots:
	@echo "Available image snapshots:"; \
	ls -1 .snapshots/ 2>/dev/null || echo "No snapshots found"

# Set site mode to blog
set-blog-mode:
	@echo "Setting site mode to 'blog'..."
	@sqlite3 $(DB_FILE) "INSERT OR REPLACE INTO param (id, short_id, name, value, description, ref_key, system, created_by, updated_by, created_at, updated_at) SELECT COALESCE((SELECT id FROM param WHERE ref_key = 'site.mode'), lower(hex(randomblob(16)))), COALESCE((SELECT short_id FROM param WHERE ref_key = 'site.mode'), ''), 'Site Mode', 'blog', 'Site operation mode: structured (multi-section) or blog (single chronological feed)', 'site.mode', 0, COALESCE((SELECT created_by FROM param WHERE ref_key = 'site.mode'), '00000000000000000000000000000000'), '00000000000000000000000000000000', COALESCE((SELECT created_at FROM param WHERE ref_key = 'site.mode'), datetime('now')), datetime('now');" 2>/dev/null || echo "Database not ready yet, will set mode after server starts"
	@echo "Site mode set to 'blog'"

# Set site mode to structured
set-structured-mode:
	@echo "Setting site mode to 'structured'..."
	@sqlite3 $(DB_FILE) "INSERT OR REPLACE INTO param (id, short_id, name, value, description, ref_key, system, created_by, updated_by, created_at, updated_at) SELECT COALESCE((SELECT id FROM param WHERE ref_key = 'site.mode'), lower(hex(randomblob(16)))), COALESCE((SELECT short_id FROM param WHERE ref_key = 'site.mode'), ''), 'Site Mode', 'structured', 'Site operation mode: structured (multi-section) or blog (single chronological feed)', 'site.mode', 0, COALESCE((SELECT created_by FROM param WHERE ref_key = 'site.mode'), '00000000000000000000000000000000'), '00000000000000000000000000000000', COALESCE((SELECT created_at FROM param WHERE ref_key = 'site.mode'), datetime('now')), datetime('now');"
	@echo "Site mode set to 'structured'"

# Run in blog mode
run-blog: kill-ports setenv build
	@echo "Running $(APP_NAME) in blog mode..."
	@$(BINARY) & \
	SERVER_PID=$$!; \
	echo "Waiting for database to be ready..."; \
	sleep 3; \
	echo "Setting blog mode and converting posts..."; \
	sqlite3 $(DB_FILE) "INSERT OR REPLACE INTO param (id, short_id, name, value, description, ref_key, system, created_by, updated_by, created_at, updated_at) SELECT COALESCE((SELECT id FROM param WHERE ref_key = 'site.mode'), lower(hex(randomblob(16)))), COALESCE((SELECT short_id FROM param WHERE ref_key = 'site.mode'), ''), 'Site Mode', 'blog', 'Site operation mode: structured (multi-section) or blog (single chronological feed)', 'site.mode', 0, COALESCE((SELECT created_by FROM param WHERE ref_key = 'site.mode'), '00000000000000000000000000000000'), '00000000000000000000000000000000', COALESCE((SELECT created_at FROM param WHERE ref_key = 'site.mode'), datetime('now')), datetime('now');" 2>/dev/null; \
	sqlite3 $(DB_FILE) "UPDATE content SET kind = 'blog' WHERE heading IN ('My First Blog Post', 'Learning Go: Week One', 'Weekend Thoughts on Minimalism', 'Building My First SSG', 'Coffee and Code', 'Debugging Like a Pro', 'The Joy of Simple Solutions', 'Working Remotely: One Year In', 'What I''m Learning Next');" 2>/dev/null; \
	echo "Blog mode configured. Server running with PID $$SERVER_PID"; \
	wait $$SERVER_PID

# Run in structured mode (explicit)
run-structured: kill-ports setenv build set-structured-mode
	@echo "Restoring article types for structured mode..."
	@sqlite3 $(DB_FILE) "UPDATE content SET kind = 'article' WHERE heading IN ('My First Blog Post', 'Learning Go: Week One', 'Weekend Thoughts on Minimalism', 'Building My First SSG', 'Coffee and Code', 'Debugging Like a Pro', 'The Joy of Simple Solutions', 'Working Remotely: One Year In', \"What I'm Learning Next\");" 2>/dev/null || true
	@echo "Running $(APP_NAME) in structured mode..."
	@$(BINARY)

# Show current site mode
show-mode:
	@echo "Current site mode:"
	@sqlite3 $(DB_FILE) "SELECT value FROM param WHERE ref_key = 'site.mode';" 2>/dev/null || echo "Not set (defaults to 'structured')"

# Run with specific header styles (development testing)
run-stacked: kill-ports build
	@echo "Running with style: stacked"
	@$(BINARY) -ssg.header.style=stacked

run-overlay: kill-ports build
	@echo "Running with style: overlay"
	@$(BINARY) -ssg.header.style=overlay

run-boxed: kill-ports build
	@echo "Running with style: boxed"
	@CLIO_SSG_HEADER_STYLE=boxed $(BINARY)

run-text-only: kill-ports build
	@echo "Running with style: text-only"
	@CLIO_SSG_HEADER_STYLE=text-only $(BINARY)

# Run the application with command-line flags (development testing)
runflags: kill-ports build
	@echo "Running $(APP_NAME) with command-line flags..."
	@$(BINARY) -server.web.host=localhost -server.web.port=9080 -server.api.host=localhost -server.api.port=9081

# Generate a CSRF key (development utility)
gencsrfkey:
	@if command -v openssl >/dev/null 2>&1; then \
		echo "CSRF Key: $$(openssl rand -base64 32)"; \
	elif command -v dd >/dev/null 2>&1; then \
		echo "CSRF Key: $$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | base64)"; \
	else \
		echo "Neither openssl nor dd are available. Please install one of them."; \
		exit 1; \
	fi

# Generate migration (development utility)
new-migration:
	@read -p "Migration name: " name; \
	timestamp=$$(date +"%Y%m%d%H%M%S"); \
	kebab=$$(echo "$$name" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g' | sed -E 's/^-+|-+$$//g'); \
	filename="./assets/migration/sqlite/$${timestamp}-$${kebab}.sql"; \
	mkdir -p ./assets/migration/sqlite; \
	touch "$$filename"; \
	echo "Created $$filename"

.PHONY: backup-db reset-db seed-images-structured seed-images-blog seed-images-site seed-structured seed-blog verify-seeding clean-seed restore-seed snapshot-images restore-images list-snapshots set-blog-mode set-structured-mode run-blog run-structured show-mode run-stacked run-overlay run-boxed run-text-only runflags gencsrfkey new-migration
