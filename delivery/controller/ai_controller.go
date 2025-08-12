package controller

import (
	"g3-g65-bsp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentRequest struct {
	Title string `json:"title" binding:"required"`
}

type EnhanceContent struct {
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
}

type AIcontroller struct {
	aiusecase domain.AIUseCase
}

func NewAIcontroller(aus domain.AIUseCase) *AIcontroller {
	return &AIcontroller{
		aiusecase: aus,
	}
}

func (ac *AIcontroller) HandleAIContentrequest(c *gin.Context) {
	var content ContentRequest
	if err := c.ShouldBindJSON(&content); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "incorrect request"})
	}

	res, err := ac.aiusecase.GenerateIntialSuggestion(c.Request.Context(), content.Title)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to generate content"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"content": res})
}

func (ac *AIcontroller) HandleAIEnhancement(c *gin.Context) {
	var content EnhanceContent
	if err := c.ShouldBindJSON(&content); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "incorrect request"})
		return
	}

	res, err := ac.aiusecase.GenerateBasedOnTags(c.Request.Context(), content.Content, content.Tags)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"content": res})
}
