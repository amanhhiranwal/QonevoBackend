package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qonevo-backend/internal/config"
	"qonevo-backend/internal/controllers"
	"qonevo-backend/internal/db"
	"qonevo-backend/internal/middleware"
	"qonevo-backend/internal/repositories"
	v1 "qonevo-backend/internal/routes/v1"
	"qonevo-backend/internal/services"
)

func main() {

	// ============================================
	// LOAD CONFIGURATION
	// ============================================

	cfg := config.Load()

	// ============================================
	// DATABASE CONNECTION
	// ============================================

	database := db.Connect(
		cfg.DatabaseURL,
	)

	defer func() {

		if err := database.Close(); err != nil {

			log.Printf(
				"failed to close database: %v",
				err,
			)
		}
	}()

	// ============================================
	// REPOSITORIES
	// ============================================

	userRepo := repositories.NewUserRepository(
		database,
	)

	productRepo := repositories.NewProductRepository(
		database,
	)

	// ============================================
	// SERVICES
	// ============================================

	authService := services.NewAuthService(
		userRepo,
		cfg,
	)

	// ============================================
	// S3 SERVICE
	// ============================================

	s3Service, err := services.NewS3Service(
		cfg.AWSS3Bucket,
		cfg.AWSRegion,
	)

	if err != nil {

		log.Fatalf(
			"failed to initialize s3 service: %v",
			err,
		)
	}

	// ============================================
	// PRODUCT SERVICE
	// ============================================

	productService := services.NewProductService(
		productRepo,
		s3Service,
	)

	// ============================================
	// CONTROLLERS
	// ============================================

	pageController := controllers.NewPageController(
		productService,
	)

	authController := controllers.NewAuthController(
		authService,
		cfg,
	)

	productController := controllers.NewProductController(
		productService,
		pageController,
		s3Service,
	)

	// ============================================
	// ROUTER
	// ============================================

	mux := http.NewServeMux()

	// ============================================
	// REGISTER ROUTES
	// ============================================

	v1.RegisterV1Routes(
		mux,
		authController,
		pageController,
		productController,
	)

	// ============================================
	// GLOBAL MIDDLEWARE
	// ============================================

	handler := middleware.Logger(
		middleware.SecurityHeaders(
			mux,
		),
	)

	// ============================================
	// SERVER CONFIG
	// ============================================

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// ============================================
	// START SERVER
	// ============================================

	go func() {

		log.Printf(
			"server running on http://localhost:%s",
			cfg.AppPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server failed: %v",
				err,
			)
		}
	}()

	// ============================================
	// GRACEFUL SHUTDOWN
	// ============================================

	quit := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println(
		"shutting down server...",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Fatalf(
			"server forced shutdown: %v",
			err,
		)
	}

	log.Println(
		"server exited properly",
	)
}
