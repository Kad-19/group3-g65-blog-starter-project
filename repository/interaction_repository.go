package repository

import (
	"context"
	"g3-g65-bsp/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type LikesModel struct {
	Count int      `bson:"count"`
	Users []string `bson:"users"`
}

type InteractionRepository struct {
	collection *mongo.Collection
}

func NewMongoInteractionRepository(collection *mongo.Collection) domain.InteractionRepository {
	return &InteractionRepository{
		collection: collection,
	}
}

func (r *InteractionRepository) LikeBlog(ctx context.Context, userID string, blogID string, preftype string) error {
	// Implementation for liking a blog post
	return nil
}

func (r *InteractionRepository) CommentOnBlog(ctx context.Context, userID string, blogID string, comment *domain.Comment) error {
	// Implementation for commenting on a blog post
	return nil
}

