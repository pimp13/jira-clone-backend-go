package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/db"
	"github.com/pimp13/jira-clone-backend-go/internal/module/auth"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
	"github.com/pimp13/jira-clone-backend-go/internal/module/project"
	"github.com/pimp13/jira-clone-backend-go/internal/module/workspace"
	"github.com/pimp13/jira-clone-backend-go/pkg/logger"
)

type App struct {
	// Application running port
	port uint

	// Application api route prefix
	prefix string

	// Default api version
	version string

	// Logger use by zerolog
	logger logger.Logger

	// Engine for running routes
	engine *echo.Echo

	// Database orm for connect to database
	entClient *ent.Client

	// Configuration and ENV global variables
	cfg *config.Config

	// Is production for check status running server
	isProduction bool
}

func NewApp() (*App, error) {
	cfg := config.NewConfig()
	isProduction := cfg.App.Env == "production"

	logger := logger.New(isProduction)

	entClient, err := db.NewEntClient(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	e := echo.New()

	logger.Info().
		Bool("IsProductionMode", isProduction).
		Msg("Application initialized!")

	return &App{
		port:         cfg.App.Port,
		prefix:       "/api",
		version:      "v1",
		logger:       logger,
		engine:       e,
		entClient:    entClient,
		cfg:          cfg,
		isProduction: isProduction,
	}, nil
}

func (a *App) Bootstrap() error {
	a.setupMiddlewares()
	a.setupRoutes()
	a.setupServices()

	return nil
}

func (a *App) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", a.port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	server := &http.Server{
		Addr:    addr,
		Handler: a.engine,
	}

	go func() {
		a.logger.Info().Str("address", addr).Msg("🚀 Server starting on")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("server failed!")
		}
	}()

	<-quit
	a.logger.Warn().Msg("🛑 Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Error().Err(err).Msg("graceful shutdown error!")
		return err
	}

	a.logger.Info().Msg("✅ Server gracefully stopped")
	return nil
}

func (a *App) setupMiddlewares() {
	a.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			a.cfg.App.FrontendURL,
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"Authorization",
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	if !a.isProduction {
		a.engine.Use(middleware.Logger())
	}
	a.engine.Use(middleware.Recover())

	a.engine.Static("/public", "public")
}

func (a *App) setupRoutes() {
	// Swagger
	a.engine.GET("/api/docs/*", echoSwagger.WrapHandler)
	a.engine.GET("/api/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/api/docs/index.html")
	})
}

func (a *App) setupServices() {
	// V1 Routes
	api_v1 := a.engine.Group(fmt.Sprintf("%s/%s", a.prefix, a.version))

	// JWT & Auth
	jwtSvc := jwt.NewJWTService(a.entClient, a.cfg)
	authMiddleware := auth.NewAuthMiddleware(jwtSvc)
	authSvc := auth.NewAuthService(a.entClient, jwtSvc, a.logger)
	authCtrl := auth.NewAuthController(authSvc, authMiddleware)
	authCtrl.Routes(api_v1)

	// Workspace
	wsSvc := workspace.NewWorkspaceService(a.entClient, a.logger)
	wsCtrl := workspace.NewWorkspaceController(wsSvc, authMiddleware)
	wsCtrl.Routes(api_v1)

	// Project
	projectService := project.NewProjectService(a.entClient, a.logger)
	projectController := project.NewProjectController(projectService, authMiddleware)
	projectController.Routes(api_v1)

}
