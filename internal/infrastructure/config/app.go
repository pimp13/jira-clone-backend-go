package config

type App struct {
	Port        uint
	Url         string
	FrontendURL string
	Env         string
	SecretKey   string
}

func newApp() *App {
	return &App{
		Port:        GetEnv("PORT", uint(8080)),
		Url:         GetEnv("APP_URL", "http://localhost"),
		FrontendURL: GetEnv("FRONTEND_URL", "http://localhost:3001"),
		Env:         GetEnv("APP_ENV", "development"),
		SecretKey:   GetEnv("SECRET_KEY", "secret-key"),
	}
}
