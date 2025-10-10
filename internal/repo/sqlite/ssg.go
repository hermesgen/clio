package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

var (
	featSSG         = "ssg"
	resLayout       = "layout"
	resContent      = "content"
	resMeta         = "meta"
	resSection      = "section"
	resTag          = "tag"
	resParam        = "param"
	resImage        = "image"
	resImageVariant = "image_variant"
)

// sanitizeURLPath sanitizes a file path for safe use in URLs
func sanitizeURLPath(path string) string {
	// Split path into directory and filename components
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part != "" { // Don't sanitize empty parts (preserves leading slashes)
			// Replace problematic characters with hyphens
			re := regexp.MustCompile(`[^a-zA-Z0-9\-_.]`)
			sanitized := re.ReplaceAllString(part, "-")

			// Remove multiple consecutive hyphens
			re2 := regexp.MustCompile(`-+`)
			sanitized = re2.ReplaceAllString(sanitized, "-")

			// Remove leading/trailing hyphens
			sanitized = strings.Trim(sanitized, "-")

			parts[i] = strings.ToLower(sanitized)
		}
	}
	return strings.Join(parts, "/")
}

// Content related

func (repo *ClioRepo) CreateContent(ctx context.Context, c *ssg.Content) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	// Create Content
	contentQuery, err := repo.BaseRepo.Query().Get(featSSG, resContent, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create content query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, contentQuery, c); err != nil {
		return fmt.Errorf("cannot create content: %w", err)
	}

	// Create Meta
	c.Meta.ContentID = c.ID
	c.Meta.GenID()
	c.Meta.GenCreateValues(c.CreatedBy)
	metaQuery, err := repo.BaseRepo.Query().Get(featSSG, resMeta, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create meta query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, metaQuery, c.Meta); err != nil {
		return fmt.Errorf("cannot create meta: %w", err)
	}

	return nil
}

func (repo *ClioRepo) GetContent(ctx context.Context, id uuid.UUID) (ssg.Content, error) {
	// This is a placeholder. A specific query "GetWithMeta" is needed for optimal performance.
	// For now, we will filter from the large GetAll query.
	contents, err := repo.GetAllContentWithMeta(ctx)
	if err != nil {
		return ssg.Content{}, err
	}
	for _, content := range contents {
		if content.ID == id {
			return content, nil
		}
	}
	return ssg.Content{}, errors.New("content not found")
}

func (repo *ClioRepo) UpdateContent(ctx context.Context, c *ssg.Content) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	// Update Content
	contentQuery, err := repo.BaseRepo.Query().Get(featSSG, resContent, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update content query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, contentQuery, c); err != nil {
		return fmt.Errorf("cannot update content: %w", err)
	}

	// Update Meta
	metaQuery, err := repo.BaseRepo.Query().Get(featSSG, resMeta, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update meta query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, metaQuery, c.Meta); err != nil {
		return fmt.Errorf("cannot update meta: %w", err)
	}

	return nil
}

