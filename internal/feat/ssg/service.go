package ssg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

// Service defines the interface for the ssg service.
type Service interface {
	CreateContent(ctx context.Context, content *Content) error
	GetAllContentWithMeta(ctx context.Context) ([]Content, error)
	GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]Content, int, error)
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, content *Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error

	CreateSection(ctx context.Context, section Section) error
	GetSection(ctx context.Context, id uuid.UUID) (Section, error)
	GetSections(ctx context.Context) ([]Section, error)
	UpdateSection(ctx context.Context, section Section) error
	DeleteSection(ctx context.Context, id uuid.UUID) error

	CreateLayout(ctx context.Context, layout Layout) error
	GetLayout(ctx context.Context, id uuid.UUID) (Layout, error)
	GetAllLayouts(ctx context.Context) ([]Layout, error)
	UpdateLayout(ctx context.Context, layout Layout) error
	DeleteLayout(ctx context.Context, id uuid.UUID) error

	CreateTag(ctx context.Context, tag Tag) error
	GetTag(ctx context.Context, id uuid.UUID) (Tag, error)
	GetTagByName(ctx context.Context, name string) (Tag, error)
	GetAllTags(ctx context.Context) ([]Tag, error)
	UpdateTag(ctx context.Context, tag Tag) error
	DeleteTag(ctx context.Context, id uuid.UUID) error

	CreateParam(ctx context.Context, param *Param) error
	GetParam(ctx context.Context, id uuid.UUID) (Param, error)
	GetParamByName(ctx context.Context, name string) (Param, error)
	GetParamByRefKey(ctx context.Context, refKey string) (Param, error)
	ListParams(ctx context.Context) ([]Param, error)
	UpdateParam(ctx context.Context, param *Param) error
	DeleteParam(ctx context.Context, id uuid.UUID) error

	// Image related
	CreateImage(ctx context.Context, image *Image) error
	GetImage(ctx context.Context, id uuid.UUID) (Image, error)
	GetImageByShortID(ctx context.Context, shortID string) (Image, error)
	ListImages(ctx context.Context) ([]Image, error)
	UpdateImage(ctx context.Context, image *Image) error
	DeleteImage(ctx context.Context, id uuid.UUID) error

	// ImageVariant related
	CreateImageVariant(ctx context.Context, variant *ImageVariant) error
	GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error)
	ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error)
	UpdateImageVariant(ctx context.Context, variant *ImageVariant) error
	DeleteImageVariant(ctx context.Context, id uuid.UUID) error

	// Content Image Management
	UploadContentImage(ctx context.Context, contentID uuid.UUID, file multipart.File, header *multipart.FileHeader, imageType ImageType, altText, caption string) (*ImageProcessResult, error)
	GetContentImages(ctx context.Context, contentID uuid.UUID) ([]Image, error)
	DeleteContentImage(ctx context.Context, contentID uuid.UUID, imagePath string) error

	// Section Image Management
	UploadSectionImage(ctx context.Context, sectionID uuid.UUID, file multipart.File, header *multipart.FileHeader, imageType ImageType, altText, caption string) (*ImageProcessResult, error)
	DeleteSectionImage(ctx context.Context, sectionID uuid.UUID, imageType ImageType) error

	// ContentTag related
	AddTagToContent(ctx context.Context, contentID uuid.UUID, tagName string) error
	RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error
	GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error)
	GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error)

	GenerateMarkdown(ctx context.Context) error
	GenerateHTMLFromContent(ctx context.Context) error
	Publish(ctx context.Context, commitMessage string) (string, error)
	Plan(ctx context.Context) (PlanReport, error)
}

// BaseService is the concrete implementation of the Service interface.
type BaseService struct {
	*hm.Service
	assetsFS embed.FS
	repo     Repo
	gen      *Generator
	pub      Publisher
	pm       *ParamManager
	im       *ImageManager
}

func NewService(assetsFS embed.FS, repo Repo, gen *Generator, publisher Publisher, pm *ParamManager, im *ImageManager, params hm.XParams) *BaseService {
	return &BaseService{
		Service:  hm.NewService("ssg-svc", params),
		assetsFS: assetsFS,
		repo:     repo,
		gen:      gen,
		pub:      publisher,
		pm:       pm,
		im:       im,
	}
}

