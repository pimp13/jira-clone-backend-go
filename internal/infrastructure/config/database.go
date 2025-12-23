package config

type DB struct {
	Connection string
	Name       string
	Pass       string
	User       string
	Port       uint
	Host       string
}

func newDB() *DB {
	return &DB{
		Connection: GetEnv("DB_CONNECTION", "postgres"),
		Name:       GetEnv("DB_NAME", "dbname"),
		Pass:       GetEnv("DB_PASS", "pass"),
		User:       GetEnv("DB_USER", "username"),
		Port:       GetEnv("DB_PORT", uint(5432)),
		Host:       GetEnv("DB_HOST", "127.0.0.1"),
	}
}
