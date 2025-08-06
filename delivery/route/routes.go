package route

import (
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

func InteractionRouter (r *gin.Engine, interactionController *controller.InteractionController, jwt *auth.JWT) {
    interactionGroup := r.Group("/blogs")
    interactionGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    {
        interactionGroup.POST("/like/:id", interactionController.LikeBlog)
        interactionGroup.POST("/comment/:id", interactionController.CommentOnBlog)
    }
}

func BlogRouter(r *gin.Engine, blogController *controller.BlogController, jwt *auth.JWT) {
    blogGroup := r.Group("/blogs")
    blogGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
    {
        blogGroup.POST("/", blogController.CreateBlog)
        blogGroup.GET("/", blogController.ListBlogs)
        blogGroup.GET(":id", blogController.GetBlogByID)
        blogGroup.PUT(":id", blogController.UpdateBlog)
        blogGroup.DELETE(":id", blogController.DeleteBlog)
    }
}

func AuthRouter(r *gin.Engine, authController *controller.AuthController, jwt *auth.JWT) {
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authController.Register)
        authGroup.POST("/login", authController.Login)
        authGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
        {
            authGroup.POST("/refresh", authController.Refresh)
            authGroup.POST("/logout", authController.Logout)      // Single device
            authGroup.POST("/logout-all", authController.LogoutAll) // All devices
        }
    }
}

// NewRouter initializes the Gin engine and registers all routes
func NewRouter() *gin.Engine {
    r := gin.Default()
    return r
}
