package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
	JWTExpiry int
}

func Load() *Config {
	godotenv.Load()

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))

	dbURL := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") +
		"?sslmode=" + os.Getenv("DB_SSLMODE")

	return &Config{
		Port:      os.Getenv("APP_PORT"),
		DBUrl:     dbURL,
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiry: exp,
	}
}