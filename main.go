package main

import (
	"context"
	"embed"
	"net/http"

	"github.com/adrianpk/clio/internal/am"
	"github.com/adrianpk/clio/internal/am/github"
	"github.com/adrianpk/clio/internal/core"
	"github.com/adrianpk/clio/internal/feat/auth"
	"github.com/adrianpk/clio/internal/feat/ssg"
	"github.com/adrianpk/clio/internal/repo/sqlite"
	webssg "github.com/adrianpk/clio/internal/web/ssg"
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
	ctx := context.Background()
	log := am.NewLogger("info")
	cfg := am.LoadCfg(namespace, am.Flags)
	
	// opts := am.DefOpts(log, cfg) // REMOVED - migrating to XParams
	
	// XParams for components that only need log + config
	xparams := am.XParams{Cfg: cfg, Log: log}
	
	// Opts for legacy components still using variadic pattern
	opts := []am.Option{}

	fm := am.NewFlashManager()
	workspace := core.NewWorkspace(opts...)
	app := core.NewAppWithParams(name, version, assetsFS, xparams)
	queryManager := am.NewQueryManagerWithParams(assetsFS, engine, xparams)
	templateManager := am.NewTemplateManagerWithParams(assetsFS, xparams)
	repo := sqlite.NewClioRepo(queryManager)
	migrator := am.NewMigrator(assetsFS, engine)
	fileServer := am.NewFileServerWithParams(assetsFS, xparams)

	app.MountFileServer("/", fileServer)

	// Serve uploaded images from the filesystem
	imagesPath := cfg.StrValOrDef(am.Key.SSGImagesPath, "_workspace/documents/assets/images")
	imageFileServer := http.FileServer(http.Dir(imagesPath))
	app.Router.Handle("/static/images/*", http.StripPrefix("/static/images/", imageFileServer))

	apiRouter := am.NewAPIRouterWithParams("api-router", xparams)

	// GitAuth feature
	authSeeder := auth.NewSeeder(assetsFS, engine, repo)
	authService := auth.NewServiceWithParams(repo, xparams)
	authAPIHandler := auth.NewAPIHandler("auth-api-handler", authService, opts...)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, nil) // No middleware for now
	apiRouter.Mount("/auth", authAPIRouter)

	// SSG feature
	gitClient := github.NewClientWithParams(xparams)
	ssgPublisher := ssg.NewPublisherWithParams(gitClient, xparams)
	ssgSeeder := ssg.NewSeeder(assetsFS, engine, repo)
	ssgGenerator := ssg.NewGeneratorWithParams(xparams)
	ssgParamManager := ssg.NewParamManagerWithParams(repo, xparams)
	ssgImageManager := ssg.NewImageManagerWithParams(xparams)
	ssgService := ssg.NewServiceWithParams(assetsFS, repo, ssgGenerator, ssgPublisher, ssgParamManager, ssgImageManager, xparams)
	ssgAPIHandler := ssg.NewAPIHandler("ssg-api-handler", ssgService)
	ssgAPIRouter := ssg.NewAPIRouter(ssgAPIHandler, []am.Middleware{am.CORSMw})
	apiRouter.Mount("/ssg", ssgAPIRouter)

app.MountAPI("/api/v1", apiRouter)

	// Web app
	ssgWebHandler := webssg.NewWebHandler(templateManager, fm, opts...)
	ssgWebRouter := webssg.NewWebRouter(ssgWebHandler, append(fm.Middlewares(), am.LogHeadersMw))

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

	err = app.Start(ctx)
	if err != nil {
		log.Errorf("Cannot start %s(%s): %v", name, version, err)
	}
}
