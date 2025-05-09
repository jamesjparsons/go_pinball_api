package main

import (
	"log"
	"net/http"
	"os"

	_ "backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/handlers"
	"backend/models"
)

// @title           Pinball League API
// @version         1.0
// @description     A Pinball League Management System API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := gorm.Open(sqlite.Open("pinball.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.League{}, &models.Season{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	leagueHandler := handlers.NewLeagueHandler(db)
	seasonHandler := handlers.NewSeasonHandler(db)

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(handlers.CORSMiddleware)
	router.Use(handlers.LoggingMiddleware)

	// Public routes
	router.POST("/api/auth/signup", authHandler.Signup)
	router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(handlers.AuthMiddleware)
	{
		protected.GET("/auth/me", authHandler.GetCurrentUser)

		// League routes
		protected.GET("/leagues", leagueHandler.ListLeagues)
		protected.GET("/leagues/:leagueID", leagueHandler.GetLeague)
		protected.POST("/leagues/create", leagueHandler.CreateLeague)

		// Season routes
		protected.GET("/leagues/:leagueID/seasons", seasonHandler.ListSeasons)
		protected.POST("/leagues/:leagueID/seasons", seasonHandler.CreateSeason)
	}

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