// Publish delegates the publishing task to the underlying pub.
func (svc *BaseService) Publish(ctx context.Context, commitMessage string) (string, error) {
	svc.Log().Info("Service starting publish process")

	// For now, we build the config from the application's configuration.
	cfg := PublisherConfig{
		RepoURL: svc.pm.Get(ctx, SSGKey.PublishRepoURL, ""),
		Branch:  svc.pm.Get(ctx, SSGKey.PublishBranch, ""),
		Auth: hm.GitAuth{
			// NOTE: This is oversimplified. We need to work out a bit more here.
			Method: hm.AuthToken,
			Token:  svc.pm.Get(ctx, SSGKey.PublishAuthToken, ""),
		},
		CommitAuthor: hm.GitCommit{
			UserName:  svc.pm.Get(ctx, SSGKey.PublishCommitUserName, ""),
			UserEmail: svc.pm.Get(ctx, SSGKey.PublishCommitUserEmail, ""),
			Message:   svc.pm.Get(ctx, SSGKey.PublishCommitMessage, ""),
		},
	}

	// Override commit message if provided in the request body
	if commitMessage != "" {
		cfg.CommitAuthor.Message = commitMessage
	}

	// Get the output directory for HTML files, which is the source for publishing
	sourceDir := svc.Cfg().StrValOrDef(SSGKey.HTMLPath, "_workspace/documents/html")

	commitURL, err := svc.pub.Publish(ctx, cfg, sourceDir)
	if err != nil {
		return "", fmt.Errorf("cannot publish site: %w", err)
	}

	svc.Log().Info("Service publish process finished successfully", "commit_url", commitURL)
	return commitURL, nil
}

// Plan delegates the plan task to the underlying pub.
func (svc *BaseService) Plan(ctx context.Context) (PlanReport, error) {
	svc.Log().Info("Service starting plan process")

	// For now, we build the config from the application's configuration.
	cfg := PublisherConfig{
		RepoURL: svc.pm.Get(ctx, SSGKey.PublishRepoURL, ""),
		Branch:  svc.pm.Get(ctx, SSGKey.PublishBranch, ""),
		Auth: hm.GitAuth{
			// NOTE: This is oversimplified. We need to work out a bit more here.
			Method: hm.AuthToken,
			Token:  svc.pm.Get(ctx, SSGKey.PublishAuthToken, ""),
		},
		CommitAuthor: hm.GitCommit{
			UserName:  svc.pm.Get(ctx, SSGKey.PublishCommitUserName, ""),
			UserEmail: svc.pm.Get(ctx, SSGKey.PublishCommitUserEmail, ""),
			Message:   svc.pm.Get(ctx, SSGKey.PublishCommitMessage, ""),
		},
	}

	// Get the output directory for HTML files, which is the source for planning
	sourceDir := svc.Cfg().StrValOrDef(SSGKey.HTMLPath, "_workspace/documents/html")

	report, err := svc.pub.Plan(ctx, cfg, sourceDir)
	if err != nil {
		return PlanReport{}, fmt.Errorf("cannot plan site: %w", err)
	}

	svc.Log().Info("Service plan process finished successfully", "summary", report.Summary)
	return report, nil
}

// GenerateMarkdown generates markdown files from the content in the database.
func (svc *BaseService) GenerateMarkdown(ctx context.Context) error {
	svc.Log().Info("Service starting markdown generation")

	contents, err := svc.repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return fmt.Errorf("cannot get all content with meta: %w", err)
	}

	if err := svc.gen.Generate(contents); err != nil {
		return fmt.Errorf("cannot generate markdown: %w", err)
	}

	svc.Log().Info("Service markdown generation finished")
	return nil
}

