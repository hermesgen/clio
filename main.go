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
	migrator := hm.NewMigrator(assetsFS, engine, xparams)
	migrator.SetPath("assets/migration/sqlite")
	dbManager := core.NewAdminDBManager(assetsFS, engine, migrator, xparams)
	sessionManager := auth.NewSessionManager(xparams)
	dynamicImageServer := core.NewDynamicImageServer(xparams)
	apiRouter := hm.NewAPIRouter("api-router", xparams)
	gitClient := github.NewClient(xparams)
	ssgPublisher := ssg.NewPublisher(gitClient, xparams)
	ssgGenerator := ssg.NewGenerator(xparams)
	qm := hm.NewQueryManager(assetsFS, engine, xparams)
	clioRepo := sqlite.NewClioRepo(qm, xparams)
	siteManager := ssg.NewSiteManager(clioRepo, assetsFS, engine, xparams)
	siteContextMw := ssg.NewSiteContextMw(sessionManager, siteManager, xparams)
	authSeeder := auth.NewSeeder(assetsFS, engine, clioRepo, xparams)
	ssgSeeder := ssg.NewSeeder(assetsFS, engine, clioRepo, xparams)
	paramManager := ssg.NewParamManager(clioRepo, xparams)
	imageManager := ssg.NewImageManager(xparams)
	ssgAPIService := ssg.NewService(assetsFS, clioRepo, ssgGenerator, ssgPublisher, paramManager, imageManager, xparams)
	ssgAPIHandler := ssg.NewAPIHandler("ssg-api-handler", ssgAPIService, siteManager, xparams)
	ssgAPIRouter := ssg.NewAPIRouter(ssgAPIHandler, []hm.Middleware{hm.CORSMw, siteContextMw.APIHandler}, xparams)

	authAPIHandler := auth.NewAPIHandler("auth-api-handler", clioRepo, xparams)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, []hm.Middleware{}, xparams)

	app.Add(workspace)
	app.Add(dbManager)
	app.Add(migrator)
	app.Add(qm)
	app.Add(clioRepo)
	app.Add(fm)
	app.Add(fileServer)
	app.Add(dynamicImageServer)
	app.Add(templateManager)
	app.Add(sessionManager)
	app.Add(siteManager)
	app.Add(gitClient)
	app.Add(ssgPublisher)
	app.Add(ssgGenerator)
	app.Add(apiRouter)
	app.Add(authSeeder)
	app.Add(ssgSeeder)
	app.Add(authAPIHandler)
	app.Add(authAPIRouter)
	app.Add(ssgAPIHandler)
	app.Add(ssgAPIRouter)

	ssgWebHandler := webssg.NewWebHandler(templateManager, fm, paramManager, siteManager, sessionManager, xparams)
	ssgWebRouter := webssg.NewWebRouter(ssgWebHandler, append(fm.Middlewares(), siteContextMw.WebHandler), xparams)

	// TODO: This also needs to be handled by lifecycle hooks
	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			ssgWebHandler.RootRedirect(w, r)
		}
	})

	if err := app.Setup(ctx); err != nil {
		log.Errorf("Cannot setup application: %v", err)
		return
	}

	app.Router.HandleFunc("/static/images/*", dynamicImageServer.Handler())
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
