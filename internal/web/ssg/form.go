package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

// ContentForm represents the form data for a content.
type ContentForm struct {
	*hm.BaseForm

	// Content fields
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	SectionID   string `json:"section_id"`
	Kind        string `json:"kind"`
	Heading     string `json:"heading"`
	Body        string `json:"body"`
	Image       string `json:"image"`
	Draft       bool   `json:"draft"`
	Featured    bool   `json:"featured"`
	PublishedAt string `json:"published_at"`
	Tags        string `json:"tags"`

	// Meta fields
	Description     string `json:"description"`
	Keywords        string `json:"keywords"`
	Robots          string `json:"robots"`
	CanonicalURL    string `json:"canonical_url"`
	Sitemap         string `json:"sitemap"`
	TableOfContents bool   `json:"table_of_contents"`
	Share           bool   `json:"share"`
	Comments        bool   `json:"comments"`
}

// NewContentForm creates a new ContentForm from a request.
func NewContentForm(r *http.Request) ContentForm {
	return ContentForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// ContentFormFromRequest creates a ContentForm from an HTTP request.
func ContentFormFromRequest(r *http.Request) (ContentForm, error) {
	if err := r.ParseForm(); err != nil {
		return ContentForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewContentForm(r) // Initialize with BaseForm
	form.ID = r.Form.Get("id")
	form.UserID = r.Form.Get("user_id")
	form.SectionID = r.Form.Get("section_id")
	form.Kind = r.Form.Get("kind")
	form.Heading = r.Form.Get("heading")
	form.Body = r.Form.Get("body")
	form.Image = r.Form.Get("image")
	form.Tags = r.Form.Get("tags")
	form.Draft, _ = strconv.ParseBool(r.Form.Get("draft"))
	form.Featured, _ = strconv.ParseBool(r.Form.Get("featured"))
	form.PublishedAt = r.Form.Get("published_at")

	// Meta fields
	form.Description = r.Form.Get("description")
	form.Keywords = r.Form.Get("keywords")
	form.Robots = r.Form.Get("robots")
	form.CanonicalURL = r.Form.Get("canonical_url")
	form.Sitemap = r.Form.Get("sitemap")
	form.TableOfContents, _ = strconv.ParseBool(r.Form.Get("table_of_contents"))
	form.Share, _ = strconv.ParseBool(r.Form.Get("share"))
	form.Comments, _ = strconv.ParseBool(r.Form.Get("comments"))

	return form, nil
}

// ToFeatContent converts a ContentForm to a feat.Content model.
func ToFeatContent(form ContentForm) feat.Content {
	content := feat.NewContent(form.Heading, form.Body)

	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			content.ID = id
		}
	}

	if form.UserID != "" {
		userID, err := uuid.Parse(form.UserID)
		if err == nil {
			content.UserID = userID
		}
	}

	if form.SectionID != "" {
		sectionID, err := uuid.Parse(form.SectionID)
		if err == nil {
			content.SectionID = sectionID
		}
	}

	content.Kind = form.Kind
	// TODO: Handle image via relationship
	content.Draft = form.Draft
	content.Featured = form.Featured

	if form.PublishedAt != "" {
		// Try parsing multiple formats, starting with RFC3339
		formats := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02"}
		for _, format := range formats {
			pubAt, err := time.Parse(format, form.PublishedAt)
			if err == nil {
				content.PublishedAt = &pubAt
				break
			}
		}
	}

	// New part for tags
	if form.Tags != "" {
		var tags []struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal([]byte(form.Tags), &tags); err == nil {
			for _, t := range tags {
				content.Tags = append(content.Tags, feat.Tag{Name: t.Value})
			}
		} else {
			// Handle plain comma-separated tags
			tagNames := strings.Split(form.Tags, ",")
			for _, name := range tagNames {
				if trimmedName := strings.TrimSpace(name); trimmedName != "" {
					content.Tags = append(content.Tags, feat.Tag{Name: trimmedName})
				}
			}
		}
	}

	// Meta
	meta := feat.NewMeta(content.ID)
	meta.Description = form.Description
	meta.Keywords = form.Keywords
	meta.Robots = form.Robots
	meta.CanonicalURL = form.CanonicalURL
	meta.Sitemap = form.Sitemap
	meta.TableOfContents = form.TableOfContents
	meta.Share = form.Share
	meta.Comments = form.Comments
	content.Meta = meta

	return content
}

