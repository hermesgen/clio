package main

import (
	"context"
	"embed"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hermesgen/clio/internal/core"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/clio/internal/repo/sqlite"
	webssg "github.com/hermesgen/clio/internal/web/ssg"
	"github.com/hermesgen/hm"
	"github.com/hermesgen/hm/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	name      = "clio"
	version   = "v1"
	namespace = "clio"
	engine    = "sqlite"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	flag.Parse()

	ctx := context.Background()
	log := hm.NewLogger("info")
	cfg := hm.LoadCfg(namespace, hm.Flags)
	xparams := hm.XParams{Cfg: cfg, Log: log}

	fm := hm.NewFlashManager(xparams)
	workspace := core.NewWorkspace(xparams)
	app := core.NewApp(name, version, assetsFS, xparams)
	app.PreviewHandler = core.NewMultiSitePreviewHandler(xparams)
	templateManager := hm.NewTemplateManager(assetsFS, xparams)
	fileServer := hm.NewFileServer(assetsFS, xparams)

	sitesMigrator := hm.NewMigrator(assetsFS, engine, xparams)
	sitesMigrator.SetPath("assets/migration/sqlite-sites")
	adminDBManager := core.NewAdminDBManager(assetsFS, engine, sitesMigrator, xparams)

	repoFactory := func(qm *hm.QueryManager, params hm.XParams) ssg.Repo {
		return sqlite.NewClioRepo(qm, params)
	}

	repoManager := ssg.NewRepoManager(assetsFS, engine, repoFactory, xparams)
	sessionManager := auth.NewSessionManager(xparams)
	siteManager := ssg.NewSiteManager(adminDBManager, assetsFS, engine, repoFactory, xparams)

	apiRouter := hm.NewAPIRouter("api-router", xparams)
	gitClient := github.NewClient(xparams)
	ssgPublisher := ssg.NewPublisher(gitClient, xparams)
	ssgGenerator := ssg.NewGenerator(xparams)

	app.Add(workspace)
	app.Add(adminDBManager)
	app.Add(sitesMigrator)
	app.Add(fm)
	app.Add(fileServer)
	app.Add(templateManager)
	app.Add(sessionManager)
	app.Add(repoManager)
	app.Add(siteManager)
	app.Add(gitClient)
	app.Add(ssgPublisher)
	app.Add(ssgGenerator)
	app.Add(apiRouter)

	siteContextMw := ssg.NewSiteContextMw(sessionManager, siteManager, repoManager, xparams)

	authRepoFactory := func(qm *hm.QueryManager, db *sqlx.DB, params hm.XParams) auth.Repo {
		repo := sqlite.NewClioRepo(qm, params)
		repo.SetDB(db)
		return repo
	}
	authAPIHandler := auth.NewAPIHandler("auth-api-handler", adminDBManager, assetsFS, engine, authRepoFactory, xparams)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, []hm.Middleware{}, xparams)
	app.Add(authAPIHandler)
	app.Add(authAPIRouter)

	paramManager := ssg.NewParamManager(nil, xparams)
	imageManager := ssg.NewImageManager(xparams)
	ssgAPIService := ssg.NewService(assetsFS, nil, ssgGenerator, ssgPublisher, paramManager, imageManager, xparams)
	ssgAPIHandler := ssg.NewAPIHandler("ssg-api-handler", ssgAPIService, siteManager, xparams)
	ssgAPIRouter := ssg.NewAPIRouter(ssgAPIHandler, []hm.Middleware{hm.CORSMw, siteContextMw.APIHandler}, xparams)
	app.Add(ssgAPIHandler)
	app.Add(ssgAPIRouter)

	ssgWebHandler := webssg.NewWebHandler(templateManager, fm, paramManager, siteManager, sessionManager, xparams)
	ssgWebRouter := webssg.NewWebRouter(ssgWebHandler, append(fm.Middlewares(), siteContextMw.WebHandler), xparams)

	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			ssgWebHandler.RootRedirect(w, r)
		}
	})

	if err := app.Setup(ctx); err != nil {
		log.Errorf("Cannot setup application: %v", err)
		return
	}

	app.MountAPI("/api/v1/auth", authAPIRouter)
	app.MountAPI("/api/v1/ssg", ssgAPIRouter)
	app.MountWeb("/ssg", ssgWebRouter)
	app.MountFileServer("/", fileServer)

	log.Infof("%s(%s) setup completed", name, version)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		err := app.Start(ctx)
		if err != nil {
			log.Errorf("Cannot start %s(%s): %v", name, version, err)
		}
	}()

	log.Infof("%s(%s) started successfully", name, version)

	<-stop

	log.Infof("Shutting down %s(%s)...", name, version)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := app.Stop(shutdownCtx)
	if err != nil {
		log.Errorf("Error during shutdown: %v", err)
	} else {
		log.Infof("%s(%s) stopped gracefully", name, version)
	}
}
