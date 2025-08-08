package route

import (
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/infrastructure/cache"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/middleware"
	"time"

	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"

	"github.com/gin-gonic/gin"
)

func InteractionRouter (r *gin.Engine, interactionController *controller.InteractionController, jwt *auth.JWT, contentCreationLimiter *limiter.Limiter) {
    interactionGroup := r.Group("/blogs")
    interactionGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    interactionGroup.Use(tollbooth_gin.LimitHandler(contentCreationLimiter)) // Apply rate limiting middleware
    {
        interactionGroup.POST("/like/:id", interactionController.LikeBlog)
        interactionGroup.POST("/comment/:id", interactionController.CommentOnBlog)
        interactionGroup.PUT("/comment/:id/:comment_id", interactionController.UpdateComment)
        interactionGroup.DELETE("/comment/:id/:comment_id", interactionController.DeleteComment)
    }
}

func BlogRouter(r *gin.Engine, blogController *controller.BlogController, jwt *auth.JWT, cacheService *cache.Service, contentCreationLimiter *limiter.Limiter, contentReadLimiter *limiter.Limiter) {
    cachingMiddleware := middleware.CachePage(*cacheService, 2*time.Minute)
    revalidateMiddleware := middleware.RevalidateCache(*cacheService)
    blogGroup := r.Group("/blogs")
    blogGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    {
        blogGroup.POST("/", tollbooth_gin.LimitHandler(contentCreationLimiter), blogController.CreateBlog)
        blogGroup.GET("/", tollbooth_gin.LimitHandler(contentReadLimiter), cachingMiddleware, blogController.ListBlogs)
        blogGroup.GET(":id", tollbooth_gin.LimitHandler(contentReadLimiter), blogController.GetBlogByID)
        blogGroup.PUT(":id", tollbooth_gin.LimitHandler(contentCreationLimiter), revalidateMiddleware, blogController.UpdateBlog)
        blogGroup.DELETE(":id", tollbooth_gin.LimitHandler(contentCreationLimiter), revalidateMiddleware, blogController.DeleteBlog)
    }
}

func AuthRouter(r *gin.Engine, authController *controller.AuthController, jwt *auth.JWT, authLimiter *limiter.Limiter) {
    authGroup := r.Group("/auth")
    authGroup.Use(tollbooth_gin.LimitHandler(authLimiter)) // Apply rate limiting middleware
    {
        authGroup.POST("/register", authController.Register)
        authGroup.POST("/login", authController.Login)
        authGroup.GET("/activate", authController.ActivateUser)
        authGroup.POST("/resend-activation", authController.ResendActivationEmail)
        authGroup.POST("/forgot-password", authController.ForgotPassword)
        authGroup.POST("/reset-password", authController.ResetPassword)
        authGroup.POST("/refresh_token", authController.RefreshAccessToken)
        
        authGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
        {
            authGroup.POST("/logout", authController.Logout)      // Single device
            authGroup.POST("/logout-all", authController.LogoutAll) // All devices
        }
    }
}

func OAuthRouter(r *gin.Engine, oauthController *controller.OAuthController, authLimiter *limiter.Limiter) {
    oauthGroup := r.Group("/auth/google")
    oauthGroup.Use(tollbooth_gin.LimitHandler(authLimiter)) // Apply rate limiting middleware
    {
        oauthGroup.GET("/login", oauthController.HandleGoogleLogin)
        oauthGroup.GET("/callback", oauthController.HandleGoogleCallback)
    }
}

func UserRouter(r *gin.Engine, userController *controller.UserController, jwt *auth.JWT, contentCreationLimiter *limiter.Limiter, contentReadLimiter *limiter.Limiter) {
    userGroup := r.Group("/user")
    {
        userGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
        {
            userGroup.POST("/update-profile", tollbooth_gin.LimitHandler(contentCreationLimiter), userController.HandleUpdateUser)
            userGroup.POST("/promote", tollbooth_gin.LimitHandler(contentCreationLimiter), middleware.RoleMiddleware(), userController.HandlePromote)
            userGroup.POST("/demote", tollbooth_gin.LimitHandler(contentCreationLimiter), middleware.RoleMiddleware(), userController.HandleDemote)
			userGroup.GET("/allusers", tollbooth_gin.LimitHandler(contentReadLimiter), userController.HandleGetAllUsers)
		}
	}
}

func AIRouter(r *gin.Engine, aicontroller *controller.AIcontroller, jwt *auth.JWT, contentCreationLimiter *limiter.Limiter) {
	aigroup := r.Group("/ai")
	{
		aigroup.Use(middleware.AuthMiddleware(jwt))
		{
			aigroup.POST("/content", tollbooth_gin.LimitHandler(contentCreationLimiter), aicontroller.HandleAIContentrequest)
			aigroup.POST("/enhance", tollbooth_gin.LimitHandler(contentCreationLimiter), aicontroller.HandleAIEnhancement)
        }
    }
}

// HealthRouter registers a health check endpoint
func HealthRouter(r *gin.Engine) {
    r.GET("/health", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{"status": "ok"})
    })
}


// NewRouter initializes the Gin engine and registers all routes
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("utils/*.html")
    HealthRouter(r) // Register health check endpoint
	return r
}
