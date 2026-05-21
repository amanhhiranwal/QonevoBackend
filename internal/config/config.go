package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {

	// ============================================
	// APPLICATION
	// ============================================

	AppPort string
	AppEnv  string
	AppURL  string

	// ============================================
	// DATABASE
	// ============================================

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	DatabaseURL string

	// ============================================
	// JWT AUTH
	// ============================================

	JWTSecret      string
	JWTExpiryHours int

	// ============================================
	// AWS S3
	// ============================================

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Bucket        string

	// ============================================
	// FILE UPLOAD
	// ============================================

	MaxUploadSize int64
}

// ============================================
// LOAD CONFIG
// ============================================

func Load() *Config {

	// ============================================
	// LOAD ENV
	// ============================================

	err := godotenv.Load()

	if err != nil {
		log.Println("warning: .env file not found")
	}

	// ============================================
	// JWT EXPIRY
	// ============================================

	jwtExpiry, err := strconv.Atoi(
		getEnv("JWT_EXPIRY_HOURS", "24"),
	)

	if err != nil {
		jwtExpiry = 24
	}

	// ============================================
	// MAX FILE SIZE
	// ============================================

	maxUploadSize, err := strconv.ParseInt(
		getEnv("MAX_UPLOAD_SIZE", "10485760"),
		10,
		64,
	)

	if err != nil {
		maxUploadSize = 10485760
	}

	// ============================================
	// DATABASE URL
	// ============================================

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "qonevo-backend"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// ============================================
	// RETURN CONFIG
	// ============================================

	return &Config{

		// ========================================
		// APP
		// ========================================

		AppPort: getEnv(
			"APP_PORT",
			"8080",
		),

		AppEnv: getEnv(
			"APP_ENV",
			"development",
		),

		AppURL: getEnv(
			"APP_URL",
			"http://localhost:8080",
		),

		// ========================================
		// DATABASE
		// ========================================

		DBHost: getEnv(
			"DB_HOST",
			"localhost",
		),

		DBPort: getEnv(
			"DB_PORT",
			"5432",
		),

		DBUser: getEnv(
			"DB_USER",
			"postgres",
		),

		DBPassword: getEnv(
			"DB_PASSWORD",
			"",
		),

		DBName: getEnv(
			"DB_NAME",
			"qonevo-backend",
		),

		DBSSLMode: getEnv(
			"DB_SSLMODE",
			"disable",
		),

		DatabaseURL: dbURL,

		// ========================================
		// JWT
		// ========================================

		JWTSecret: getEnv(
			"JWT_SECRET",
			"change-me",
		),

		JWTExpiryHours: jwtExpiry,

		// ========================================
		// AWS
		// ========================================

		AWSAccessKeyID: getEnv(
			"AWS_ACCESS_KEY_ID",
			"",
		),

		AWSSecretAccessKey: getEnv(
			"AWS_SECRET_ACCESS_KEY",
			"",
		),

		AWSRegion: getEnv(
			"AWS_REGION",
			"",
		),

		AWSS3Bucket: getEnv(
			"AWS_S3_BUCKET",
			"",
		),

		// ========================================
		// FILES
		// ========================================

		MaxUploadSize: maxUploadSize,
	}
}

// ============================================
// ENV HELPER
// ============================================

func getEnv(
	key string,
	fallback string,
) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}