// GenerateHTMLFromContent generates HTML files from the content in the database.
func (svc *BaseService) GenerateHTMLFromContent(ctx context.Context) error {
	svc.Log().Info("Service starting HTML generation")

	contents, err := svc.repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return fmt.Errorf("cannot get all content with meta: %w", err)
	}

	// Set placeholder for content without image
	for range contents {
		// TODO: Handle placeholder image via relationships
		// if contents[i].Image == "" {
		//	contents[i].Image = "/static/img/placeholder.png"
		// }
	}

	sections, err := svc.repo.GetSections(ctx)
	if err != nil {
		return fmt.Errorf("cannot get sections: %w", err)
	}

	var menuSections []Section
	for _, s := range sections {
		if s.Name != "root" {
			menuSections = append(menuSections, s)
		}
	}

	layoutPath := svc.Cfg().StrValOrDef(SSGKey.LayoutPath, "assets/ssg/layout/layout.html")
	tmpl, err := template.ParseFS(svc.assetsFS,
		layoutPath,
		"assets/ssg/partial/list.tmpl",
		"assets/ssg/partial/blocks.tmpl",
		"assets/ssg/partial/article-blocks.tmpl",
		"assets/ssg/partial/blog-blocks.tmpl",
		"assets/ssg/partial/series-blocks.tmpl",
		"assets/ssg/partial/pagination.tmpl",
		"assets/ssg/partial/google-search.tmpl",
	)
	if err != nil {
		return fmt.Errorf("cannot parse template from embedded fs: %w", err)
	}

	htmlPath := svc.Cfg().StrValOrDef(SSGKey.HTMLPath, "_workspace/documents/html")

	if err := CopyStaticAssets(svc.assetsFS, htmlPath); err != nil {
		return fmt.Errorf("cannot copy static assets: %w", err)
	}

	// Copy dynamic images from assets/images to html/static/images
	workspaceDir := svc.Cfg().StrValOrDef(SSGKey.WorkspacePath, "_workspace")
	docsDir := filepath.Join(workspaceDir, "documents") // This gives us _workspace/documents
	svc.Log().Info("Copying dynamic images", "from", filepath.Join(docsDir, "assets", "images"), "to", filepath.Join(htmlPath, "static", "images"))
	if err := CopyDynamicImages(docsDir, htmlPath); err != nil {
		svc.Log().Error("Failed to copy dynamic images", "error", err)
		return fmt.Errorf("cannot copy dynamic images: %w", err)
	}
	svc.Log().Info("Dynamic images copied successfully")

	headerStyle := svc.Cfg().StrValOrDef(SSGKey.HeaderStyle, "boxed", true)
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".webp"}

	// Prepare SearchData
	searchData := SearchData{
		Provider: "google", // O el proveedor que corresponda
		Enabled:  svc.Cfg().BoolVal(SSGKey.SearchGoogleEnabled, false),
		ID:       svc.Cfg().StrValOrDef(SSGKey.SearchGoogleID, ""),
	}
	svc.Log().Info("SearchData values", "enabled", searchData.Enabled, "id", searchData.ID) // Línea de log modificada

	for _, content := range contents {
		svc.Log().Debug("Processing content for HTML generation", "slug", content.Slug(), "section_path", content.SectionPath)
		if content.Draft {
			svc.Log().Debug("Skipping draft content", "slug", content.Slug())
			continue
		}

		headerImagePath := ""

		if content.HeaderImageURL != "" {
			headerImagePath = content.HeaderImageURL
		} else {
			contentDir := filepath.Join(htmlPath, content.SectionPath, content.Slug())
			contentImgDir := filepath.Join(contentDir, "img")

			foundSpecificHeader := false
			for _, ext := range imageExtensions {
				checkPath := filepath.Join("assets", "content", content.SectionPath, content.Slug(), "img", "header"+ext)
				if f, err := svc.assetsFS.Open(checkPath); err == nil {
					f.Close()
					if err := os.MkdirAll(contentImgDir, 0755); err != nil {
						return fmt.Errorf("cannot create img directory: %w", err)
					}
					dst := filepath.Join(contentImgDir, "header"+ext)
					if err := copyFile(svc.assetsFS, checkPath, dst); err != nil {
						return fmt.Errorf("cannot copy specific header: %w", err)
					}
					headerImagePath = "img/header" + ext
					foundSpecificHeader = true
					break
				}
			}

			if !foundSpecificHeader {
				headerImagePath = "/static/img/header.png"
			}
		}

		assetPath := "/"

		contentImages, err := svc.GetContentImages(ctx, content.ID)
		if err != nil {
			svc.Log().Debug("Failed to load content images", "contentID", content.ID, "error", err)
			contentImages = []Image{}
		} else {
			svc.Log().Info("Loaded content images", "contentID", content.ID, "count", len(contentImages), "slug", content.Slug())
		}

		imageContext := &ImageContext{
			Images: make(map[string]ImageMetadata),
		}

		for _, img := range contentImages {
			svc.Log().Debug("Adding image to context", "filePath", img.FilePath, "caption", img.Caption, "altText", img.AltText)
			imageContext.Images[img.FilePath] = ImageMetadata{
				AltText:         img.AltText,
				Caption:         img.Caption,
				LongDescription: img.LongDescription,
				Title:           img.Title,
				Decorative:      img.Decorative,
			}
		}
		svc.Log().Info("Image context created", "imageCount", len(imageContext.Images))

		processor := NewMarkdownProcessor()

		htmlBody, err := processor.ToHTMLWithImageContext([]byte(content.Body), imageContext)
		if err != nil {
			svc.Log().Error("Error converting markdown to HTML", "slug", content.Slug(), "error", err)
			continue
		}

		if headerStyle == "boxed" || headerStyle == "overlay" {
			htmlBody = svc.removeFirstH1(htmlBody)
		}

		pageContent := PageContent{
			Heading:            content.Heading,
			HeaderImage:        headerImagePath,
			HeaderImageAlt:     content.HeaderImageAlt,
			HeaderImageCaption: content.HeaderImageCaption,
			Body:               template.HTML(htmlBody),
			Kind:               content.Kind,
		}

		blocks := BuildBlocks(content, contents, int(svc.Cfg().IntVal(SSGKey.BlocksMaxItems, 5)))

		data := PageData{
			HeaderStyle: headerStyle,
			AssetPath:   assetPath,
			Menu:        menuSections,
			Content:     pageContent,
			Blocks:      blocks,
			Search:      searchData,
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			svc.Log().Error("Error executing template for content", "slug", content.Slug(), "error", err)
			continue
		}

		outputPath := filepath.Join(htmlPath, content.SectionPath, content.Slug(), "index.html")

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			svc.Log().Error("Error creating directory for HTML file", "path", outputPath, "error", err)
			continue
		}

		if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			svc.Log().Error("Error writing index HTML file", "path", outputPath, "error", err)
			continue
		}
	}

	// Generate index pages
	svc.Log().Info("Building site indexes...")
	indexes := BuildIndexes(contents, sections)

	// Create a lookup map for manual index pages
	manualIndexPages := make(map[string]bool)
	for _, c := range contents {
		if strings.ToLower(c.Kind) == "page" && c.Slug() == "index" {
			manualIndexPages[c.SectionPath] = true
		}
	}

	postsPerPage := int(svc.Cfg().IntVal(SSGKey.IndexMaxItems, 9))

	for _, index := range indexes {
		// Check if a manual index page exists for this path
		if manualIndexPages[index.Path] {
			svc.Log().Info(fmt.Sprintf("Skipping index generation for '%s': manual index page found.", index.Path))
			continue
		}

		// Get section header image for this index
		var sectionHeaderImage string
		for _, section := range sections {
			if section.Path == index.Path {
				headerPath, err := svc.GetSectionHeaderImage(ctx, section.ID)
				if err == nil && headerPath != "" {
					sectionHeaderImage = "/static/images/" + headerPath
				}
				break
			}
		}

		// Paginate the content
		totalContent := len(index.Content)
		if totalContent == 0 {
			continue
		}
		totalPages := (totalContent + postsPerPage - 1) / postsPerPage

		for page := 1; page <= totalPages; page++ {
			start := (page - 1) * postsPerPage
			end := start + postsPerPage
			if end > totalContent {
				end = totalContent
			}
			pageContent := index.Content[start:end]

			// Determine output path for the index page
			var outputPath string
			if page == 1 {
				outputPath = filepath.Join(htmlPath, index.Path, "index.html")
			} else {
				outputPath = filepath.Join(htmlPath, index.Path, "page", fmt.Sprintf("%d", page), "index.html")
			}

			assetPath := "/"

			// Prepare pagination data
			pagination := &PaginationData{
				CurrentPage: page,
				TotalPages:  totalPages,
			}
			if page > 1 {
				if page == 2 {
					pagination.PrevPageURL = assetPath + strings.TrimSuffix(index.Path, "/")
				} else {
					pagination.PrevPageURL = fmt.Sprintf("%spage/%d", assetPath, page-1)
				}
			}
			if page < totalPages {
				pagination.NextPageURL = fmt.Sprintf("%spage/%d", assetPath, page+1)
			}

			data := PageData{
				HeaderStyle:        headerStyle,
				AssetPath:          assetPath,
				Menu:               menuSections,
				IsIndex:            true,
				ListPageContent:    pageContent,
				Pagination:         pagination,
				Search:             searchData,
				SectionHeaderImage: sectionHeaderImage,
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				svc.Log().Error("Error executing template for index", "path", index.Path, "error", err)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				svc.Log().Error("Error creating directory for index file", "path", outputPath, "error", err)
				continue
			}

			if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
				svc.Log().Error("Error writing index HTML file", "path", outputPath, "error", err)
				continue
			}
		}
	}

	svc.Log().Info("Service HTML generation finished")
	return nil
}

