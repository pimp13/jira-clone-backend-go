package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/db"
	"github.com/pimp13/jira-clone-backend-go/internal/module/auth"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type App struct {
	port      uint
	prefix    string
	engine    *echo.Echo
	entClient *ent.Client
	cfg       *config.Config
}

func NewApp() *App {
	initConfig := config.NewConfig()
	initEcho := echo.New()
	initDb, err := db.NewEntClient(initConfig.DB)
	if err != nil {
		log.Fatalf("error in connect to db: %v", err)
	}
	newApp := &App{
		prefix:    "/api",
		cfg:       initConfig,
		port:      initConfig.App.Port,
		engine:    initEcho,
		entClient: initDb,
	}

	return newApp
}

func (a *App) Bootstrap() error {
	a.setupMiddlewares()

	a.engine.GET("/api/docs/*", echoSwagger.WrapHandler)
	a.engine.GET("/api/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusOK, "/api/docs/index.html")
	})

	api_v1 := a.engine.Group(a.prefix + "/v1")

	// Register services
	// AUTH SERVICES
	jwtService := jwt.NewJWTService(a.entClient, a.cfg)
	authMiddleware := auth.NewAuthMiddleware(jwtService)
	authService := auth.NewAuthService(a.entClient, jwtService)
	authController := auth.NewAuthController(authService, authMiddleware)
	authController.Routes(api_v1)

	addr := fmt.Sprintf(":%v", a.port)
	return a.engine.Start(addr)
}

func (a *App) setupMiddlewares() {

	a.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			a.cfg.App.FrontendURL,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"Authorization",
			"X-Requested-With",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowCredentials: true,
	}))

	a.engine.Use(middleware.Logger())

	a.engine.Use(middleware.Recover())

	a.engine.Use(middleware.Static("/public"))

	a.engine.Static("/public", "public")

}
