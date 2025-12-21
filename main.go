package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
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
	var initBlogMode bool
	flag.BoolVar(&initBlogMode, "init-blog-mode", false, "Initialize site in blog mode (sets site.mode=blog in database)")
	flag.Parse()

	ctx := context.Background()
	log := hm.NewLogger("info")
	cfg := hm.LoadCfg(namespace, hm.Flags)

	// XParams for components that only need log + config
	xparams := hm.XParams{Cfg: cfg, Log: log}

	fm := hm.NewFlashManager(xparams)
	workspace := core.NewWorkspace(xparams)
	app := core.NewApp(name, version, assetsFS, xparams)
	queryManager := hm.NewQueryManager(assetsFS, engine, xparams)
	templateManager := hm.NewTemplateManager(assetsFS, xparams)
	repo := sqlite.NewClioRepo(queryManager, xparams)
	migrator := hm.NewMigrator(assetsFS, engine, xparams)
	fileServer := hm.NewFileServer(assetsFS, xparams)

	// Redirect root to SSG list content
	app.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/ssg/list-content", http.StatusFound)
			return
		}
	})

	app.MountFileServer("/", fileServer)

	// Serve uploaded images from the filesystem
	imagesPath := cfg.StrValOrDef(ssg.SSGKey.ImagesPath, "_workspace/documents/assets/images")
	imageFileServer := http.FileServer(http.Dir(imagesPath))
	app.Router.Handle("/static/images/*", http.StripPrefix("/static/images/", imageFileServer))

	apiRouter := hm.NewAPIRouter("api-router", xparams)

	// GitAuth feature
	authSeeder := auth.NewSeeder(assetsFS, engine, repo, xparams)
	authService := auth.NewService(repo, xparams)
	authAPIHandler := auth.NewAPIHandler("auth-api-handler", authService, xparams)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, nil, xparams) // No middleware for now
	apiRouter.Mount("/auth", authAPIRouter)

	// SSG feature
	gitClient := github.NewClient(xparams)
	ssgPublisher := ssg.NewPublisher(gitClient, xparams)
	ssgSeeder := ssg.NewSeeder(assetsFS, engine, repo, xparams)
	ssgGenerator := ssg.NewGenerator(xparams)
	ssgParamManager := ssg.NewParamManager(repo, xparams)
	ssgImageManager := ssg.NewImageManager(xparams)
	ssgService := ssg.NewService(assetsFS, repo, ssgGenerator, ssgPublisher, ssgParamManager, ssgImageManager, xparams)
	ssgAPIHandler := ssg.NewAPIHandler("ssg-api-handler", ssgService, xparams)
	ssgAPIRouter := ssg.NewAPIRouter(ssgAPIHandler, []hm.Middleware{hm.CORSMw}, xparams)
	apiRouter.Mount("/ssg", ssgAPIRouter)

	app.MountAPI("/api/v1", apiRouter)

	// Web app
	ssgWebHandler := webssg.NewWebHandler(templateManager, fm, ssgParamManager, xparams)
	ssgWebRouter := webssg.NewWebRouter(ssgWebHandler, append(fm.Middlewares(), hm.LogHeadersMw), xparams)

	app.MountWeb("/ssg", ssgWebRouter)

	// Add deps
	app.Add(workspace)
	app.Add(migrator)
	app.Add(fm)
	app.Add(fileServer)
	app.Add(queryManager)
	app.Add(templateManager)
	app.Add(repo)
	app.Add(ssgWebHandler)
	app.Add(ssgWebRouter)
	app.Add(authSeeder)
	app.Add(ssgSeeder)
	app.Add(gitClient)
	app.Add(ssgPublisher)
	app.Add(ssgGenerator)
	app.Add(ssgParamManager)
	app.Add(ssgService)
	app.Add(ssgAPIHandler)
	app.Add(ssgAPIRouter)
	app.Add(apiRouter)
	app.Add(authSeeder)

	err := app.Setup(ctx)
	if err != nil {
		log.Errorf("Cannot setup %s(%s): %v", name, version, err)
		return
	}

	// Handle --init-blog-mode flag
	if initBlogMode {
		err = ssgParamManager.SetSiteMode(ctx, "blog")
		if err != nil {
			log.Errorf("Failed to set blog mode: %v", err)
			fmt.Fprintf(os.Stderr, "Error: Failed to set blog mode: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully initialized site in blog mode")
		fmt.Println("Site mode has been set to 'blog' in the database")
		return
	}

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

	err = app.Stop(shutdownCtx)
	if err != nil {
		log.Errorf("Error during shutdown: %v", err)
	} else {
		log.Infof("%s(%s) stopped gracefully", name, version)
	}
}
