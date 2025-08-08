package main

import (
	"g3-g65-bsp/config"
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/delivery/route"
	"g3-g65-bsp/infrastructure"
	"g3-g65-bsp/infrastructure/ai"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/cache"
	"g3-g65-bsp/infrastructure/database"
	"g3-g65-bsp/infrastructure/email"
	"g3-g65-bsp/infrastructure/image"
	"g3-g65-bsp/repository"
	"g3-g65-bsp/usecase"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
)

func main() {
	// Initialize configuration
	infrastructure.InitLogger()
	config.LoadConfig()
	dbName := config.AppConfig.DbName
	accessSecret := config.AppConfig.AccessTokenSecret
	refreshSecret := config.AppConfig.RefreshTokenSecret
	accessExpiry := config.AppConfig.AccessTokenExpiry
	refreshExpiry := config.AppConfig.RefreshTokenExpiry

	// Initialize MongoDB connection
	db := database.InitMongoDB().Database(dbName)
	blogCollection := db.Collection("blogs")

	// Initialize repository, usecase, controller for authentication
	authRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	unActiveUserRepo := repository.NewUnactiveUserRepo(db)
	passwordResetRepo := repository.NewPasswordReset(db)
	emailService := email.NewEmailService()
	jwt := auth.NewJWT(accessSecret, refreshSecret, accessExpiry, refreshExpiry)
	authUsecase := usecase.NewAuthUsecase(authRepo, tokenRepo, jwt, unActiveUserRepo, emailService, passwordResetRepo)
	authController := controller.NewAuthController(authUsecase, jwt)

	// Initialize repository, usecase, controller for blogs
	repoCacheService := cache.NewInMemoryCache(5*time.Minute, 10*time.Minute)
	blogRepo := repository.NewBlogRepository(blogCollection, repoCacheService)
	blogUsecase := usecase.NewBlogUsecase(blogRepo, authRepo)
	blogController := controller.NewBlogController(blogUsecase)

	// Initialize interaction usecase and controller
	interactionUsecase := usecase.NewInteractionUsecase(blogRepo, authRepo)
	interactionController := controller.NewInteractionController(interactionUsecase)

	// Initialize OAuth usecase and controller
	oauthUsecase := usecase.NewOAuthUsecase(authRepo, tokenRepo, jwt)
	oauthController := controller.NewOAuthController(oauthUsecase)

	// Initialize repository, usecase, controller for user management
	imageUpload := image.NewCloudinaryService()
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, imageUpload)
	userController := controller.NewUserController(userUsecase)

	// Initialize Cache
	cacheService := cache.NewInMemoryCache(5*time.Minute, 10*time.Minute)

	aiservice := ai.NewGeminiService()
	aiusecase := usecase.NewAIUsecaseImpl(aiservice)
	aicontroller := controller.NewAIcontroller(aiusecase)

    // Initialize router
    r := route.NewRouter()
	contentCreationLimiter := tollbooth.NewLimiter(0.5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	contentReadLimiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	route.BlogRouter(r, blogController, jwt, &cacheService, contentCreationLimiter, contentReadLimiter)
	route.InteractionRouter(r, interactionController, jwt, contentCreationLimiter)

	// Register authentication routes
	authLimiter := tollbooth.NewLimiter(0.16, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Minute})
	route.AuthRouter(r, authController, jwt, authLimiter)

	// Register OAuth routes
	route.OAuthRouter(r, oauthController, authLimiter)

	// user management routes
	route.UserRouter(r, userController, jwt, contentCreationLimiter, contentReadLimiter)

	//ai features routes
	route.AIRouter(r, aicontroller, jwt, contentCreationLimiter)

	// Start the server on port 8080
	if err := r.Run("localhost:8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
