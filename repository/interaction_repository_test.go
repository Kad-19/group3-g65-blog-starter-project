package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"g3-g65-bsp/domain"
)

func setupInteractionTest(t *testing.T) (*InteractionRepository, *mongo.Collection, func()) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skip("skipping mongo tests, could not create mongo client: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skip("skipping mongo tests, could not connect to mongodb: " + err.Error())
	}

	collection := client.Database("testdb").Collection("blogs")
	repo := NewMongoInteractionRepository(collection)

	// Insert a blog post
	_, err = collection.InsertOne(context.Background(), bson.M{
		"_id":      "blog1",
		"authorid": "author1",
		"title":    "title1",
		"content":  "content1",
		"metrics": bson.M{
			"viewcount": 0,
			"likes":     bson.M{"count": 0, "users": []string{}},
			"dislikes":  bson.M{"count": 0, "users": []string{}},
		},
		"comments": []domain.Comment{},
	})
	assert.NoError(t, err)

	return repo.(*InteractionRepository), collection, func() {
		collection.Drop(context.Background())
		client.Disconnect(context.Background())
	}
}

func TestInteractionRepository_LikeBlog(t *testing.T) {
	repo, collection, teardown := setupInteractionTest(t)
	defer teardown()

	ctx := context.Background()
	blogID := "blog1"
	userID := "user1"

	// First like
	err := repo.LikeBlog(ctx, userID, blogID, "like")
	assert.NoError(t, err)

	var blog domain.Blog
	err = collection.FindOne(ctx, bson.M{"_id": blogID}).Decode(&blog)
	assert.NoError(t, err)
	assert.Equal(t, 1, blog.Metrics.Likes.Count)
	assert.Contains(t, blog.Metrics.Likes.Users, userID)

	// Unlike
	err = repo.LikeBlog(ctx, userID, blogID, "like")
	assert.NoError(t, err)

	err = collection.FindOne(ctx, bson.M{"_id": blogID}).Decode(&blog)
	assert.NoError(t, err)
	assert.Equal(t, 0, blog.Metrics.Likes.Count)
	assert.NotContains(t, blog.Metrics.Likes.Users, userID)
}

func TestInteractionRepository_CommentOnBlog(t *testing.T) {
	repo, collection, teardown := setupInteractionTest(t)
	defer teardown()

	ctx := context.Background()
	blogID := "blog1"
	userID := "user1"
	now := time.Now()
	comment := &domain.Comment{
		AuthorID:  userID,
		Content:   "Great post!",
		CreatedAt: &now,
	}

	err := repo.CommentOnBlog(ctx, userID, blogID, comment)
	assert.NoError(t, err)

	var blog domain.Blog
	err = collection.FindOne(ctx, bson.M{"_id": blogID}).Decode(&blog)
	assert.NoError(t, err)
	assert.Len(t, blog.Comments, 1)
	assert.Equal(t, "Great post!", blog.Comments[0].Content)
}
