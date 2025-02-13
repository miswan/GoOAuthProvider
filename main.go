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
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database schema
	err = db.AutoMigrate(&models.User{}, &models.Client{}, &models.AuthCode{}, &models.RefreshToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize storage with database
	store := storage.NewPostgresStorage(db)

	// Initialize services
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	clientService := services.NewClientService(store)

	// Initialize handlers
	oauthHandler := handlers.NewOAuthHandler(oauthService)
	userHandler := handlers.NewUserHandler(userService)
	clientHandler := handlers.NewClientHandler(clientService)

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

	// Start server
	log.Fatal(e.Start(":8000"))
}