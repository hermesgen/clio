# Clio

[![Go Reference](https://pkg.go.dev/badge/github.com/adrianpk/clio.svg)](https://pkg.go.dev/github.com/adrianpk/clio)
[![codecov](https://codecov.io/gh/adrianpk/clio/branch/main/graph/badge.svg)](https://codecov.io/gh/adrianpk/clio)

![Main Index](docs/img/main-index.png)
<p align="right"><i><a href="docs/gallery/content.md">view gallery...</a></i></p>

Clio is a lightweight static site generator written in Go. It enables you to create and publish static content with a simple, direct workflow. The project is a work in progress.

## What you can do with Clio

- Write content in Markdown and generate clean HTML
- Preview your site locally before publishing
- Publish to platforms like GitHub Pages
- Create a simple personal blog
- Build sites with multiple sections and blogs
- Manage multiple independent sites from a single installation

Clio maintains a versioned record of both the Markdown source and the generated content, giving you a complete history of your site's evolution.

The workflow is straightforward: write in Markdown, preview, and publish when ready.

---

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests with HTML coverage report
make coverage-html

# Check if coverage meets 85% threshold
make coverage-check
```

### Quality Checks

```bash
# Run all quality checks (format, vet, test, coverage, lint)
make check

# Run CI pipeline (strict)
make ci
```

---

For a detailed feature list, see [features](docs/features.md).
For a complete overview of plans, see the [roadmap](docs/roadmap.md).

### Changelog

All notable changes to this project are documented in the [changelog file](docs/changelog.md).
