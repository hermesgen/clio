package core

import (
	"context"
	"embed"
	"os"
	"os/signal"

	"github.com/hermesgen/hm"
)

type App struct {
	*hm.App
}

func NewApp(name, version string, fs embed.FS, params hm.XParams) *App {
	core := hm.NewApp(name, version, fs, params)
	app := &App{
		App: core,
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	err := app.App.Start(ctx)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	return app.Core.Stop(ctx)
}

func (app *App) Setup(ctx context.Context) error {
	app.setupPluralization()
	return app.App.Setup(ctx)
}

func (app *App) setupPluralization() {
	// NOTE: This should be externally configurable in the future.
	hm.AddUncountableRule("content")
}
