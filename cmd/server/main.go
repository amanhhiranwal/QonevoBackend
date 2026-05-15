package main

import (
	"log"
	"net/http"

	"qonevo-backend/internal/config"
	"qonevo-backend/internal/controllers"
	"qonevo-backend/internal/db"
	"qonevo-backend/internal/repositories"
	v1 "qonevo-backend/internal/routes/v1"
	"qonevo-backend/internal/services"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBUrl)

	repo := repositories.NewUserRepo(database)
	service := services.NewAuthService(repo, cfg.JWTSecret, cfg.JWTExpiry)
	controller := controllers.NewAuthController(service)

	mux := http.NewServeMux()
	v1.AuthRoutes(mux, controller)

	log.Println("Server running on :", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}