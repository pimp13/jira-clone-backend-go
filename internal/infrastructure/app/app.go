package app

import (
	"context"
	"fmt"
	"log"
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
	"github.com/pimp13/jira-clone-backend-go/internal/module/workspace"
)

type App struct {
	// Application running port
	port uint
	// Application api route prefix
	prefix string
	// Default api version
	version string
	// Engine for running routes
	engine *echo.Echo
	// Database orm for connect to database
	entClient *ent.Client
	// Configuration and ENV global variables
	cfg *config.Config
}

func NewApp() (*App, error) {
	cfg := config.NewConfig()

	entClient, err := db.NewEntClient(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	e := echo.New()

	return &App{
		port:      cfg.App.Port,
		prefix:    "/api",
		version:   "v1",
		engine:    e,
		entClient: entClient,
		cfg:       cfg,
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
		log.Printf("🚀 Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
		return err
	}

	log.Println("✅ Server gracefully stopped")
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

	a.engine.Use(middleware.Logger())
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
	api := a.engine.Group(fmt.Sprintf("%s/%s", a.prefix, a.version))

	// JWT & Auth
	jwtSvc := jwt.NewJWTService(a.entClient, a.cfg)
	authMiddleware := auth.NewAuthMiddleware(jwtSvc)

	authSvc := auth.NewAuthService(a.entClient, jwtSvc)
	authCtrl := auth.NewAuthController(authSvc, authMiddleware)
	authCtrl.Routes(api)

	// Workspace
	wsSvc := workspace.NewWorkspaceService(a.entClient)
	wsCtrl := workspace.NewWorkspaceController(wsSvc, authMiddleware)
	wsCtrl.Routes(api)
}