// Content related

func (svc *BaseService) CreateContent(ctx context.Context, content *Content) error {
	return svc.repo.CreateContent(ctx, content)
}

func (svc *BaseService) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {
	return svc.repo.GetContent(ctx, id)
}

func (svc *BaseService) UpdateContent(ctx context.Context, content *Content) error {
	return svc.repo.UpdateContent(ctx, content)
}

func (svc *BaseService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteContent(ctx, id)
}

func (svc *BaseService) GetAllContentWithMeta(ctx context.Context) ([]Content, error) {
	return svc.repo.GetAllContentWithMeta(ctx)
}

func (svc *BaseService) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]Content, int, error) {
	return svc.repo.GetContentWithPaginationAndSearch(ctx, offset, limit, searchQuery)
}

// Section related
func (svc *BaseService) CreateSection(ctx context.Context, section Section) error {
	return svc.repo.CreateSection(ctx, section)
}

func (svc *BaseService) GetSection(ctx context.Context, id uuid.UUID) (Section, error) {
	return svc.repo.GetSection(ctx, id)
}

func (svc *BaseService) GetSections(ctx context.Context) ([]Section, error) {
	return svc.repo.GetSections(ctx)
}

func (svc *BaseService) UpdateSection(ctx context.Context, section Section) error {
	return svc.repo.UpdateSection(ctx, section)
}