func (repo *ClioRepo) DeleteContent(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resContent, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *ClioRepo) GetAllContentWithMeta(ctx context.Context) ([]ssg.Content, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resContent, "GetAllContentWithMeta")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contentMap := make(map[uuid.UUID]*ssg.Content)
	var contentOrder []uuid.UUID

	for rows.Next() {
		var c ssg.Content
		var m ssg.Meta
		var t ssg.Tag
		var sectionPath, sectionName sql.NullString
		var publishedAt sql.NullTime

		var metaID sql.NullString
		var description, keywords, robots, canonicalURL, sitemap sql.NullString
		var tableOfContents, share, comments sql.NullBool

		var tagID, tagShortID, tagName, tagSlug sql.NullString
		var contentImageID, imagePurpose, imageFilePath, imageAltText, imageLongDescription sql.NullString

		err := rows.Scan(
			&c.ID, &c.UserID, &c.SectionID, &c.Kind, &c.Heading, &c.Body, &c.Draft, &c.Featured, &publishedAt, &c.ShortID,
			&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
			&sectionPath, &sectionName,
			&metaID, &description, &keywords, &robots, &canonicalURL, &sitemap, &tableOfContents, &share, &comments,
			&tagID, &tagShortID, &tagName, &tagSlug,
			&contentImageID, &imagePurpose, &imageFilePath, &imageAltText, &imageLongDescription,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if _, ok := contentMap[c.ID]; !ok {
			c.SectionPath = sectionPath.String
			c.SectionName = sectionName.String
			if publishedAt.Valid {
				c.PublishedAt = &publishedAt.Time
			}

			if metaID.Valid {
				m.ID, _ = uuid.Parse(metaID.String)
				m.ContentID = c.ID
				m.Description = description.String
				m.Keywords = keywords.String
				m.Robots = robots.String
				m.CanonicalURL = canonicalURL.String
				m.Sitemap = sitemap.String
				m.TableOfContents = tableOfContents.Bool
				m.Share = share.Bool
				m.Comments = comments.Bool
				c.Meta = m
			}

			contentMap[c.ID] = &c
			contentOrder = append(contentOrder, c.ID)
		}

		// Handle images
		if contentImageID.Valid && imageFilePath.Valid {
			// Sanitize the file path for URL safety
			sanitizedPath := sanitizeURLPath(imageFilePath.String)
			// Convert file_path to URL (add /static/images prefix)
			imageURL := "/static/images" + sanitizedPath
			if imagePurpose.String == "thumbnail" {
				contentMap[c.ID].ThumbnailURL = imageURL
			} else if imagePurpose.String == "header" {
				contentMap[c.ID].HeaderImageURL = imageURL
				contentMap[c.ID].HeaderImageAlt = imageAltText.String
				contentMap[c.ID].HeaderImageCaption = imageLongDescription.String
			} else if imagePurpose.String == "content" {
				// Use content images as fallback for thumbnails if no thumbnail exists
				if contentMap[c.ID].ThumbnailURL == "" {
					contentMap[c.ID].ThumbnailURL = imageURL
				}
			}
		}

		if tagID.Valid {
			t.ID, _ = uuid.Parse(tagID.String)
			t.SetShortID(tagShortID.String)
			t.Name = tagName.String
			t.SlugField = tagSlug.String
			contentMap[c.ID].Tags = append(contentMap[c.ID].Tags, t)
		}
	}

	contents := make([]ssg.Content, len(contentOrder))
	for i, id := range contentOrder {
		contents[i] = *contentMap[id]
	}

	return contents, nil
}

