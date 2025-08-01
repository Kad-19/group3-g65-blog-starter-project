package controller

import (
	"g3-g65-bsp/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BlogDTO is a data transfer object for Blog with JSON restrictions
type BlogDTO struct {
    ID        string        `json:"id,omitempty"`
    AuthorID  string        `json:"author_id" binding:"required"`
    Title     string        `json:"title" binding:"required"`
    Content   string        `json:"content" binding:"required"`
    Tags      []string      `json:"tags"`
    Metrics   MetricsDTO    `json:"metrics"`
    Comments  []CommentDTO  `json:"comments"`
    CreatedAt *time.Time    `json:"created_at,omitempty"`
    UpdatedAt *time.Time    `json:"updated_at,omitempty"`
}

type MetricsDTO struct {
    ViewCount int      `json:"view_count"`
    Likes     LikesDTO `json:"likes"`
}

type LikesDTO struct {
    Count int      `json:"count"`
    Users []string `json:"users"`
}

type CommentDTO struct {
    ID             string     `json:"id,omitempty"`
    AuthorID       string     `json:"author_id" binding:"required"`
    AuthorUsername string     `json:"author_username"`
    Content        string     `json:"content" binding:"required"`
    CreatedAt      *time.Time `json:"created_at,omitempty"`
}


// ConvertToDomain converts BlogDTO to domain.Blog
func (dto *BlogDTO) ConvertToDomain() *domain.Blog {
    comments := make([]domain.Comment, len(dto.Comments))
    for i, c := range dto.Comments {
        var createdAt time.Time
        if c.CreatedAt != nil {
            createdAt = *c.CreatedAt
        }
        comments[i] = domain.Comment{
            ID:             c.ID,
            AuthorID:       c.AuthorID,
            AuthorUsername: c.AuthorUsername,
            Content:        c.Content,
            CreatedAt:      createdAt,
        }
    }
    var createdAt, updatedAt time.Time
    if dto.CreatedAt != nil {
        createdAt = *dto.CreatedAt
    }
    if dto.UpdatedAt != nil {
        updatedAt = *dto.UpdatedAt
    }
    return &domain.Blog{
        ID:        dto.ID,
        AuthorID:  dto.AuthorID,
        Title:     dto.Title,
        Content:   dto.Content,
        Tags:      dto.Tags,
        Metrics: domain.Metrics{
            ViewCount: dto.Metrics.ViewCount,
            Likes: domain.Likes{
                Count: dto.Metrics.Likes.Count,
                Users: dto.Metrics.Likes.Users,
            },
        },
        Comments:  comments,
        CreatedAt: createdAt,
        UpdatedAt: updatedAt,
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
            CreatedAt:      &createdAt,
        }
    }
    createdAt := blog.CreatedAt
    updatedAt := blog.UpdatedAt
    return &BlogDTO{
        ID:        blog.ID,
        AuthorID:  blog.AuthorID,
        Title:     blog.Title,
        Content:   blog.Content,
        Tags:      blog.Tags,
        Metrics: MetricsDTO{
            ViewCount: blog.Metrics.ViewCount,
            Likes: LikesDTO{
                Count: blog.Metrics.Likes.Count,
                Users: blog.Metrics.Likes.Users,
            },
        },
        Comments:  comments,
        CreatedAt: &createdAt,
        UpdatedAt: &updatedAt,
    }
}

type BlogController struct {
    blogUsecase domain.BlogUsecase
}

func NewBlogController(blogUsecase domain.BlogUsecase) *BlogController {
    return &BlogController{blogUsecase: blogUsecase}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
    var blog BlogDTO
    if err := ctx.ShouldBindJSON(&blog); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }
    id, err := c.blogUsecase.CreateBlog(ctx, blog.ConvertToDomain())
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    blog.ID = id
    ctx.JSON(http.StatusCreated, blog)
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
    id := ctx.Param("id")
    var blog BlogDTO
    if err := ctx.ShouldBindJSON(&blog); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }
    blog.ID = id
    err := c.blogUsecase.UpdateBlog(ctx, blog.ConvertToDomain())
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
        return
    }
    ctx.JSON(http.StatusOK, blog)
}

func (c *BlogController) DeleteBlog(ctx *gin.Context) {
    id := ctx.Param("id")
    err := c.blogUsecase.DeleteBlog(ctx, id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
        return
    }
    ctx.Status(http.StatusNoContent)
}

func (c *BlogController) ListBlogs(ctx *gin.Context) {
    filter := make(map[string]interface{})
    // Add any filtering logic here if needed
    blogs, err := c.blogUsecase.ListBlogs(ctx, filter)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    dtos := make([]*BlogDTO, len(blogs))
    for i, b := range blogs {
        dtos[i] = ConvertFromDomain(b)
    }
    ctx.JSON(http.StatusOK, dtos)
}