func (svc *BaseService) DeleteSection(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteSection(ctx, id)
}

// Layout related
func (svc *BaseService) CreateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.CreateLayout(ctx, layout)
}

func (svc *BaseService) GetLayout(ctx context.Context, id uuid.UUID) (Layout, error) {
	return svc.repo.GetLayout(ctx, id)
}

func (svc *BaseService) GetAllLayouts(ctx context.Context) ([]Layout, error) {
	return svc.repo.GetAllLayouts(ctx)
}

func (svc *BaseService) UpdateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.UpdateLayout(ctx, layout)
}

func (svc *BaseService) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteLayout(ctx, id)
}

// Tag related
func (svc *BaseService) CreateTag(ctx context.Context, tag Tag) error {
	return svc.repo.CreateTag(ctx, tag)
}

func (svc *BaseService) GetTag(ctx context.Context, id uuid.UUID) (Tag, error) {
	return svc.repo.GetTag(ctx, id)
}

func (svc *BaseService) GetTagByName(ctx context.Context, name string) (Tag, error) {
	return svc.repo.GetTagByName(ctx, name)
}

func (svc *BaseService) GetAllTags(ctx context.Context) ([]Tag, error) {
	return svc.repo.GetAllTags(ctx)
}

func (svc *BaseService) UpdateTag(ctx context.Context, tag Tag) error {
	return svc.repo.UpdateTag(ctx, tag)
}

func (svc *BaseService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteTag(ctx, id)
}

// Param related
func (svc *BaseService) CreateParam(ctx context.Context, param *Param) error {
	return svc.repo.CreateParam(ctx, param)
}

func (svc *BaseService) GetParam(ctx context.Context, id uuid.UUID) (Param, error) {
	return svc.repo.GetParam(ctx, id)
}

func (svc *BaseService) GetParamByName(ctx context.Context, name string) (Param, error) {
	return svc.repo.GetParamByName(ctx, name)
}

func (svc *BaseService) GetParamByRefKey(ctx context.Context, refKey string) (Param, error) {
	return svc.repo.GetParamByRefKey(ctx, refKey)
}

func (svc *BaseService) ListParams(ctx context.Context) ([]Param, error) {
	return svc.repo.ListParams(ctx)
}

func (svc *BaseService) UpdateParam(ctx context.Context, param *Param) error {
	return svc.repo.UpdateParam(ctx, param)
}

func (svc *BaseService) DeleteParam(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteParam(ctx, id)
}