func (repo *ClioRepo) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]ssg.Content, int, error) {
	countQuery, err := repo.BaseRepo.Query().Get(featSSG, resContent, "GetContentCountWithSearch")
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get count query: %w", err)
	}

	var totalCount int
	row := repo.db.QueryRowxContext(ctx, countQuery, searchQuery, searchQuery)
	err = row.Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get total count: %w", err)
	}

	query, err := repo.BaseRepo.Query().Get(featSSG, resContent, "GetContentWithPaginationAndSearch")
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get pagination query: %w", err)
	}

	rows, err := repo.db.QueryxContext(ctx, query, searchQuery, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot execute pagination query: %w", err)
	}
	defer rows.Close()

	contentMap := make(map[uuid.UUID]*ssg.Content)
	var contentOrder []uuid.UUID

	for rows.Next() {
		var c ssg.Content
		var m ssg.Meta
		var t ssg.Tag
		var sectionPath, sectionName sql.NullString
		var publishedAt sql.NullTime

		var metaID sql.NullString
		var description, keywords, robots, canonicalURL, sitemap sql.NullString
		var tableOfContents, share, comments sql.NullBool

		var tagID, tagShortID, tagName, tagSlug sql.NullString
		var contentImageID, imagePurpose, imageFilePath, imageAltText, imageLongDescription sql.NullString

		err := rows.Scan(
			&c.ID, &c.UserID, &c.SectionID, &c.Kind, &c.Heading, &c.Body, &c.Draft, &c.Featured, &publishedAt, &c.ShortID,
			&c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt,
			&sectionPath, &sectionName,
			&metaID, &description, &keywords, &robots, &canonicalURL, &sitemap, &tableOfContents, &share, &comments,
			&tagID, &tagShortID, &tagName, &tagSlug,
			&contentImageID, &imagePurpose, &imageFilePath, &imageAltText, &imageLongDescription,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}

		if _, ok := contentMap[c.ID]; !ok {
			c.SectionPath = sectionPath.String
			c.SectionName = sectionName.String
			if publishedAt.Valid {
				c.PublishedAt = &publishedAt.Time
			}

			if metaID.Valid {
				m.ID, _ = uuid.Parse(metaID.String)
				m.ContentID = c.ID
				m.Description = description.String
				m.Keywords = keywords.String
				m.Robots = robots.String
				m.CanonicalURL = canonicalURL.String
				m.Sitemap = sitemap.String
				m.TableOfContents = tableOfContents.Bool
				m.Share = share.Bool
				m.Comments = comments.Bool
				c.Meta = m
			}

			contentMap[c.ID] = &c
			contentOrder = append(contentOrder, c.ID)
		}

		if contentImageID.Valid && imageFilePath.Valid {
			sanitizedPath := sanitizeURLPath(imageFilePath.String)
			imageURL := "/static/images" + sanitizedPath
			if imagePurpose.String == "thumbnail" {
				contentMap[c.ID].ThumbnailURL = imageURL
			} else if imagePurpose.String == "header" {
				contentMap[c.ID].HeaderImageURL = imageURL
				contentMap[c.ID].HeaderImageAlt = imageAltText.String
				contentMap[c.ID].HeaderImageCaption = imageLongDescription.String
			} else if imagePurpose.String == "content" {
				if contentMap[c.ID].ThumbnailURL == "" {
					contentMap[c.ID].ThumbnailURL = imageURL
				}
			}
		}

		if tagID.Valid {
			t.ID, _ = uuid.Parse(tagID.String)
			t.SetShortID(tagShortID.String)
			t.Name = tagName.String
			t.SlugField = tagSlug.String
			contentMap[c.ID].Tags = append(contentMap[c.ID].Tags, t)
		}
	}

	contents := make([]ssg.Content, len(contentOrder))
	for i, id := range contentOrder {
		contents[i] = *contentMap[id]
	}

	return contents, totalCount, nil
}

// Section related

func (repo *ClioRepo) CreateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resSection, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		section.GetID(),
		section.GetShortID(),
		section.Name,
		section.Description,
		section.Path,
		section.LayoutID,
		section.GetCreatedBy(),
		section.GetUpdatedBy(),
		section.GetCreatedAt(),
		section.GetUpdatedAt(),
	)
	return err
}

func (repo *ClioRepo) GetSections(ctx context.Context) ([]ssg.Section, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resSection, "GetAll")
	if err != nil {
		return nil, err
	}
	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []ssg.Section
	for rows.Next() {
		var s ssg.Section
		var layoutName sql.NullString
		err := rows.Scan(
			&s.ID, &s.ShortID, &s.Name, &s.Description, &s.Path, &s.LayoutID,
			&s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt, &layoutName,
		)
		if err != nil {
			return nil, err
		}
		s.LayoutName = layoutName.String
		sections = append(sections, s)
	}
	return sections, nil
}

