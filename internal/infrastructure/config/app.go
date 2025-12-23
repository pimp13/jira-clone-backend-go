package config

type App struct {
	Port      uint
	Url       string
	Env       string
	SecretKey string
}

func newApp() *App {
	return &App{
		Port:      GetEnv("PORT", uint(8080)),
		Url:       GetEnv("APP_URL", "http://localhost"),
		Env:       GetEnv("APP_ENV", "development"),
		SecretKey: GetEnv("SECRET_KEY", "secret-key"),
	}
}