// ToContentForm converts a feat.Content model to a ContentForm.
func ToContentForm(r *http.Request, content feat.Content) ContentForm {
	form := NewContentForm(r) // Initialize with BaseForm
	form.ID = content.GetID().String()
	form.UserID = content.UserID.String()
	form.SectionID = content.SectionID.String()
	form.Kind = content.Kind
	form.Heading = content.Heading
	form.Body = content.Body
	form.Image = "" // TODO: Get image via relationship
	form.Draft = content.Draft
	form.Featured = content.Featured
	if content.PublishedAt != nil {
		form.PublishedAt = content.PublishedAt.Format(time.RFC3339) // Preserve original format with timezone
	}

	// Create a comma-separated string of tag names
	tagNames := make([]string, len(content.Tags))
	for i, tag := range content.Tags {
		tagNames[i] = tag.Name
	}
	form.Tags = strings.Join(tagNames, ",")

	// Meta
	form.Description = content.Meta.Description
	form.Keywords = content.Meta.Keywords
	form.Robots = content.Meta.Robots
	form.CanonicalURL = content.Meta.CanonicalURL
	form.Sitemap = content.Meta.Sitemap
	form.TableOfContents = content.Meta.TableOfContents
	form.Share = content.Meta.Share
	form.Comments = content.Meta.Comments

	return form
}

// Validate validates the ContentForm.
func (f *ContentForm) Validate() {
	validation := f.Validation()
	if f.Heading == "" {
		validation.AddFieldError("heading", f.Heading, "Heading cannot be empty")
	}
	f.SetValidation(validation)
}

// LayoutForm represents the form for creating or updating a layout.
type LayoutForm struct {
	*hm.BaseForm
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

// NewLayoutForm creates a new LayoutForm.
func NewLayoutForm(r *http.Request) LayoutForm {
	return LayoutForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// LayoutFormFromRequest creates a LayoutForm from an HTTP request.
func LayoutFormFromRequest(r *http.Request) (LayoutForm, error) {
	if err := r.ParseForm(); err != nil {
		return LayoutForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewLayoutForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Code = r.Form.Get("code")

	return form, nil
}

// ToFeatLayout converts a LayoutForm to a feat.Layout model.
func ToFeatLayout(form LayoutForm) feat.Layout {
	layout := feat.Newlayout(form.Name, form.Description, form.Code)
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			layout.ID = id
		}
	}
	return layout
}

// ToLayoutForm converts a feat.Layout model to a LayoutForm.
func ToLayoutForm(r *http.Request, layout feat.Layout) LayoutForm {
	form := NewLayoutForm(r)
	form.ID = layout.GetID().String()
	form.Name = layout.Name
	form.Description = layout.Description
	form.Code = layout.Code

	return form
}

// Validate validates the LayoutForm.
func (f *LayoutForm) Validate() {
	validation := f.Validation()
	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name is required")
	}
	if f.Code == "" {
		validation.AddFieldError("code", f.Code, "Code is required")
	}
	f.SetValidation(validation)
}

// SectionForm represents the form for creating or updating a section.
type SectionForm struct {
	*hm.BaseForm
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	LayoutID    string `json:"layout_id"`
	Header      string `json:"header"`
	BlogHeader  string `json:"blog_header"`
}

// NewSectionForm creates a new SectionForm.
func NewSectionForm(r *http.Request) SectionForm {
	return SectionForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// SectionFormFromRequest creates a SectionForm from an HTTP request.
func SectionFormFromRequest(r *http.Request) (SectionForm, error) {
	if err := r.ParseForm(); err != nil {
		return SectionForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewSectionForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Path = r.Form.Get("path")
	form.LayoutID = r.Form.Get("layout_id")
	form.Header = r.Form.Get("header")
	form.BlogHeader = r.Form.Get("blog_header")

	return form, nil
}

// ToFeatSection converts a SectionForm to a feat.Section model.
func ToFeatSection(form SectionForm) feat.Section {
	layoutID, _ := uuid.Parse(form.LayoutID)
	section := feat.NewSection(form.Name, form.Description, form.Path, layoutID)
	// TODO: Handle header and blog header via relationships
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			section.ID = id
		}
	}

	return section
}

