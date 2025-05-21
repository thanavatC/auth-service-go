package main

import (
	"fmt"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Initialize database
	db := db.New(config.AppConfig.Database)

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

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