// Image related
func (svc *BaseService) CreateImage(ctx context.Context, image *Image) error {
	return svc.repo.CreateImage(ctx, image)
}

func (svc *BaseService) GetImage(ctx context.Context, id uuid.UUID) (Image, error) {
	return svc.repo.GetImage(ctx, id)
}

func (svc *BaseService) GetImageByShortID(ctx context.Context, shortID string) (Image, error) {
	return svc.repo.GetImageByShortID(ctx, shortID)
}

func (svc *BaseService) ListImages(ctx context.Context) ([]Image, error) {
	return svc.repo.ListImages(ctx)
}

func (svc *BaseService) UpdateImage(ctx context.Context, image *Image) error {
	return svc.repo.UpdateImage(ctx, image)
}

func (svc *BaseService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteImage(ctx, id)
}

// ImageVariant related
func (svc *BaseService) CreateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return svc.repo.CreateImageVariant(ctx, variant)
}

func (svc *BaseService) GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error) {
	return svc.repo.GetImageVariant(ctx, id)
}

func (svc *BaseService) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error) {
	return svc.repo.ListImageVariantsByImageID(ctx, imageID)
}

func (svc *BaseService) UpdateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return svc.repo.UpdateImageVariant(ctx, variant)
}

func (svc *BaseService) DeleteImageVariant(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteImageVariant(ctx, id)
}

// ContentTag related
func (svc *BaseService) AddTagToContent(ctx context.Context, contentID uuid.UUID, tagName string) error {
	tag, err := svc.repo.GetTagByName(ctx, tagName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error getting tag by name: %w", err)
	}

	if tag.IsZero() {
		newTag := NewTag(tagName)
		newTag.GenCreateValues()
		err = svc.repo.CreateTag(ctx, newTag)
		if err != nil {
			return fmt.Errorf("error creating tag: %w", err)
		}
		tag = newTag
	}

	return svc.repo.AddTagToContent(ctx, contentID, tag.ID)
}

func (svc *BaseService) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	return svc.repo.RemoveTagFromContent(ctx, contentID, tagID)
}

func (svc *BaseService) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error) {
	return svc.repo.GetTagsForContent(ctx, contentID)
}

func (svc *BaseService) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error) {
	return svc.repo.GetContentForTag(ctx, tagID)
}

var firstH1Regex = regexp.MustCompile(`(?i)<h1[^>]*>.*?</h1>`)

// removeFirstH1 removes the first <h1>...</h1> tag from an HTML string.
func (svc *BaseService) removeFirstH1(htmlContent string) string {
	return firstH1Regex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		// Only replace the first occurrence
		if strings.HasPrefix(htmlContent, match) {
			return ""
		}
		return match
	})
}

// Content Image Management

// UploadContentImage handles uploading images for content (header or content images)
func (svc *BaseService) UploadContentImage(ctx context.Context, contentID uuid.UUID, file multipart.File, header *multipart.FileHeader, imageType ImageType, altText, caption string) (*ImageProcessResult, error) {
	svc.Log().Debugf("Uploading content image: contentID=%s, type=%s", contentID, imageType)

	content, err := svc.repo.GetContent(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	var section *Section
	if content.SectionID != uuid.Nil {
		s, err := svc.repo.GetSection(ctx, content.SectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get section: %w", err)
		}
		section = &s
	}

	result, err := svc.im.ProcessUpload(ctx, file, header, &content, section, imageType, altText, caption)
	if err != nil {
		return nil, fmt.Errorf("failed to process upload: %w", err)
	}

	contentHash, err := calculateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Create Image record with accessibility metadata - always create new record
	image := Image{
		Title:           result.Filename,
		FilePath:        result.RelativePath,
		ContentHash:     contentHash,
		AltText:         altText,
		Caption:         caption,
		LongDescription: caption, // Use caption as long description for now
	}
	image.GenCreateValues()

	if err := svc.repo.CreateImage(ctx, &image); err != nil {
		svc.im.DeleteImage(ctx, result.RelativePath)
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}

	contentImage := NewContentImage(contentID, image.GetID(), string(imageType))

	if err := svc.repo.CreateContentImage(ctx, contentImage); err != nil {
		svc.im.DeleteImage(ctx, result.RelativePath)
		svc.repo.DeleteImage(ctx, image.GetID())
		return nil, fmt.Errorf("failed to create content-image relationship: %w", err)
	}

	// TODO: Remove direct field update when we complete migration
	// if imageType == ImageTypeHeader {
	//	content.Image = result.RelativePath
	//	if err := svc.repo.UpdateContent(ctx, &content); err != nil {
	//		return nil, fmt.Errorf("failed to update content with header image: %w", err)
	//	}
	// }

	return result, nil
}

// GetContentImages returns all images for a specific content via relationships
func (svc *BaseService) GetContentImages(ctx context.Context, contentID uuid.UUID) ([]Image, error) {
	svc.Log().Debugf("Getting content images: contentID=%s", contentID)

	contentImages, err := svc.repo.GetContentImagesByContentID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content images: %w", err)
	}

	var images []Image
	for _, ci := range contentImages {
		if !ci.IsActive {
			continue
		}

		image, err := svc.repo.GetImage(ctx, ci.ImageID)
		if err != nil {
			svc.Log().Info("Failed to get image %s: %v", ci.ImageID, err)
			continue
		}

		image.Purpose = ci.Purpose
		images = append(images, image)
	}

	return images, nil
}

