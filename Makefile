# Clio - Static Site Generator
# Main Makefile for user-facing commands

# Variables
APP_NAME = clio
BUILD_DIR = build
SRC_DIR = .
MAIN_SRC = $(SRC_DIR)/main.go
BINARY = $(BUILD_DIR)/$(APP_NAME)
DB_FILE = _workspace/db/clio.db
DB_BACKUP_DIR = bak

CSS_SOURCES = assets/static/css/prose.css assets/ssg/**/*.html assets/ssg/**/*.tmpl assets/static/css/main.css

# Backup database with timestamp (used by dev.mk)
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

# Include development targets
-include dev.mk

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

# Kill processes on application ports
kill-ports:
	@echo "Killing processes on ports 8080, 8081, 8082..."
	@for port in 8080 8081 8082; do \
		pid=$$(lsof -ti :$$port 2>/dev/null); \
		if [ -n "$$pid" ]; then \
			echo "Killing process $$pid on port $$port"; \
			kill -9 $$pid 2>/dev/null || true; \
		fi; \
	done
	@echo "Ports cleared."

# Run the application with environment variables
run: kill-ports setenv build
	@echo "Running $(APP_NAME) with environment variables..."
	@$(BINARY)

# Generate markdown files
generate-markdown:
	@echo "Triggering markdown generation..."
	@./scripts/curl/ssg/generate-markdown.sh

# Clean HTML output directory
clean-html:
	@echo "Cleaning HTML output directory..."
	@rm -rf _workspace/documents/html
	@mkdir -p _workspace/documents/html
	@echo "HTML directory cleaned"

# Generate html files
generate-html:
	@echo "Triggering HTML generation..."
	@./scripts/curl/ssg/generate-html.sh

# Clean and generate HTML (useful when switching modes)
regenerate-html: clean-html generate-html

# Publish site
publish:
	@echo "Publishing site..."
	@./scripts/curl/ssg/publish.sh

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

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Phony targets
.PHONY: all build run setenv clean generate-markdown generate-html clean-html regenerate-html publish test build-css kill-ports lint format
