package route

import (
    "g3-g65-bsp/delivery/controller"
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

// NewRouter initializes the Gin engine and registers all routes
func NewRouter() *gin.Engine {
    r := gin.Default()
    return r
}
