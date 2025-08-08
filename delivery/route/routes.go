package route

import (
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/infrastructure/cache"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/middleware"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"

	"github.com/gin-gonic/gin"
)

func InteractionRouter (r *gin.Engine, interactionController *controller.InteractionController, jwt *auth.JWT) {
    interactionGroup := r.Group("/blogs")
    interactionGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    {
        interactionGroup.POST("/like/:id", interactionController.LikeBlog)
        interactionGroup.POST("/comment/:id", interactionController.CommentOnBlog)
        interactionGroup.PUT("/comment/:id/:comment_id", interactionController.UpdateComment)
        interactionGroup.DELETE("/comment/:id/:comment_id", interactionController.DeleteComment)
    }
}

func BlogRouter(r *gin.Engine, blogController *controller.BlogController, jwt *auth.JWT, cacheService *cache.Service) {
    cachingMiddleware := middleware.CachePage(*cacheService, 2*time.Minute)
    revalidateMiddleware := middleware.RevalidateCache(*cacheService)
    blogGroup := r.Group("/blogs")
    blogGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    {
        blogGroup.POST("/", blogController.CreateBlog)
        blogGroup.GET("/", cachingMiddleware, blogController.ListBlogs)
        blogGroup.GET(":id", cachingMiddleware, blogController.GetBlogByID)
        blogGroup.PUT(":id", revalidateMiddleware, blogController.UpdateBlog)
        blogGroup.DELETE(":id", revalidateMiddleware, blogController.DeleteBlog)
    }
}

func AuthRouter(r *gin.Engine, authController *controller.AuthController, jwt *auth.JWT) {
    authGroup := r.Group("/auth")
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

func OAuthRouter(r *gin.Engine, oauthController *controller.OAuthController) {
    oauthGroup := r.Group("/auth/google")
    {
        oauthGroup.GET("/login", oauthController.HandleGoogleLogin)
        oauthGroup.GET("/callback", oauthController.HandleGoogleCallback)
    }
}

func UserRouter(r *gin.Engine, userController *controller.UserController, jwt *auth.JWT){
    userGroup := r.Group("/user")
    {
        userGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
        {
            userGroup.POST("/update-profile", userController.HandleUpdateUser)
            userGroup.POST("/promote", middleware.RoleMiddleware(), userController.HandlePromote)
            userGroup.POST("/demote", middleware.RoleMiddleware(), userController.HandleDemote)
			userGroup.GET("/allusers", userController.HandleGetAllUsers)
		}
	}
}

func AIRouter(r *gin.Engine, aicontroller *controller.AIcontroller, jwt *auth.JWT) {
	aigroup := r.Group("/ai")
	{
		aigroup.Use(middleware.AuthMiddleware(jwt))
		{
			aigroup.POST("/content", aicontroller.HandleAIContentrequest)
			aigroup.POST("/enhance", aicontroller.HandleAIEnhancement)
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
    lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour, // How long each rate limiter for a given key lives.
	})
    r.Use(tollbooth_gin.LimitHandler(lmt)) // Apply rate limiting middleware
    HealthRouter(r) // Register health check endpoint
	return r
}
