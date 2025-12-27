package main

import (
	"database/sql"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type ImageMapping struct {
	FileName   string
	ContentRef string
	SectionRef string // If set, this is a section image (use section name, e.g., "root")
	AltText    string
	Title      string
	Purpose    string // header, content, thumbnail, blog_header
}

func main() {
	// Mapeo de archivos a contenido (usando el heading exacto del seed)
	mappings := []ImageMapping{
		{
			FileName:   "building-kubernetes-operators.png",
			ContentRef: "Building Kubernetes Operators with Go",
			AltText:    "Kubernetes operators illustration",
			Title:      "Building Kubernetes Operators",
			Purpose:    "header",
		},
		{
			FileName:   "the-philosophy-of-movies.png",
			ContentRef: "The Philosophy of Movies: Deconstructing Cinematic Meaning",
			AltText:    "Cinema and philosophy",
			Title:      "The Philosophy of Movies",
			Purpose:    "header",
		},
		{
			FileName:   "food-section-culinary-world.png",
			SectionRef: "Food",
			AltText:    "Global culinary diversity and gastronomic exploration",
			Title:      "Food Section Header",
			Purpose:    "header",
		},
		{
			FileName:   "a-culinary-journey.png",
			ContentRef: "A Culinary Journey: Exploring the World of Food",
			AltText:    "Culinary journey through global cuisines",
			Title:      "A Culinary Journey",
			Purpose:    "header",
		},
		{
			FileName:   "the-art-of-sushi.png",
			ContentRef: "The Art of Sushi: A Japanese Culinary Tradition",
			AltText:    "Traditional Japanese sushi",
			Title:      "The Art of Sushi",
			Purpose:    "header",
		},
		{
			FileName:   "italian-cuisine.png",
			ContentRef: "Italian Cuisine: A Flavorful Feast",
			AltText:    "Italian dishes",
			Title:      "Italian Cuisine",
			Purpose:    "header",
		},
		{
			FileName:   "thai-street-food.png",
			ContentRef: "Thai Street Food: A Culinary Adventure",
			AltText:    "Thai street food scene",
			Title:      "Thai Street Food",
			Purpose:    "header",
		},
		{
			FileName:   "mexican-delight.png",
			ContentRef: "Mexican Delights: Flavor and Tradition",
			AltText:    "Mexican food",
			Title:      "Mexican Delights",
			Purpose:    "header",
		},
		{
			FileName:   "the-secret-to-perfect-italian-pasta.png",
			ContentRef: "The Secret to Perfect Italian Pasta: A Simple Guide",
			AltText:    "Perfect Italian pasta",
			Title:      "The Secret to Perfect Italian Pasta",
			Purpose:    "header",
		},
		{
			FileName:   "exploring-indian-vegetarian-cuisine.png",
			ContentRef: "Exploring Indian Vegetarian Cuisine: A World of Spices",
			AltText:    "Indian vegetarian dishes",
			Title:      "Exploring Indian Vegetarian Cuisine",
			Purpose:    "header",
		},
		{
			FileName:   "index-header.png",
			SectionRef: "root",
			AltText:    "Welcome to the root section",
			Title:      "Index Header",
			Purpose:    "header",
		},
		{
			FileName:   "tech-section-deep-dive.png",
			SectionRef: "Tech",
			AltText:    "Deep dive into technology topics",
			Title:      "Tech Section Header",
			Purpose:    "header",
		},
		{
			FileName:   "exploring-depths-philosophy.png",
			SectionRef: "Philosophy",
			AltText:    "Exploring philosophical depths",
			Title:      "Philosophy Section Header",
			Purpose:    "header",
		},
		{
			FileName:   "tech-deep-dive-exploration.png",
			ContentRef: "Deep Dive into the Tech Section",
			AltText:    "Exploring layered technological depths and system architecture",
			Title:      "Tech Deep Dive",
			Purpose:    "header",
		},
		{
			FileName:   "philosophy-depths-exploration.png",
			ContentRef: "Exploring the Depths of Philosophy",
			AltText:    "Excavating philosophical strata and conceptual foundations",
			Title:      "Philosophy Depths",
			Purpose:    "header",
		},
		{
			FileName:   "journey-begins-here.png",
			ContentRef: "Welcome to the Root Section",
			AltText:    "Beginning of exploration across technology, food, and philosophy",
			Title:      "Journey Begins Here",
			Purpose:    "header",
		},
		{
			FileName:   "the-elegance-of-go.png",
			ContentRef: "The Elegance of Go: Simplicity and Concurrency",
			AltText:    "Go programming elegance and concurrency",
			Title:      "The Elegance of Go",
			Purpose:    "header",
		},
		{
			FileName:   "generics-in-go.png",
			ContentRef: "Unlocking Flexibility: Generics in Go 1.18+",
			AltText:    "Generic programming in Go",
			Title:      "Generics in Go 1.18+",
			Purpose:    "header",
		},
		{
			FileName:   "nats-messaging-system.png",
			ContentRef: "NATS: The Messaging System for Cloud-Native",
			AltText:    "NATS messaging system",
			Title:      "NATS Messaging",
			Purpose:    "header",
		},
		{
			FileName:   "hashicorp-nomad-orchestration.png",
			ContentRef: "HashiCorp Nomad: Orchestrating Workloads with Simplicity",
			AltText:    "HashiCorp Nomad orchestration",
			Title:      "HashiCorp Nomad",
			Purpose:    "header",
		},
		{
			FileName:   "kubernetes-101-foundation.png",
			ContentRef: "Kubernetes 101: A Foundation for Container Orchestration",
			AltText:    "Kubernetes fundamentals",
			Title:      "Kubernetes 101",
			Purpose:    "header",
		},
		{
			FileName:   "go-modules-dependency-management.png",
			ContentRef: "Mastering Go Modules: Dependency Management",
			AltText:    "Go modules and dependencies",
			Title:      "Go Modules",
			Purpose:    "header",
		},
		{
			FileName:   "python-data-science.png",
			ContentRef: "Python for Data Science: A Comprehensive Guide",
			AltText:    "Python for data science",
			Title:      "Python Data Science",
			Purpose:    "header",
		},
		{
			FileName:   "microservices-architecture-principles.png",
			ContentRef: "Microservices Architecture: Design Principles",
			AltText:    "Microservices architecture",
			Title:      "Microservices Principles",
			Purpose:    "header",
		},
		{
			FileName:   "go-concurrency-patterns.png",
			ContentRef: "Go Concurrency Patterns: Goroutines and Channels",
			AltText:    "Go concurrency patterns",
			Title:      "Go Concurrency",
			Purpose:    "header",
		},
		{
			FileName:   "dockerizing-microservices.png",
			ContentRef: "Dockerizing Your Microservices: A Practical Guide",
			AltText:    "Docker and microservices",
			Title:      "Dockerizing Microservices",
			Purpose:    "header",
		},
		{
			FileName:   "serverless-architectures.png",
			ContentRef: "Serverless Architectures: Beyond Functions as a Service",
			AltText:    "Serverless architecture concepts",
			Title:      "Serverless Architectures",
			Purpose:    "header",
		},
		{
			FileName:   "idealism-vs-realism.png",
			ContentRef: "Idealism vs. Realism: Two Ways of Seeing the World",
			AltText:    "Idealism versus realism philosophical perspectives",
			Title:      "Idealism vs Realism",
			Purpose:    "header",
		},
		{
			FileName:   "logic-everyday-arguments.png",
			ContentRef: "The Logic of Everyday Arguments: A Practical Guide",
			AltText:    "Logic and argumentation",
			Title:      "Logic of Arguments",
			Purpose:    "header",
		},
		{
			FileName:   "existentialism-freedom-meaning.png",
			ContentRef: "Existentialism: Freedom, Responsibility, and Meaning",
			AltText:    "Existentialist philosophy concepts",
			Title:      "Existentialism",
			Purpose:    "header",
		},
		{
			FileName:   "ethics-of-ai-moral-landscape.png",
			ContentRef: "The Ethics of AI: Navigating the Moral Landscape",
			AltText:    "AI ethics and morality",
			Title:      "Ethics of AI",
			Purpose:    "header",
		},
		{
			FileName:   "stoicism-modern-life.png",
			ContentRef: "Stoicism for Modern Life: Ancient Wisdom, Contemporary Challenges",
			AltText:    "Stoic philosophy for modern living",
			Title:      "Stoicism for Modern Life",
			Purpose:    "header",
		},
		{
			FileName:   "future-of-work-remote-office.png",
			ContentRef: "The Future of Work: Remote vs. Office",
			AltText:    "Remote and office work paradigms",
			Title:      "Future of Work",
			Purpose:    "header",
		},
		{
			FileName:   "impact-of-ai-everyday-life.png",
			ContentRef: "The Impact of AI on Everyday Life",
			AltText:    "AI impact on daily life",
			Title:      "AI Everyday Impact",
			Purpose:    "header",
		},
		{
			FileName:   "understanding-modern-philosophy.png",
			ContentRef: "Understanding Modern Philosophy: A Brief Guide",
			AltText:    "Modern philosophy overview",
			Title:      "Modern Philosophy",
			Purpose:    "header",
		},
		{
			FileName:   "journey-into-stoicism.png",
			ContentRef: "My Journey into Stoicism: Finding Inner Peace",
			AltText:    "Personal journey into Stoic practice",
			Title:      "Journey into Stoicism",
			Purpose:    "header",
		},
		{
			FileName:   "art-of-minimalist-living.png",
			ContentRef: "The Art of Minimalist Living: Less is More",
			AltText:    "Minimalist lifestyle principles",
			Title:      "Minimalist Living",
			Purpose:    "header",
		},
		{
			FileName:   "exploring-ethics-artificial-intelligence.png",
			ContentRef: "Exploring the Ethics of Artificial Intelligence",
			AltText:    "Ethical considerations in AI",
			Title:      "Ethics of AI",
			Purpose:    "header",
		},
	}

	// Get site slug from environment or default to structured
	siteSlug := os.Getenv("SITE_SLUG")
	if siteSlug == "" {
		siteSlug = "structured"
	}

	// Use unified database path
	dbPath := os.Getenv("DB_FILE")
	if dbPath == "" {
		dbPath = "_workspace/db/clio.db"
	}

	// Source images directory (where seed images are stored)
	sourceImagesDir := os.Getenv("SOURCE_IMAGES_DIR")
	if sourceImagesDir == "" {
		sourceImagesDir = "/home/adrian/html/static/images"
	}

	// Destination directory (workspace assets)
	destImagesBase := fmt.Sprintf("_workspace/sites/%s/documents/assets/images", siteSlug)

	log.Printf("Seeding images for site: %s", siteSlug)
	log.Printf("Database path: %s", dbPath)
	log.Printf("Source images: %s", sourceImagesDir)
	log.Printf("Destination base: %s", destImagesBase)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Get site ID from slug
	var siteID string
	err = db.QueryRow("SELECT id FROM site WHERE slug = ? LIMIT 1", siteSlug).Scan(&siteID)
	if err != nil {
		log.Fatalf("Error getting site ID for slug '%s': %v", siteSlug, err)
	}
	log.Printf("Site ID: %s", siteID)

	// Get superadmin user ID for created_by/updated_by
	var superadminID string
	err = db.QueryRow("SELECT id FROM user WHERE username = 'superadmin' LIMIT 1").Scan(&superadminID)
	if err != nil {
		log.Fatalf("Error getting superadmin user: %v", err)
	}

	log.Println("Starting image import...")

	for _, mapping := range mappings {
		var targetID string
		var targetType string // "content" or "section"

		if mapping.SectionRef != "" {
			// This is a section image
			targetType = "section"
			log.Printf("Processing: %s -> Section '%s'", mapping.FileName, mapping.SectionRef)

			err := db.QueryRow(`
				SELECT id FROM section
				WHERE name = ? AND site_id = ?
				LIMIT 1
			`, mapping.SectionRef, siteID).Scan(&targetID)

			if err != nil {
				log.Printf("  Section not found for ref '%s', skipping...", mapping.SectionRef)
				continue
			}
		} else {
			// This is a content image
			targetType = "content"
			log.Printf("Processing: %s -> Content '%s'", mapping.FileName, mapping.ContentRef)

			err := db.QueryRow(`
				SELECT id FROM content
				WHERE heading = ? AND site_id = ?
				LIMIT 1
			`, mapping.ContentRef, siteID).Scan(&targetID)

			if err != nil {
				log.Printf("  Content not found for ref '%s', skipping...", mapping.ContentRef)
				continue
			}
		}

		// Get section path and content slug for proper directory structure
		var sectionPath, contentShortID, contentHeading string
		if targetType == "section" {
			err := db.QueryRow("SELECT path FROM section WHERE id = ?", targetID).Scan(&sectionPath)
			if err != nil {
				log.Printf("  ⚠ Error getting section path: %v", err)
				continue
			}
		} else {
			err := db.QueryRow(`
				SELECT COALESCE(s.path, ''), c.short_id, c.heading
				FROM content c
				LEFT JOIN section s ON c.section_id = s.id
				WHERE c.id = ?
			`, targetID).Scan(&sectionPath, &contentShortID, &contentHeading)
			if err != nil {
				log.Printf("  ⚠ Error getting content info: %v", err)
				continue
			}
		}

		// Generate content slug using same normalization as NormalizeSlug
		var contentSlug string
		if targetType == "content" {
			// Apply same normalization as internal/feat/ssg/slug.go:NormalizeSlug
			s := strings.ToLower(contentHeading)
			s = strings.ReplaceAll(s, " ", "-")

			// Remove non-alphanumeric except hyphens
			var b strings.Builder
			for _, r := range s {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
					b.WriteRune(r)
				}
			}
			s = b.String()

			// Collapse consecutive hyphens
			for strings.Contains(s, "--") {
				s = strings.ReplaceAll(s, "--", "-")
			}

			// Trim leading/trailing hyphens
			s = strings.Trim(s, "-")

			contentSlug = s + "-" + contentShortID
		}

		// Determine destination directory following uploader convention
		var destDir string
		var relativeFilePath string
		assetsBase := destImagesBase

		// Clean section path - remove leading slash for relative paths
		cleanSectionPath := strings.TrimPrefix(sectionPath, "/")

		if targetType == "section" {
			if sectionPath == "/" || sectionPath == "" {
				destDir = assetsBase
				relativeFilePath = mapping.FileName
			} else {
				destDir = filepath.Join(assetsBase, cleanSectionPath)
				relativeFilePath = filepath.Join(cleanSectionPath, mapping.FileName)
			}
		} else {
			if sectionPath == "/" || sectionPath == "" {
				destDir = filepath.Join(assetsBase, contentSlug)
				relativeFilePath = filepath.Join(contentSlug, mapping.FileName)
			} else {
				destDir = filepath.Join(assetsBase, cleanSectionPath, contentSlug)
				relativeFilePath = filepath.Join(cleanSectionPath, contentSlug, mapping.FileName)
			}
		}

		// Create destination directory
		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Printf("  ⚠ Error creating directory %s: %v", destDir, err)
			continue
		}

		// Copy image from seed directory to proper location
		srcPath := filepath.Join(sourceImagesDir, mapping.FileName)
		dstPath := filepath.Join(destDir, mapping.FileName)

		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Printf("  ⚠ Error opening source image %s: %v", srcPath, err)
			continue
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			srcFile.Close()
			log.Printf("  ⚠ Error creating destination file %s: %v", dstPath, err)
			continue
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			log.Printf("  ⚠ Error copying file: %v", err)
			continue
		}

		// Get image dimensions
		imgFile, _ := os.Open(dstPath)
		imgConfig, _, err := image.DecodeConfig(imgFile)
		imgFile.Close()
		if err != nil {
			log.Printf("  ⚠ Error decoding image: %v", err)
			continue
		}

		fileInfo, _ := os.Stat(dstPath)
		fileSize := fileInfo.Size()

		// Create image record with relative path
		imageID := uuid.New().String()
		shortID := uuid.New().String()[:8]
		now := time.Now().UTC().Format("2006-01-02 15:04:05")

		_, err = db.Exec(`
			INSERT INTO image (
				id, site_id, short_id, file_name, file_path,
				alt_text, title, width, height,
				created_by, updated_by, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, imageID, siteID, shortID, mapping.FileName, relativeFilePath,
			mapping.AltText, mapping.Title, imgConfig.Width, imgConfig.Height,
			superadminID, superadminID, now, now)

		if err != nil {
			log.Printf("  ⚠ Error inserting image: %v", err)
			continue
		}

		// Create relationship with content or section
		relationID := uuid.New().String()
		isHeader := 0
		if mapping.Purpose == "header" {
			isHeader = 1
		}

		if targetType == "section" {
			_, err = db.Exec(`
				INSERT INTO section_images (
					id, section_id, image_id, is_header, is_featured, order_num, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?)
			`, relationID, targetID, imageID, isHeader, 0, 0, now)

			if err != nil {
				log.Printf("  Error creating section-image relationship: %v", err)
				continue
			}
		} else {
			_, err = db.Exec(`
				INSERT INTO content_images (
					id, content_id, image_id, is_header, is_featured, order_num, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?)
			`, relationID, targetID, imageID, isHeader, 0, 0, now)

			if err != nil {
				log.Printf("  Error creating content-image relationship: %v", err)
				continue
			}
		}

		log.Printf("  ✓ Imported: %s (%dx%d, %.1f KB)",
			mapping.FileName, imgConfig.Width, imgConfig.Height, float64(fileSize)/1024)
	}

	log.Println("\n✓ Image import completed!")
}