func (repo *ClioRepo) GetSection(ctx context.Context, id uuid.UUID) (ssg.Section, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resSection, "Get")
	if err != nil {
		return ssg.Section{}, err
	}

	row := repo.db.QueryRowxContext(ctx, query, id)

	var (
		sectionID   uuid.UUID
		name        string
		description string
		path        string
		layoutID    uuid.UUID
		shortID     string
		createdBy   uuid.UUID
		updatedBy   uuid.UUID
		createdAt   time.Time
		updatedAt   time.Time
		layoutName  sql.NullString
	)

	err = row.Scan(
		&sectionID, &shortID, &name, &description, &path, &layoutID,
		&createdBy, &updatedBy, &createdAt, &updatedAt, &layoutName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Section{}, errors.New("section not found")
		}
		return ssg.Section{}, err
	}

	section := ssg.NewSection(name, description, path, layoutID)
	section.SetID(sectionID)
	// TODO: Remove header and blogHeader field assignments
	section.LayoutName = layoutName.String
	section.SetShortID(shortID)
	section.SetCreatedBy(createdBy)
	section.SetUpdatedBy(updatedBy)
	section.SetCreatedAt(createdAt)
	section.SetUpdatedAt(updatedAt)

	return section, nil
}

func (repo *ClioRepo) UpdateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resSection, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, section)
	return err
}

func (repo *ClioRepo) DeleteSection(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resSection, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Layout related

func (repo *ClioRepo) CreateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, "layout", "Create")
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx, query,
		layout.GetID(),
		layout.GetShortID(),
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetCreatedBy(),
		layout.GetUpdatedBy(),
		layout.GetCreatedAt(),
		layout.GetUpdatedAt(),
		layout.GetHeaderImageID(),
	)
	return err
}

func (repo *ClioRepo) GetAllLayouts(ctx context.Context) ([]ssg.Layout, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resLayout, "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layouts []ssg.Layout
	for rows.Next() {
		var (
			id            uuid.UUID
			shortID       string
			name          string
			description   string
			code          string
			createdBy     uuid.UUID
			updatedBy     uuid.UUID
			createdAt     time.Time
			updatedAt     time.Time
			headerImageID sql.NullString
		)

		err := rows.Scan(
			&id, &shortID, &name, &description, &code,
			&createdBy, &updatedBy, &createdAt, &updatedAt, &headerImageID,
		)
		if err != nil {
			return nil, err
		}

		layout := ssg.Newlayout(name, description, code)
		layout.SetID(id)
		layout.SetShortID(shortID)
		layout.SetCreatedBy(createdBy)
		layout.SetUpdatedBy(updatedBy)
		layout.SetCreatedAt(createdAt)
		layout.SetUpdatedAt(updatedAt)

		// Set header image ID if present
		if headerImageID.Valid {
			imageID, err := uuid.Parse(headerImageID.String)
			if err == nil {
				layout.SetHeaderImageID(&imageID)
			}
		}

		layouts = append(layouts, layout)
	}

	return layouts, nil
}

func (repo *ClioRepo) GetLayout(ctx context.Context, id uuid.UUID) (ssg.Layout, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, "layout", "Get")
	if err != nil {
		return ssg.Layout{}, err
	}

	var layout ssg.Layout
	err = repo.db.GetContext(ctx, &layout, query, id)
	if err != nil {
		return ssg.Layout{}, err
	}

	return layout, nil
}

func (repo *ClioRepo) UpdateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resLayout, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetUpdatedBy(),
		layout.GetUpdatedAt(),
		layout.GetHeaderImageID(),
		layout.GetID(),
	)
	return err
}

func (repo *ClioRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resLayout, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Tag related

func (repo *ClioRepo) CreateTag(ctx context.Context, tag ssg.Tag) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, tag)
	return err
}

func (repo *ClioRepo) GetTag(ctx context.Context, id uuid.UUID) (ssg.Tag, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "Get")
	if err != nil {
		return ssg.Tag{}, err
	}

	var tag ssg.Tag
	err = repo.db.GetContext(ctx, &tag, query, id)
	if err != nil {
		return ssg.Tag{}, err
	}

	return tag, nil
}

