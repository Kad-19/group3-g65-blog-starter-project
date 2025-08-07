package controller

import (
	"g3-g65-bsp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeRequest struct {
	Preftype string `json:"preftype" binding:"required,oneof=like dislike"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (c *CommentRequest) ConvertToDomain() *domain.Comment {
	return &domain.Comment{
		Content:        c.Content,
	}
}

type InteractionController struct {
	usecase domain.InteractionUsecase
}

func NewInteractionController(usecase domain.InteractionUsecase) *InteractionController {
	return &InteractionController{
		usecase: usecase,
	}	
}

func (c *InteractionController) LikeBlog(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	blogID := ctx.Param("id")
	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog_id is required"})
		return
	}

	var req LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := c.usecase.LikeBlog(ctx, userID, blogID, req.Preftype); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like blog"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "blog liked successfully"})
}

func (c *InteractionController) CommentOnBlog(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	blogID := ctx.Param("id")

	if blogID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "blog_id is required"})
		return
	}

	var comment CommentRequest
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := c.usecase.CommentOnBlog(ctx, userID, blogID, comment.ConvertToDomain()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to comment on blog"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment added successfully"})
}