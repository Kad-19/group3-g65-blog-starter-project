package controller

import (
	"fmt"
	"g3-g65-bsp/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BlogDTO is a data transfer object for Blog with JSON restrictions
type BlogDTO struct {
	ID             string       `json:"id,omitempty"`
	AuthorID       string       `json:"author_id"`
	AuthorUsername string       `json:"author_username"`
	Title          string       `json:"title" binding:"required"`
	Content        string       `json:"content" binding:"required"`
	Tags           []string     `json:"tags"`
	Metrics        *MetricsDTO  `json:"metrics"`
	Comments       []CommentDTO `json:"comments"`
	CreatedAt      *time.Time   `json:"created_at,omitempty"`
	UpdatedAt      *time.Time   `json:"updated_at,omitempty"`
}

type MetricsDTO struct {
	ViewCount int       `json:"view_count"`
	Likes     *LikesDTO `json:"likes"`
	Dislikes  *LikesDTO `json:"dislikes"`
}

type LikesDTO struct {
	Count int      `json:"count"`
	Users []string `json:"users"`
}

type CommentDTO struct {
	ID             string     `json:"id,omitempty"`
	AuthorID       string     `json:"author_id"`
	AuthorUsername string     `json:"author_username"`
	Content        string     `json:"content"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

// ConvertToDomain converts BlogDTO to domain.Blog
func (dto *BlogDTO) ConvertToDomain() *domain.Blog {
	return &domain.Blog{
		Title:   dto.Title,
		Content: dto.Content,
		Tags:    dto.Tags,
	}
}

// ConvertFromDomain converts a domain.Blog to BlogDTO
func ConvertFromDomain(blog *domain.Blog) *BlogDTO {
	comments := make([]CommentDTO, len(blog.Comments))
	for i, c := range blog.Comments {
		createdAt := c.CreatedAt
		comments[i] = CommentDTO{
			ID:             c.ID,
			AuthorID:       c.AuthorID,
			AuthorUsername: c.AuthorUsername,
			Content:        c.Content,
			CreatedAt:      createdAt,
		}
	}
	createdAt := blog.CreatedAt
	updatedAt := blog.UpdatedAt
	return &BlogDTO{
		ID:             blog.ID,
		AuthorID:       blog.AuthorID,
		AuthorUsername: blog.AuthorUsername,
		Title:          blog.Title,
		Content:        blog.Content,
		Tags:           blog.Tags,
		Metrics: &MetricsDTO{
			ViewCount: blog.Metrics.ViewCount,
			Likes: &LikesDTO{
				Count: blog.Metrics.Likes.Count,
				Users: blog.Metrics.Likes.Users,
			},
			Dislikes: &LikesDTO{
				Count: blog.Metrics.Dislikes.Count,
				Users: blog.Metrics.Dislikes.Users,
			},
		},
		Comments:  comments,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type BlogController struct {
	blogUsecase domain.BlogUsecase
}

func NewBlogController(blogUsecase domain.BlogUsecase) *BlogController {
	return &BlogController{blogUsecase: blogUsecase}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	userid := ctx.GetString("user_id")
	if userid == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
	}
	var blog BlogDTO
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	newBlog, err := c.blogUsecase.CreateBlog(ctx, blog.ConvertToDomain(), userid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ConvertFromDomain(newBlog))
}

func (c *BlogController) GetBlogByID(ctx *gin.Context) {
	id := ctx.Param("id")
	blog, err := c.blogUsecase.GetBlogByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}
	ctx.JSON(http.StatusOK, ConvertFromDomain(blog))
}

func (c *BlogController) UpdateBlog(ctx *gin.Context) {
	userid, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	id := ctx.Param("id")

	var blog BlogDTO
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	updatedBlog, err := c.blogUsecase.UpdateBlog(ctx, blog.ConvertToDomain(), userid.(string), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ConvertFromDomain(updatedBlog))
}

func (c *BlogController) DeleteBlog(ctx *gin.Context) {
	userid, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	role, ok := ctx.Get("role")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	id := ctx.Param("id")
	err := c.blogUsecase.DeleteBlog(ctx, id, userid.(string), role.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (c *BlogController) ListBlogs(ctx *gin.Context) {
	filter := make(map[string]any)

	// Parse tags (comma-separated)
	if tagsStr := ctx.Query("tags"); tagsStr != "" {
		tags := []string{}
		for _, t := range ctx.QueryArray("tags") {
			for _, tag := range splitAndTrim(t, ",") {
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}
		if len(tags) > 0 {
			filter["tags"] = tags
		}
	}

	// Parse date range
	if from := ctx.Query("created_at_from"); from != "" {
		filter["created_at_from"] = from
	}
	if to := ctx.Query("created_at_to"); to != "" {
		filter["created_at_to"] = to
	}

	// Parse min_views
	if minViews := ctx.Query("min_views"); minViews != "" {
		if mv, err := parseInt(minViews); err == nil {
			filter["min_views"] = mv
		}
	}

	// Parse search query
	if search := ctx.Query("search"); search != "" {
		filter["search"] = search
	}

	// Parse pagination
	page, limit := 1, 10
	if p := ctx.Query("page"); p != "" {
		if v, err := parseInt(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := ctx.Query("limit"); l != "" {
		if v, err := parseInt(l); err == nil && v > 0 {
			limit = v
		}
	}

	// Parse sorting
	if sort := ctx.Query("sortBy"); sort != "" {
		filter["sortBy"] = sort
	}

	if order := ctx.Query("order"); order != "" {
		filter["order"] = order
	}

	blogs, pagination, err := c.blogUsecase.ListBlogs(ctx, filter, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dtos := make([]*BlogDTO, len(blogs))
	for i, b := range blogs {
		dtos[i] = ConvertFromDomain(b)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":       dtos,
		"pagination": pagination,
	})
}

// splitAndTrim splits a string by sep and trims spaces
func splitAndTrim(s, sep string) []string {
	arr := make([]string, 0)
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			arr = append(arr, trimmed)
		}
	}
	return arr
}

func split(s, sep string) []string {
	var res []string
	i := 0
	for i < len(s) {
		j := i
		for j < len(s) && string(s[j]) != sep {
			j++
		}
		res = append(res, s[i:j])
		i = j + 1
	}
	return res
}

func trim(s string) string {
	i, j := 0, len(s)-1
	for i <= j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
		i++
	}
	for j >= i && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n') {
		j--
	}
	if i > j {
		return ""
	}
	return s[i : j+1]
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
