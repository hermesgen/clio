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
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByName(ctx, name)
}

func (pm *ParamManager) FindParamByRef(ctx context.Context, refKey string) (Param, error) {
	// TODO: This is a simple wrapper for now but a simple caching strategy can be added here.
	return pm.repo.GetParamByRefKey(ctx, refKey)
}

func (pm *ParamManager) Get(ctx context.Context, refKey string, defVal string) string {
	param, err := pm.repo.GetParamByRefKey(ctx, refKey)
	if err == nil && !param.IsZero() {
		return param.Value
	}

	// Fallback to configuration
	return pm.Cfg().StrValOrDef(refKey, defVal)
}

// GetSiteMode returns the current site mode (normal or blog).
// Returns "normal" by default if not set.
func (pm *ParamManager) GetSiteMode(ctx context.Context) string {
	mode := pm.Get(ctx, "site.mode", "normal")
	if mode != "normal" && mode != "blog" {
		pm.Log().Error("Invalid site mode, defaulting to normal", "mode", mode)
		return "normal"
	}
	return mode
}

// SetSiteMode sets the site mode to either "normal" or "blog".
func (pm *ParamManager) SetSiteMode(ctx context.Context, mode string) error {
	if mode != "normal" && mode != "blog" {
		return fmt.Errorf("invalid site mode: must be 'normal' or 'blog'")
	}

	param, err := pm.repo.GetParamByRefKey(ctx, "site.mode")
	if err != nil || param.IsZero() {
		// Create new param
		param = NewParam("Site Mode", mode)
		param.Description = "Site operation mode: 'normal' (multi-section) or 'blog' (single chronological feed)"
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
