package ssg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
	"github.com/hermesgen/clio/internal/feat/auth"
)

type Repo interface {
	hm.Repo

	CreateContent(ctx context.Context, content *Content) error
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, content *Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error
	GetAllContentWithMeta(ctx context.Context) ([]Content, error)

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
	GetImageByContentHash(ctx context.Context, contentHash string) (Image, error)
	UpdateImage(ctx context.Context, image *Image) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
	ListImages(ctx context.Context) ([]Image, error)

	// ImageVariant related
	CreateImageVariant(ctx context.Context, variant *ImageVariant) error
	GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error)
	UpdateImageVariant(ctx context.Context, variant *ImageVariant) error
	DeleteImageVariant(ctx context.Context, id uuid.UUID) error
	ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error)

	// ContentImage relationship methods
	CreateContentImage(ctx context.Context, contentImage *ContentImage) error
	DeleteContentImage(ctx context.Context, id uuid.UUID) error
	GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]ContentImage, error)

	// SectionImage relationship methods
	CreateSectionImage(ctx context.Context, sectionImage *SectionImage) error
	DeleteSectionImage(ctx context.Context, id uuid.UUID) error
	GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]SectionImage, error)

	AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error
	RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error
	GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error)
	GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error)

	GetUserByUsername(ctx context.Context, username string) (auth.User, error)
}
