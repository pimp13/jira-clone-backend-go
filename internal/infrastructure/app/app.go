package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
)

type App struct {
	port   uint
	prefix string
	engine *echo.Echo
	// entClient *ent.Client
	cfg *config.Config
}

func NewApp() *App {
	initConfig := config.NewConfig()
	initEcho := echo.New()
	newApp := &App{
		prefix: "/api",
		cfg:    initConfig,
		port:   initConfig.App.Port,
		engine: initEcho,
	}

	return newApp
}

func (a *App) Bootstrap() error {
	addr := fmt.Sprintf(":%v", a.port)
	return a.engine.Start(addr)
}
