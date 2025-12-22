package auth

import (
	"context"
	"embed"
	"net/http"

	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
)

const (
	resUserName    = "user"
	resUserNameCap = "User"
)

type DBProvider interface {
	GetDB() *sqlx.DB
}

type RepoFactory func(qm *hm.QueryManager, db *sqlx.DB, params hm.XParams) Repo

type APIHandler struct {
	*hm.APIHandler
	svc         Service
	dbProvider  DBProvider
	assetsFS    embed.FS
	engine      string
	repoFactory RepoFactory
}

func NewAPIHandler(name string, dbProvider DBProvider, assetsFS embed.FS, engine string, repoFactory RepoFactory, params hm.XParams) *APIHandler {
	h := hm.NewAPIHandler(name, params)
	return &APIHandler{
		APIHandler:  h,
		dbProvider:  dbProvider,
		assetsFS:    assetsFS,
		engine:      engine,
		repoFactory: repoFactory,
	}
}

func (h *APIHandler) Setup(ctx context.Context) error {
	params := hm.XParams{Cfg: h.Cfg(), Log: h.Log()}
	qm := hm.NewQueryManager(h.assetsFS, h.engine, params)
	if err := qm.Setup(ctx); err != nil {
		return err
	}

	repo := h.repoFactory(qm, h.dbProvider.GetDB(), params)
	h.svc = NewService(repo, params)
	return nil
}

func (h *APIHandler) OK(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.OK(w, message, wrappedData)
}

func (h *APIHandler) Created(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.Created(w, message, wrappedData)
}

func (h *APIHandler) wrapData(data interface{}) interface{} {
	switch v := data.(type) {
	// Single entities
	case User:
		return map[string]interface{}{"user": v}

	// Slices of entities
	case []User:
		return map[string]interface{}{"users": v}

	// Default case for nil, maps, or other types
	default:
		return data
	}
}
