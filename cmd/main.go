package main

import (
	"log"

	"github.com/SPVJ/fs-common-lib/core/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thanavatC/auth-service-go/config"
	"github.com/thanavatC/auth-service-go/controller"
	"github.com/thanavatC/auth-service-go/model"
	"github.com/thanavatC/auth-service-go/repository"
	"github.com/thanavatC/auth-service-go/router"
	"github.com/thanavatC/auth-service-go/service"
)

func main() {
	// Initialize configuration
	if err := config.LoadConfig(""); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	db := db.New(config.AppConfig.Database)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if db == nil {
		log.Fatal("Failed to initialize database connection")
	}

	// Initialize components
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, viper.GetString("jwt.secret"))
	authController := controller.NewAuthController(authService)

	// Initialize router
	r := gin.Default()

	// Setup routes
	router.AuthRouter(authController, r)

	// Start server
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
