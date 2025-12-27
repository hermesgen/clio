# Clio - Static Site Generator
# Main Makefile for user-facing commands

# Variables
APP_NAME = clio
BUILD_DIR = build
SRC_DIR = .
MAIN_SRC = $(SRC_DIR)/main.go
BINARY = $(BUILD_DIR)/$(APP_NAME)
SITES_BASE = _workspace/sites
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

# Help target
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run all tests"
	@echo "  test-v           - Run tests with verbose output"
	@echo "  test-short       - Run tests in short mode"
	@echo "  coverage         - Run tests with coverage report"
	@echo "  coverage-html    - Generate HTML coverage report"
	@echo "  coverage-func    - Show function-level coverage"
	@echo "  coverage-check   - Check coverage meets 85% threshold"
	@echo "  coverage-100     - Check coverage is 100%"
	@echo "  coverage-summary - Display coverage table by package"
	@echo "  lint             - Run golangci-lint"
	@echo "  format           - Format code"
	@echo "  vet              - Run go vet"
	@echo "  check            - Run all quality checks (fmt, vet, test, coverage-check, lint)"
	@echo "  ci               - Run CI pipeline (strict, 100% coverage)"
	@echo "  clean            - Clean build and coverage files"

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
		pids=$$(lsof -ti :$$port 2>/dev/null); \
		if [ -n "$$pids" ]; then \
			for pid in $$pids; do \
				echo "Killing process $$pid on port $$port"; \
				kill -9 $$pid 2>/dev/null || true; \
			done; \
		fi; \
	done
	@pkill -9 clio 2>/dev/null || true
	@echo "Ports cleared."

# Run the application with environment variables
run: kill-ports setenv build
	@echo "Running $(APP_NAME) with environment variables..."
	@$(BINARY)

SITE ?= default

generate-markdown:
	@./scripts/curl/ssg/generate-markdown.sh $(SITE)

clean-html:
	@rm -rf _workspace/sites/$(SITE)/documents/html
	@mkdir -p _workspace/sites/$(SITE)/documents/html

generate-html:
	@./scripts/curl/ssg/generate-html.sh $(SITE)

regenerate-html: clean-html generate-html

publish:
	@./scripts/curl/ssg/publish.sh $(SITE)

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
	go test ./...

# Run tests with verbose output
test-v:
	go test -v ./...

# Run tests in short mode
test-short:
	go test -short ./...

# Run tests with coverage
coverage:
	go test -cover ./...

# Generate coverage profile and show percentage
coverage-profile:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

# Generate HTML coverage report
coverage-html: coverage-profile
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Show function-level coverage
coverage-func: coverage-profile
	go tool cover -func=coverage.out

# Check coverage percentage and fail if below threshold (85%)
coverage-check: coverage-profile
	@COVERAGE=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 85" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below 85% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets the 85% threshold"; \
	fi

# Check coverage percentage and fail if not 100%
coverage-100: coverage-profile
	@COVERAGE=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	if [ "$$COVERAGE" != "100.0" ]; then \
		echo "❌ Coverage $$COVERAGE% is not 100%"; \
		go tool cover -func=coverage.out | grep -v "100.0%"; \
		exit 1; \
	else \
		echo "🎉 Perfect! 100% test coverage achieved!"; \
	fi

# Display coverage summary table by package
coverage-summary:
	@echo "🧪 Running coverage tests by package..."
	@echo ""
	@echo "Coverage by package:"
	@echo "┌────────────────────────────────────────────────────────┬──────────┐"
	@echo "│ Package                                                │ Coverage │"
	@echo "├────────────────────────────────────────────────────────┼──────────┤"
	@for pkg in $$(go list ./... | grep -v "/build/"); do \
		pkgname=$$(echo $$pkg | sed 's|github.com/hermesgen/clio||' | sed 's|^/||'); \
		if [ -z "$$pkgname" ]; then pkgname="."; fi; \
		result=$$(go test -cover $$pkg 2>&1); \
		cov=$$(echo "$$result" | grep -oE '[0-9]+\.[0-9]+% of statements' | grep -v '^0\.0%' | tail -1 | grep -oE '[0-9]+\.[0-9]+%'); \
		if [ -z "$$cov" ]; then \
			if echo "$$result" | grep -qE '\[no test files\]|no test files'; then \
				cov="no tests"; \
			elif echo "$$result" | grep -q "FAIL"; then \
				cov="FAIL"; \
			else \
				cov="0.0%"; \
			fi; \
		fi; \
		printf "│ %-54s │ %8s │\n" "$$pkgname" "$$cov"; \
	done
	@echo "└────────────────────────────────────────────────────────┴──────────┘"

# Run go vet
vet:
	go vet ./...

# Run all quality checks
check: format vet test coverage-check lint
	@echo "✅ All quality checks passed!"

# CI pipeline - strict checks including 100% coverage
ci: format vet test coverage-100 lint
	@echo "🚀 CI pipeline passed!"

# Clean the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean -testcache
	@rm -f coverage.out coverage.html
	@echo "Clean complete."

# Phony targets
.PHONY: all build run setenv clean generate-markdown generate-html clean-html regenerate-html publish test test-v test-short coverage coverage-profile coverage-html coverage-func coverage-check coverage-100 coverage-summary vet check ci build-css kill-ports lint format