// GetContentHeaderImage returns the header image for a specific content
func (svc *BaseService) GetContentHeaderImage(ctx context.Context, contentID uuid.UUID) (string, error) {
	contentImages, err := svc.repo.GetContentImagesByContentID(ctx, contentID)
	if err != nil {
		return "", fmt.Errorf("failed to get content images: %w", err)
	}

	for _, ci := range contentImages {
		if ci.Purpose == "header" && ci.IsActive {
			image, err := svc.repo.GetImage(ctx, ci.ImageID)
			if err != nil {
				svc.Log().Info("Failed to get header image %s: %v", ci.ImageID, err)
				continue
			}

			return image.FilePath, nil
		}
	}

	return "", nil // No header image found
}

// GetSectionHeaderImage returns the header image for a specific section
func (svc *BaseService) GetSectionHeaderImage(ctx context.Context, sectionID uuid.UUID) (string, error) {
	sectionImages, err := svc.repo.GetSectionImagesBySectionID(ctx, sectionID)
	if err != nil {
		return "", fmt.Errorf("failed to get layout images: %w", err)
	}

	for _, si := range sectionImages {
		if si.Purpose == "header" && si.IsActive {
			image, err := svc.repo.GetImage(ctx, si.ImageID)
			if err != nil {
				svc.Log().Info("Failed to get section header image %s: %v", si.ImageID, err)
				continue
			}

			return image.FilePath, nil
		}
	}

	return "", nil // No header image found
}

// GetSectionBlogHeaderImage returns the blog header image for a specific section
func (svc *BaseService) GetSectionBlogHeaderImage(ctx context.Context, sectionID uuid.UUID) (string, error) {
	sectionImages, err := svc.repo.GetSectionImagesBySectionID(ctx, sectionID)
	if err != nil {
		return "", fmt.Errorf("failed to get layout images: %w", err)
	}

	for _, si := range sectionImages {
		if si.Purpose == "blog_header" && si.IsActive {
			image, err := svc.repo.GetImage(ctx, si.ImageID)
			if err != nil {
				svc.Log().Info("Failed to get section blog header image %s: %v", si.ImageID, err)
				continue
			}

			return image.FilePath, nil
		}
	}

	return "", nil // No blog header image found
}

