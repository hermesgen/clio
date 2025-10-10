package ssg

import (
	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

// Param model for web layer.
type Param struct {
	ID          uuid.UUID `json:"id"`
	ShortID     string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	RefKey      string    `json:"ref_key"`
	System      int       `json:"system"`
}

// NewParam creates a new Param for the web layer.
func NewParam(name, value string) Param {
	return Param{
		Name:  name,
		Value: value,
	}
}

// Type returns the type of the entity.
func (p *Param) Type() string {
	return "param"
}

// GetID returns the unique identifier of the entity.
func (p *Param) GetID() uuid.UUID {
	return p.ID
}

// GenID delegates to the functional helper.
func (p *Param) GenID() {
	hm.GenID(p)
}

// SetID sets the unique identifier of the entity.
func (p *Param) SetID(id uuid.UUID, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ID == uuid.Nil || (shouldForce && id != uuid.Nil) {
		p.ID = id
	}
}

// GetShortID returns the short ID portion of the slug.
func (p *Param) GetShortID() string {
	return p.ShortID
}

// GenShortID delegates to the functional helper.
func (p *Param) GenShortID() {
	hm.GenShortID(p)
}

// SetShortID sets the short ID of the entity.
func (p *Param) SetShortID(shortID string, force ...bool) {
	shouldForce := len(force) > 0 && force[0]
	if p.ShortID == "" || shouldForce {
		p.ShortID = shortID
	}
}

// TypeID returns a universal identifier for a specific model instance.
func (p *Param) TypeID() string {
	return hm.Normalize(p.Type()) + "-" + p.GetShortID()
}

// IsZero returns true if the Param is uninitialized.
func (p *Param) IsZero() bool {
	return p.ID == uuid.Nil
}

// Slug returns a slug for the param.
func (p *Param) Slug() string {
	return hm.Normalize(p.Name) + "-" + p.GetShortID()
}

// IsSystem returns true if the param is a system parameter.
func (p *Param) IsSystem() bool {
	return p.System == 1
}

// ToWebParam converts a feat.Param model to a web.Param model.
func ToWebParam(featParam feat.Param) Param {
	return Param{
		ID:          featParam.ID,
		ShortID:     featParam.ShortID,
		Name:        featParam.Name,
		Description: featParam.Description,
		Value:       featParam.Value,
		RefKey:      featParam.RefKey,
		System:      featParam.System,
	}
}

// ToWebParams converts a slice of feat.Param models to a slice of web.Param models.
func ToWebParams(featParams []feat.Param) []Param {
	webParams := make([]Param, len(featParams))
	for i, featParam := range featParams {
		webParams[i] = ToWebParam(featParam)
	}
	return webParams
}
