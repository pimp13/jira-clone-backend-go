package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	*DB
	*App
}

func newConfig() *Config {
	return &Config{
		DB:  newDB(),
		App: newApp(),
	}
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = newConfig()
	})
	return cfg
}

var Envs = NewConfig()

type supported interface {
	~string | ~bool | ~int | ~int64 | ~float64 | ~uint | ~uint64
}

func GetEnv[T supported](key string, fallback T) T {
	if value, exists := os.LookupEnv(key); exists {
		var zero T
		var result any
		var err error

		switch any(zero).(type) {
		case string:
			return any(value).(T)
		case bool:
			result, err = strconv.ParseBool(value)
		case int:
			var v int64
			v, err = strconv.ParseInt(value, 10, 0)
			result = int(v)
		case int64:
			result, err = strconv.ParseInt(value, 10, 64)
		case float64:
			result, err = strconv.ParseFloat(value, 64)
		case time.Duration:
			result, err = time.ParseDuration(value)
		case uint:
			var v uint64
			v, err = strconv.ParseUint(value, 10, 64)
			result = uint(v)
		case uint64:
			result, err = strconv.ParseUint(value, 10, 64)
		default:
			return fallback
		}

		if err == nil {
			return result.(T)
		}
	}
	return fallback
}
