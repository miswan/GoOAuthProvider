package main

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"log"
	"oauth2-provider/config"
	"oauth2-provider/handlers"
	"oauth2-provider/middleware"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
)

func main() {
	log.Println("Starting OAuth2 Provider application...")

	// Initialize Echo
	e := echo.New()
	log.Println("Echo framework initialized")

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	log.Println("Middleware configured successfully")

	log.Println("Attempting to connect to database...")
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Auto migrate database schema
	log.Println("Starting database migration...")
	err = db.AutoMigrate(&models.User{}, &models.Client{}, &models.AuthCode{}, &models.RefreshToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize storage with database
	store := storage.NewPostgresStorage(db)
	log.Println("PostgreSQL storage initialized")

	// Initialize services
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	clientService := services.NewClientService(store)
	log.Println("Services initialized")

	// Initialize handlers
	oauthHandler := handlers.NewOAuthHandler(oauthService)
	userHandler := handlers.NewUserHandler(userService)
	clientHandler := handlers.NewClientHandler(clientService)
	log.Println("Handlers initialized")

	// Routes
	// OAuth2 endpoints
	e.GET("/authorize", oauthHandler.Authorize)
	e.POST("/token", oauthHandler.Token)
	e.GET("/userinfo", oauthHandler.UserInfo, middleware.JWTAuth)

	// User management
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// Client management
	e.POST("/client/register", clientHandler.Register)
	e.GET("/client/:id", clientHandler.Get, middleware.JWTAuth)
	log.Println("Routes configured")

	// Start server
	log.Println("Starting server on port 8000...")
	if err := e.Start(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}