func (repo *ClioRepo) GetTagByName(ctx context.Context, name string) (ssg.Tag, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "GetByName")
	if err != nil {
		return ssg.Tag{}, err
	}

	var tag ssg.Tag
	err = repo.db.GetContext(ctx, &tag, query, name)
	if err != nil {
		return ssg.Tag{}, err
	}

	return tag, nil
}

func (repo *ClioRepo) GetAllTags(ctx context.Context) ([]ssg.Tag, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "GetAll")
	if err != nil {
		return nil, err
	}

	var tags []ssg.Tag
	err = repo.db.SelectContext(ctx, &tags, query)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *ClioRepo) UpdateTag(ctx context.Context, tag ssg.Tag) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, tag)
	return err
}

func (repo *ClioRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Param related

func (repo *ClioRepo) CreateParam(ctx context.Context, p *ssg.Param) (err error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create param query: %w", err)
	}
	if _, err = repo.db.NamedExecContext(ctx, query, p); err != nil {
		return fmt.Errorf("cannot create param: %w", err)
	}
	return nil
}

func (repo *ClioRepo) GetParam(ctx context.Context, id uuid.UUID) (ssg.Param, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "Get")
	if err != nil {
		return ssg.Param{}, fmt.Errorf("cannot get get param query: %w", err)
	}
	var param ssg.Param
	err = repo.db.GetContext(ctx, &param, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Param{}, errors.New("param not found")
		}
		return ssg.Param{}, fmt.Errorf("cannot get param: %w", err)
	}
	return param, nil
}

func (repo *ClioRepo) GetParamByName(ctx context.Context, name string) (ssg.Param, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "GetByName")
	if err != nil {
		return ssg.Param{}, fmt.Errorf("cannot get get param by name query: %w", err)
	}
	var param ssg.Param
	err = repo.db.GetContext(ctx, &param, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Param{}, errors.New("param not found")
		}
		return ssg.Param{}, fmt.Errorf("cannot get param by name: %w", err)
	}
	return param, nil
}

func (repo *ClioRepo) GetParamByRefKey(ctx context.Context, refKey string) (ssg.Param, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "GetByRefKey")
	if err != nil {
		return ssg.Param{}, fmt.Errorf("cannot get get param by ref key query: %w", err)
	}
	var param ssg.Param
	err = repo.db.GetContext(ctx, &param, query, refKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Param{}, errors.New("param not found")
		}
		return ssg.Param{}, fmt.Errorf("cannot get param by ref key: %w", err)
	}
	return param, nil
}

func (repo *ClioRepo) ListParams(ctx context.Context) ([]ssg.Param, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "List")
	if err != nil {
		return nil, fmt.Errorf("cannot get list params query: %w", err)
	}
	var params []ssg.Param
	err = repo.db.SelectContext(ctx, &params, query)
	if err != nil {
		return nil, fmt.Errorf("cannot list params: %w", err)
	}
	return params, nil
}

func (repo *ClioRepo) UpdateParam(ctx context.Context, p *ssg.Param) (err error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update param query: %w", err)
	}
	if _, err = repo.db.NamedExecContext(ctx, query, p); err != nil {
		return fmt.Errorf("cannot update param: %w", err)
	}
	return nil
}

func (repo *ClioRepo) DeleteParam(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resParam, "Delete")
	if err != nil {
		return fmt.Errorf("cannot get delete param query: %w", err)
	}
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete param: %w", err)
	}
	return nil
}

// Image related

func (repo *ClioRepo) CreateImage(ctx context.Context, img *ssg.Image) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	imageQuery, err := repo.BaseRepo.Query().Get(featSSG, resImage, "Create")
	if err != nil {
		return fmt.Errorf("cannot get create image query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, imageQuery, img); err != nil {
		return fmt.Errorf("cannot create image: %w", err)
	}

	return nil
}

