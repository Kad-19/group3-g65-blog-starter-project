package route

import (
	"g3-g65-bsp/delivery/controller"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

func BlogRouter(r *gin.Engine, blogController *controller.BlogController) {
    blogGroup := r.Group("/blogs")
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
        authGroup.POST("/activate", authController.ActivateUser)
        authGroup.POST("/resend-activation", authController.ResendActivationEmail)
        authGroup.POST("/forgot-password", authController.ForgotPassword)
        authGroup.POST("/reset-password", authController.ResetPassword)
        authGroup.Use(middleware.AuthMiddleware(jwt)) // Apply auth middleware
        {
            authGroup.POST("/refresh", authController.Refresh)
            authGroup.POST("/logout", authController.Logout)      // Single device
            authGroup.POST("/logout-all", authController.LogoutAll) // All devices
        }
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
        }
    }
}

// NewRouter initializes the Gin engine and registers all routes
func NewRouter() *gin.Engine {
	r := gin.Default()
	return r
}