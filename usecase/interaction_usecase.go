package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InteractionUsecase struct {
	blogRepo domain.BlogRepository
	userRepo domain.UserRepository
}

func NewInteractionUsecase(blogRepo domain.BlogRepository, userRepo domain.UserRepository) domain.InteractionUsecase {
	return &InteractionUsecase{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}
}


func (u *InteractionUsecase) LikeBlog(ctx context.Context, userID string, blogID string, preftype string) error {
	// Validate preftype
	if preftype != "like" && preftype != "dislike" {
		return domain.ErrInvalidpreftype
	}
	existingBlog, err := u.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		return err
	}

	if preftype == "like" {
		if containsUser(existingBlog.Metrics.Likes.Users, userID) {
			existingBlog.Metrics.Likes.Users = removeUser(existingBlog.Metrics.Likes.Users, userID)
			existingBlog.Metrics.Likes.Count--
		} else {
			if containsUser(existingBlog.Metrics.Dislikes.Users, userID) {
				existingBlog.Metrics.Dislikes.Users = removeUser(existingBlog.Metrics.Dislikes.Users, userID)
				existingBlog.Metrics.Dislikes.Count--

			}
			existingBlog.Metrics.Likes.Users = append(existingBlog.Metrics.Likes.Users, userID)
			existingBlog.Metrics.Likes.Count++
		}
	}

	if preftype == "dislike" {
		if containsUser(existingBlog.Metrics.Dislikes.Users, userID) {
			existingBlog.Metrics.Dislikes.Users = removeUser(existingBlog.Metrics.Dislikes.Users, userID)
			existingBlog.Metrics.Dislikes.Count--
		} else {
			if containsUser(existingBlog.Metrics.Likes.Users, userID) {
				existingBlog.Metrics.Likes.Users = removeUser(existingBlog.Metrics.Likes.Users, userID)
				existingBlog.Metrics.Likes.Count--

			}
			existingBlog.Metrics.Dislikes.Users = append(existingBlog.Metrics.Dislikes.Users, userID)
			existingBlog.Metrics.Dislikes.Count++
		}

	}

	// Update the blog metrics in the repository
	if err := u.blogRepo.UpdateBlog(ctx, existingBlog); err != nil {
		return err
	}
	return nil
}

func (u *InteractionUsecase) CommentOnBlog(ctx context.Context, userID string, blogID string, comment *domain.Comment) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	existingUser, err := u.userRepo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}
	_, e := u.blogRepo.GetBlogByID(ctx, blogID)
	if e != nil {
		return e
	}
	comment.AuthorID = userID
	comment.AuthorUsername = existingUser.Username
	now := time.Now()
	comment.CreatedAt = &now
	if err := u.blogRepo.AddComment(ctx, blogID, comment); err != nil {
		return err
	}
	return nil
}

func containsUser(users []string, userID string) bool {
    return slices.Contains(users, userID)
}

func removeUser(users []string, userID string) []string {
    result := make([]string, 0, len(users))
    for _, id := range users {
        if id != userID {
            result = append(result, id)
        }
    }
    return result
}