func (repo *ClioRepo) GetImage(ctx context.Context, id uuid.UUID) (ssg.Image, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImage, "Get")
	if err != nil {
		return ssg.Image{}, fmt.Errorf("cannot get image query: %w", err)
	}

	var img ssg.Image
	err = repo.db.GetContext(ctx, &img, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Image{}, errors.New("image not found")
		}
		return ssg.Image{}, fmt.Errorf("cannot get image: %w", err)
	}

	return img, nil
}

func (repo *ClioRepo) GetImageByShortID(ctx context.Context, shortID string) (ssg.Image, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImage, "GetImageByShortID")
	if err != nil {
		return ssg.Image{}, fmt.Errorf("cannot get image by short ID query: %w", err)
	}

	var img ssg.Image
	err = repo.db.GetContext(ctx, &img, query, shortID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Image{}, errors.New("image not found")
		}
		return ssg.Image{}, fmt.Errorf("cannot get image by short ID: %w", err)
	}

	return img, nil
}

func (repo *ClioRepo) GetImageByContentHash(ctx context.Context, contentHash string) (ssg.Image, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImage, "GetImageByContentHash")
	if err != nil {
		return ssg.Image{}, fmt.Errorf("cannot get image by content hash query: %w", err)
	}

	var img ssg.Image
	err = repo.db.GetContext(ctx, &img, query, contentHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Image{}, errors.New("image not found")
		}
		return ssg.Image{}, fmt.Errorf("cannot get image by content hash: %w", err)
	}

	return img, nil
}

func (repo *ClioRepo) ListImages(ctx context.Context) ([]ssg.Image, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImage, "List")
	if err != nil {
		return nil, fmt.Errorf("cannot get list images query: %w", err)
	}

	var images []ssg.Image
	err = repo.db.SelectContext(ctx, &images, query)
	if err != nil {
		return nil, fmt.Errorf("cannot list images: %w", err)
	}

	return images, nil
}

func (repo *ClioRepo) UpdateImage(ctx context.Context, img *ssg.Image) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	imageQuery, err := repo.BaseRepo.Query().Get(featSSG, resImage, "Update")
	if err != nil {
		return fmt.Errorf("cannot get update image query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, imageQuery, img); err != nil {
		return fmt.Errorf("cannot update image: %w", err)
	}

	return nil
}

func (repo *ClioRepo) DeleteImage(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImage, "Delete")
	if err != nil {
		return fmt.Errorf("cannot get delete image query: %w", err)
	}
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete image: %w", err)
	}
	return nil
}

// ImageVariant related

func (repo *ClioRepo) CreateImageVariant(ctx context.Context, variant *ssg.ImageVariant) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	variantQuery, err := repo.BaseRepo.Query().Get(featSSG, resImageVariant, "CreateImageVariant")
	if err != nil {
		return fmt.Errorf("cannot get create image variant query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, variantQuery, variant); err != nil {
		return fmt.Errorf("cannot create image variant: %w", err)
	}

	return nil
}

func (repo *ClioRepo) GetImageVariant(ctx context.Context, id uuid.UUID) (ssg.ImageVariant, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImageVariant, "GetImageVariantByID")
	if err != nil {
		return ssg.ImageVariant{}, fmt.Errorf("cannot get image variant query: %w", err)
	}

	var variant ssg.ImageVariant
	err = repo.db.GetContext(ctx, &variant, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.ImageVariant{}, errors.New("image variant not found")
		}
		return ssg.ImageVariant{}, fmt.Errorf("cannot get image variant: %w", err)
	}

	return variant, nil
}

func (repo *ClioRepo) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ssg.ImageVariant, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImageVariant, "GetImageVariantsByImageID")
	if err != nil {
		return nil, fmt.Errorf("cannot get image variants by image ID query: %w", err)
	}

	var variants []ssg.ImageVariant
	err = repo.db.SelectContext(ctx, &variants, query, imageID)
	if err != nil {
		return nil, fmt.Errorf("cannot list image variants by image ID: %w", err)
	}

	return variants, nil
}