// DeleteContentImage deletes a content image by path
func (svc *BaseService) DeleteContentImage(ctx context.Context, contentID uuid.UUID, imagePath string) error {
	svc.Log().Infof("Deleting content image: contentID=%s, imagePath=%s", contentID, imagePath)

	_, err := svc.repo.GetContent(ctx, contentID)
	if err != nil {
		svc.Log().Errorf("Failed to get content: %v", err)
		return fmt.Errorf("failed to get content: %w", err)
	}

	contentImages, err := svc.repo.GetContentImagesByContentID(ctx, contentID)
	if err != nil {
		return fmt.Errorf("failed to get content images: %w", err)
	}

	var imageToDelete *Image
	var relationshipToDelete *ContentImage
	for _, ci := range contentImages {
		image, err := svc.repo.GetImage(ctx, ci.ImageID)
		if err != nil {
			svc.Log().Info("Failed to get image %s: %v", ci.ImageID, err)
			continue
		}
		if image.FilePath == imagePath {
			imageToDelete = &image
			relationshipToDelete = &ci
			break
		}
	}

	if imageToDelete == nil {
		svc.Log().Info("Image not found in database for path: %s", imagePath)
		if err := svc.im.DeleteImage(ctx, imagePath); err != nil {
			return fmt.Errorf("failed to delete image file: %w", err)
		}
		return nil
	}

	if err := svc.repo.DeleteContentImage(ctx, relationshipToDelete.ID); err != nil {
		return fmt.Errorf("failed to delete content image relationship: %w", err)
	}

	if err := svc.repo.DeleteImage(ctx, imageToDelete.ID); err != nil {
		return fmt.Errorf("failed to delete image record: %w", err)
	}

	if err := svc.im.DeleteImage(ctx, imagePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

// Section Image Management

// UploadSectionImage handles uploading images for sections (section header or blog header)
func (svc *BaseService) UploadSectionImage(ctx context.Context, sectionID uuid.UUID, file multipart.File, header *multipart.FileHeader, imageType ImageType, altText, caption string) (*ImageProcessResult, error) {
	section, err := svc.repo.GetSection(ctx, sectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get section: %w", err)
	}

	result, err := svc.im.ProcessUpload(ctx, file, header, nil, &section, imageType, altText, caption)
	if err != nil {
		return nil, fmt.Errorf("failed to process upload: %w", err)
	}

	contentHash, err := calculateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Create Image record with accessibility metadata - always create new record
	image := Image{
		Title:           result.Filename,
		FilePath:        result.RelativePath,
		ContentHash:     contentHash,
		AltText:         altText,
		Caption:         caption,
		LongDescription: caption, // Use caption as long description for now
	}
	image.GenCreateValues()

	if err := svc.repo.CreateImage(ctx, &image); err != nil {
		svc.im.DeleteImage(ctx, result.RelativePath)
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}

	purposeStr := string(imageType)
	if imageType == ImageTypeSectionHeader {
		purposeStr = "header"
	} else if imageType == ImageTypeBlogHeader {
		purposeStr = "blog_header"
	}
	sectionImage := NewSectionImage(sectionID, image.GetID(), purposeStr)

	if err := svc.repo.CreateSectionImage(ctx, sectionImage); err != nil {
		svc.im.DeleteImage(ctx, result.RelativePath)
		svc.repo.DeleteImage(ctx, image.GetID())
		return nil, fmt.Errorf("failed to create section-image relationship: %w", err)
	}

	return result, nil
}

func (svc *BaseService) DeleteSectionImage(ctx context.Context, sectionID uuid.UUID, imageType ImageType) error {
	_, err := svc.repo.GetSection(ctx, sectionID)
	if err != nil {
		return fmt.Errorf("failed to get section: %w", err)
	}

	sectionImages, err := svc.repo.GetSectionImagesBySectionID(ctx, sectionID)
	if err != nil {
		return fmt.Errorf("failed to get layout images: %w", err)
	}

	var imageToDelete *Image
	var relationshipToDelete *SectionImage
	purposeStr := string(imageType)
	if imageType == ImageTypeSectionHeader {
		purposeStr = "header"
	} else if imageType == ImageTypeBlogHeader {
		purposeStr = "blog_header"
	}

	for _, si := range sectionImages {
		if si.Purpose == purposeStr && si.IsActive {
			image, err := svc.repo.GetImage(ctx, si.ImageID)
			if err != nil {
				svc.Log().Info("Failed to get image %s: %v", si.ImageID, err)
				continue
			}
			imageToDelete = &image
			relationshipToDelete = &si
			break
		}
	}

	if imageToDelete == nil {
		svc.Log().Info("No %s image found for section %s", imageType, sectionID)
		return nil // Nothing to delete
	}

	if err := svc.repo.DeleteSectionImage(ctx, relationshipToDelete.ID); err != nil {
		return fmt.Errorf("failed to delete layout image relationship: %w", err)
	}

	if err := svc.repo.DeleteImage(ctx, imageToDelete.ID); err != nil {
		return fmt.Errorf("failed to delete image record: %w", err)
	}

	if err := svc.im.DeleteImage(ctx, imageToDelete.FilePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

// calculateFileHash calculates SHA-256 hash of a multipart file
func calculateFileHash(file multipart.File) (string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer after hashing: %w", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