// ToSectionForm converts a feat.Section model to a SectionForm.
func ToSectionForm(r *http.Request, section feat.Section) SectionForm {
	form := NewSectionForm(r)
	form.ID = section.GetID().String()
	form.Name = section.Name
	form.Description = section.Description
	form.Path = section.Path
	form.LayoutID = section.LayoutID.String()
	form.Header = ""     // TODO: Get header via relationship
	form.BlogHeader = "" // TODO: Get blog header via relationship
	return form
}

// Validate validates the SectionForm.
func (f *SectionForm) Validate() {
	validation := f.Validation()
	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name is required")
	}

	if f.Path == "" {
		validation.AddFieldError("path", f.Path, "Path is required")
	}

	if f.LayoutID == "" {
		validation.AddFieldError("layout_id", f.LayoutID, "Layout is required")
	}
	f.SetValidation(validation)
}

// TagForm represents the form data for a tag.
type TagForm struct {
	*hm.BaseForm
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewTagForm creates a new TagForm from a request.
func NewTagForm(r *http.Request) TagForm {
	return TagForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// TagFormFromRequest creates a TagForm from an HTTP request.
func TagFormFromRequest(r *http.Request) (TagForm, error) {
	if err := r.ParseForm(); err != nil {
		return TagForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewTagForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")

	return form, nil
}

// ToFeatTag converts a TagForm to a feat.Tag model.
func ToFeatTag(form TagForm) feat.Tag {
	tag := feat.NewTag(form.Name)
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			tag.ID = id
		}
	}
	return tag
}

// ToTagForm converts a feat.Tag model to a TagForm.
func ToTagForm(r *http.Request, featTag feat.Tag) TagForm {
	form := NewTagForm(r)
	form.ID = featTag.GetID().String()
	form.Name = featTag.Name
	return form
}

// Validate validates the TagForm.
func (f *TagForm) Validate() {
	validation := f.Validation()
	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name cannot be empty")
	}
	f.SetValidation(validation)
}

// ParamForm represents the form data for a param.
type ParamForm struct {
	*hm.BaseForm
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	RefKey      string `json:"ref_key"`
}

// NewParamForm creates a new ParamForm from a request.
func NewParamForm(r *http.Request) ParamForm {
	return ParamForm{
		BaseForm: hm.NewBaseForm(r),
	}
}

// ParamFormFromRequest creates a ParamForm from an HTTP request.
func ParamFormFromRequest(r *http.Request) (ParamForm, error) {
	if err := r.ParseForm(); err != nil {
		return ParamForm{}, fmt.Errorf("error parsing form: %w", err)
	}

	form := NewParamForm(r)
	form.ID = r.Form.Get("id")
	form.Name = r.Form.Get("name")
	form.Description = r.Form.Get("description")
	form.Value = r.Form.Get("value")
	form.RefKey = r.Form.Get("ref_key")

	return form, nil
}

// ToFeatParam converts a ParamForm to a feat.Param model.
func ToFeatParam(form ParamForm) feat.Param {
	param := feat.NewParam(form.Name, form.Value)
	param.Description = form.Description
	param.RefKey = form.RefKey
	if form.ID != "" {
		id, err := uuid.Parse(form.ID)
		if err == nil {
			param.ID = id
		}
	}
	return param
}

// ToParamForm converts a feat.Param model to a ParamForm.
func ToParamForm(r *http.Request, featParam feat.Param) ParamForm {
	form := NewParamForm(r)
	form.ID = featParam.GetID().String()
	form.Name = featParam.Name
	form.Description = featParam.Description
	form.Value = featParam.Value
	form.RefKey = featParam.RefKey
	return form
}

// Validate validates the ParamForm.
func (f *ParamForm) Validate() {
	validation := f.Validation()
	if f.Name == "" {
		validation.AddFieldError("name", f.Name, "Name is required")
	}
	if f.Value == "" {
		validation.AddFieldError("value", f.Value, "Value is required")
	}
	f.SetValidation(validation)
}