func (repo *ClioRepo) UpdateImageVariant(ctx context.Context, variant *ssg.ImageVariant) (err error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("cannot rollback transaction: %v (original error: %w)", rbErr, err)
			}
			return
		}
		err = tx.Commit()
	}()

	variantQuery, err := repo.BaseRepo.Query().Get(featSSG, resImageVariant, "UpdateImageVariant")
	if err != nil {
		return fmt.Errorf("cannot get update image variant query: %w", err)
	}
	if _, err = tx.NamedExecContext(ctx, variantQuery, variant); err != nil {
		return fmt.Errorf("cannot update image variant: %w", err)
	}

	return nil
}

func (repo *ClioRepo) DeleteImageVariant(ctx context.Context, id uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resImageVariant, "DeleteImageVariant")
	if err != nil {
		return fmt.Errorf("cannot get delete image variant query: %w", err)
	}
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete image variant: %w", err)
	}
	return nil
}

// ContentTag related

func (repo *ClioRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "AddTagToContent")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, contentID, tagID)
	return err
}

func (repo *ClioRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "RemoveTagFromContent")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, contentID, tagID)
	return err
}

func (repo *ClioRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]ssg.Tag, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "GetTagsForContent")
	if err != nil {
		return nil, err
	}

	var tags []ssg.Tag
	err = repo.db.SelectContext(ctx, &tags, query, contentID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *ClioRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]ssg.Content, error) {
	query, err := repo.BaseRepo.Query().Get(featSSG, resTag, "GetContentForTag")
	if err != nil {
		return nil, err
	}

	var contents []ssg.Content
	err = repo.db.SelectContext(ctx, &contents, query, tagID)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// ContentImage relationship methods

func (repo *ClioRepo) CreateContentImage(ctx context.Context, contentImage *ssg.ContentImage) error {
	query := `
		INSERT INTO content_images (id, content_id, image_id, purpose, position, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := repo.db.ExecContext(ctx, query,
		contentImage.ID,
		contentImage.ContentID,
		contentImage.ImageID,
		contentImage.Purpose,
		contentImage.Position,
		contentImage.IsActive,
		contentImage.CreatedAt,
		contentImage.UpdatedAt,
	)
	return err
}

func (repo *ClioRepo) DeleteContentImage(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM content_images WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *ClioRepo) GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]ssg.ContentImage, error) {
	query := `
		SELECT id, content_id, image_id, purpose, position, is_active, created_at, updated_at
		FROM content_images
		WHERE content_id = ? AND is_active = true
		ORDER BY position
	`
	var contentImages []ssg.ContentImage
	err := repo.db.SelectContext(ctx, &contentImages, query, contentID)
	return contentImages, err
}

// SectionImage relationship methods

func (repo *ClioRepo) CreateSectionImage(ctx context.Context, sectionImage *ssg.SectionImage) error {
	query := `
		INSERT INTO section_images (id, section_id, image_id, purpose, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := repo.db.ExecContext(ctx, query,
		sectionImage.ID,
		sectionImage.SectionID,
		sectionImage.ImageID,
		sectionImage.Purpose,
		sectionImage.IsActive,
		sectionImage.CreatedAt,
		sectionImage.UpdatedAt,
	)
	return err
}

func (repo *ClioRepo) DeleteSectionImage(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM section_images WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *ClioRepo) GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]ssg.SectionImage, error) {
	query := `
		SELECT id, section_id, image_id, purpose, is_active, created_at, updated_at
		FROM section_images
		WHERE section_id = ? AND is_active = true
		ORDER BY created_at
	`
	var sectionImages []ssg.SectionImage
	err := repo.db.SelectContext(ctx, &sectionImages, query, sectionID)
	return sectionImages, err
}
