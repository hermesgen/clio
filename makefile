# Variables
APP_NAME = clio
BUILD_DIR = build
SRC_DIR = .
MAIN_SRC = $(SRC_DIR)/main.go
BINARY = $(BUILD_DIR)/$(APP_NAME)
DB_FILE = _workspace/db/clio.db
DB_BACKUP_DIR = bak

CSS_SOURCES = assets/static/css/prose.css assets/ssg/**/*.html assets/ssg/**/*.tmpl assets/static/css/main.css

# Backup database with timestamp
define backup_db
	@if [ -f "$(DB_FILE)" ]; then \
		TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
		DB_FILENAME=$$(basename $(DB_FILE)); \
		NEW_NAME="$(1)/$${TIMESTAMP}-$${DB_FILENAME}"; \
		echo "Moving $(DB_FILE) to $${NEW_NAME}..."; \
		mv "$(DB_FILE)" "$${NEW_NAME}"; \
		echo "Database moved to $${NEW_NAME}"; \
	else \
		echo "Database file $(DB_FILE) not found"; \
	fi
endef

# Default target
all: build

# Build CSS
build-css:
	@echo "Building CSS..."
	@./scripts/build-css.sh

# Build the application
build: build-css
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) $(MAIN_SRC)
	@echo "Build complete: $(BINARY)"

# Run linter
lint:
	@echo "Running linter and fixing issues..."
	@golangci-lint run --fix

# Format code
format:
	@echo "Formatting code..."
	@gofmt -w .

# Run the application with environment variables
run: setenv build
	@echo "Running $(APP_NAME) with environment variables..."
	@$(BINARY)

# Run the application with command-line flags
runflags: build
	@echo "Running $(APP_NAME) with command-line flags..."
	@$(BINARY) -server.web.host=localhost -server.web.port=9080 -server.api.host=localhost -server.api.port=9081

# Run with specific header styles
run-stacked:
	@echo "Running with style: stacked"
	@$(BINARY) -ssg.header.style=stacked

run-overlay:
	@echo "Running with style: overlay"
	@$(BINARY) -ssg.header.style=overlay

run-boxed:
	@echo "Running with style: boxed"
	@CLIO_SSG_HEADER_STYLE=boxed $(BINARY)

run-text-only:
	@echo "Running with style: text-only"
	@CLIO_SSG_HEADER_STYLE=text-only $(BINARY)

# Generate markdown files
generate-markdown:
	@echo "Triggering markdown generation..."
	@./scripts/curl/ssg/generate-markdown.sh

# Generate html files
generate-html:
	@echo "Triggering HTML generation..."
	@./scripts/curl/ssg/generate-html.sh

# Publish site
publish:
	@echo "Publishing site..."
	@./scripts/curl/ssg/publish.sh

gencsrfkey:
	@if command -v openssl >/dev/null 2>&1; then \
		echo "CSRF Key: $$(openssl rand -base64 32)"; \
	elif command -v dd >/dev/null 2>&1; then \
		echo "CSRF Key: $$(dd if=/dev/urandom bs=32 count=1 2>/dev/null | base64)"; \
	else \
		echo "Neither openssl nor dd are available. Please install one of them."; \
		exit 1; \
	fi

