package ssg

import (
	"context"
	"fmt"

	"github.com/hermesgen/hm"
)

type ParamManager struct {
	hm.Core
	repo Repo
}

// NewParamManagerWithParams creates a ParamManager with XParams.
func NewParamManager(repo Repo, params hm.XParams) *ParamManager {
	core := hm.NewCore("param-manager", params)
	return &ParamManager{
		Core: core,
		repo: repo,
	}
}

func (pm *ParamManager) FindParam(ctx context.Context, ref string) (Param, error) {
	return pm.findParamByName(ctx, ref)

}

func (pm *ParamManager) findParamByName(ctx context.Context, name string) (Param, error) {
	if pm.repo == nil {
		return Param{}, fmt.Errorf("no repository available")
	}
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByName(ctx, name)
}

func (pm *ParamManager) FindParamByRef(ctx context.Context, refKey string) (Param, error) {
	if pm.repo == nil {
		return Param{}, fmt.Errorf("no repository available")
	}
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByRefKey(ctx, refKey)
}

func (pm *ParamManager) Get(ctx context.Context, refKey string, defVal string) string {
	if pm.repo == nil {
		// No repo available, fallback to configuration
		return pm.Cfg().StrValOrDef(refKey, defVal)
	}

	param, err := pm.repo.GetParamByRefKey(ctx, refKey)
	if err == nil && !param.IsZero() {
		return param.Value
	}

	// Fallback to configuration
	return pm.Cfg().StrValOrDef(refKey, defVal)
}

// GetSiteMode returns the current site mode (structured or blog).
// Returns "structured" by default if not set.
func (pm *ParamManager) GetSiteMode(ctx context.Context) string {
	mode := pm.Get(ctx, "site.mode", "structured")
	if mode != "structured" && mode != "blog" {
		pm.Log().Error("Invalid site mode, defaulting to structured", "mode", mode)
		return "structured"
	}
	return mode
}

// SetSiteMode sets the site mode to either "structured" or "blog".
func (pm *ParamManager) SetSiteMode(ctx context.Context, mode string) error {
	if pm.repo == nil {
		return fmt.Errorf("no repository available")
	}

	if mode != "structured" && mode != "blog" {
		return fmt.Errorf("invalid site mode: must be 'structured' or 'blog'")
	}

	param, err := pm.repo.GetParamByRefKey(ctx, "site.mode")
	if err != nil || param.IsZero() {
		// Create new param
		param = NewParam("Site Mode", mode)
		param.Description = "Site operation mode: 'structured' (multi-section) or 'blog' (single chronological feed)"
		param.RefKey = "site.mode"
		param.System = 1
		param.GenCreateValues()
		return pm.repo.CreateParam(ctx, &param)
	}

	// Update existing
	param.Value = mode
	param.GenUpdateValues()
	return pm.repo.UpdateParam(ctx, &param)
}