# Set environment variables
# WIP: This is a workaround to be able to associate some styles to notifications and buttons but another approach will
# be used at the end.
setenv:
	@echo "Setting app environment to development..."
	@export CLIO_APP_ENV="dev"
	@echo "Environment variables set."
	@echo "Setting environment variables..."
	@export CLIO_SERVER_WEB_HOST=localhost
	@export CLIO_SERVER_WEB_PORT=8080
	@export CLIO_SERVER_API_HOST=localhost
	@export CLIO_SERVER_API_PORT=8081
	@export CLIO_SERVER_INDEX_ENABLED=true
	@echo "Setting a CSRF key..."
	@export CLIO_SEC_CSRF_KEY="NdZ7ULOe+NJ1bs5TzS51K+U4azOYQ6Wtv4CXlF6gJNM="
	@echo "Setting encryption key..."
	@export CLIO_SEC_ENCRYPTION_KEY="6ee4f00a50771711e34dad331fde0aaf92ef48e0357c3cf7abcdcaeb7a18fd2a"
	@echo "Setting notification styles..."
	@export CLIO_NOTIFICATION_SUCCESS_STYLE="bg-green-600 text-white px-4 py-2 rounded"
	@export CLIO_NOTIFICATION_INFO_STYLE="bg-blue-600 text-white px-4 py-2 rounded"
	@export CLIO_NOTIFICATION_WARN_STYLE="bg-yellow-600 text-white px-4 py-2 rounded"
	@export CLIO_NOTIFICATION_ERROR_STYLE="bg-red-600 text-white px-4 py-2 rounded"
	@export CLIO_NOTIFICATION_DEBUG_STYLE="bg-gray-600 text-white px-4 py-2 rounded"
	@echo "Setting button styles..."
	@export CLIO_BUTTON_STYLE_STANDARD="bg-gray-600 text-white px-4 py-2 rounded"
	@export CLIO_BUTTON_STYLE_BLUE="bg-blue-600 text-white px-4 py-2 rounded"
	@export CLIO_BUTTON_STYLE_RED="bg-red-600 text-white px-4 py-2 rounded"
	@export CLIO_BUTTON_STYLE_GREEN="bg-green-600 text-white px-4 py-2 rounded"
	@export CLIO_BUTTON_STYLE_YELLOW="bg-yellow-600 text-white px-4 py-2 rounded"
	@echo "Setting render errors..."
	@export CLIO_RENDER_WEB_ERRORS="true"
	@export CLIO_RENDER_API_ERRORS="true"
	@export CLIO_SSG_BLOCKS_MAXITEMS=5
	@export CLIO_SSG_INDEX_MAXITEMS=9
	@export CLIO_SSG_SEARCH_GOOGLE_ENABLED=true
	@export CLIO_SSG_SEARCH_GOOGLE_ID="94ad2c0b147c141fa"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate migration
new-migration:
	@read -p "Migration name: " name; \
	timestamp=$$(date +"%Y%m%d%H%M%S"); \
	kebab=$$(echo "$$name" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g' | sed -E 's/^-+|-+$$//g'); \
	filename="./assets/migration/sqlite/$${timestamp}-$${kebab}.sql"; \
	mkdir -p ./assets/migration/sqlite; \
	touch "$$filename"; \
	echo "Created $$filename"

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Backup database in current directory
backup-db:
	$(call backup_db,.)

# Reset database by moving it to backup directory
reset-db:
	@mkdir -p $(DB_BACKUP_DIR)
	$(call backup_db,$(DB_BACKUP_DIR))
	@echo "A fresh database will be created on next application start"

# Snapshot current images and their database relations
# Used to recreate image state for gallery captures
snapshot-images:
	@TIMESTAMP=$$(date +%Y%m%d_%H%M%S); \
	echo "Creating image snapshot: $$TIMESTAMP"; \
	mkdir -p .snapshots/$$TIMESTAMP; \
	if [ -d "_workspace/documents/assets/images" ]; then \
		cp -r _workspace/documents/assets/images .snapshots/$$TIMESTAMP/; \
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
		rm -rf _workspace/documents/assets/images; \
		cp -r .snapshots/$$snapshot/images _workspace/documents/assets/; \
		echo "Images restored from snapshot"; \
	else \
		echo "No images found in snapshot"; \
	fi; \
	if [ -f ".snapshots/$$snapshot/images_table.sql" ] && [ -f ".snapshots/$$snapshot/content_images_table.sql" ]; then \
		echo "Restoring database relations..."; \
		sqlite3 $(DB_FILE) "DELETE FROM content_images; DELETE FROM images;"; \
		sqlite3 $(DB_FILE) < .snapshots/$$snapshot/images_table.sql; \
		sqlite3 $(DB_FILE) < .snapshots/$$snapshot/content_images_table.sql; \
		echo "Database relations restored"; \
	else \
		echo "No database relations found in snapshot"; \
	fi; \
	echo "Restore complete"

# List available image snapshots
list-snapshots:
	@echo "Available image snapshots:"; \
	ls -1 .snapshots/ 2>/dev/null || echo "No snapshots found"

# Phony targets
.PHONY: all build run runflags setenv clean backup-db reset-db generate-markdown generate-html publish test run-stacked run-overlay run-boxed run-text-only build-css snapshot-images restore-images list-snapshots
