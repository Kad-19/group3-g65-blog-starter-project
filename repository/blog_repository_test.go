package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"g3-g65-bsp/domain"
)

func TestMongoBlogRepository_CreateBlog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		blog := &domain.Blog{
			AuthorID:       primitive.NewObjectID().Hex(),
			AuthorUsername: "testuser",
			Title:          "Test Title",
			Content:        "Test Content",
			Tags:           []string{"go", "test"},
			Metrics: &domain.Metrics{
				ViewCount: 0,
				Likes:     &domain.Likes{Count: 0, Users: []string{}},
				Dislikes:  &domain.Likes{Count: 0, Users: []string{}},
			},
			Comments: []domain.Comment{},
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		id, err := repo.CreateBlog(context.Background(), blog)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestMongoBlogRepository_GetBlogByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		id := primitive.NewObjectID()
		expectedBlog := &BlogModel{
			ID:             id,
			AuthorID:       primitive.NewObjectID(),
			AuthorUsername: "testuser",
			Title:          "Test Title",
			Content:        "Test Content",
			Metrics: &Metrics{
				ViewCount: 0,
				Likes:     &Likes{Count: 0, Users: []string{}},
				Dislikes:  &Likes{Count: 0, Users: []string{}},
			},
			Comments: []Comment{},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedBlog)))

		blog, err := repo.GetBlogByID(context.Background(), id.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expectedBlog.Title, blog.Title)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		id := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		_, err := repo.GetBlogByID(context.Background(), id.Hex())
		assert.Error(t, err)
		assert.Equal(t, ErrBlogNotFound, err)
	})
}

func TestMongoBlogRepository_DeleteBlog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		id := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}}...))

		err := repo.DeleteBlog(context.Background(), id.Hex())
		assert.NoError(t, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		id := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 0}}...))

		err := repo.DeleteBlog(context.Background(), id.Hex())
		assert.Error(t, err)
		assert.Equal(t, ErrBlogNotFound, err)
	})
}

func TestMongoBlogRepository_AddComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		blogID := primitive.NewObjectID()
		comment := &domain.Comment{
			AuthorID: primitive.NewObjectID().Hex(),
			Content:  "New Comment",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}, {Key: "nModified", Value: 1}}...))

		err := repo.AddComment(context.Background(), blogID.Hex(), comment)
		assert.NoError(t, err)
	})
}

func TestMongoBlogRepository_GetCommentByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		blogID := primitive.NewObjectID()
		commentID := primitive.NewObjectID()
		expectedComment := &Comment{
			ID:       commentID,
			AuthorID: primitive.NewObjectID(),
			Content:  "Test Comment",
		}
		expectedBlog := &BlogModel{
			ID:       blogID,
			Comments: []Comment{*expectedComment},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedBlog)))

		comment, err := repo.GetCommentByID(context.Background(), blogID.Hex(), commentID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expectedComment.Content, comment.Content)
	})
}

func TestMongoBlogRepository_UpdateComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		blogID := primitive.NewObjectID()
		comment := &domain.Comment{
			ID:      primitive.NewObjectID().Hex(),
			Content: "Updated Comment",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}, {Key: "nModified", Value: 1}}...))

		err := repo.UpdateComment(context.Background(), blogID.Hex(), comment)
		assert.NoError(t, err)
	})
}

func TestMongoBlogRepository_DeleteComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &mongoBlogRepository{collection: mt.Coll}
		blogID := primitive.NewObjectID()
		commentID := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.D{{Key: "n", Value: 1}, {Key: "nModified", Value: 1}}...))

		err := repo.DeleteComment(context.Background(), blogID.Hex(), commentID.Hex())
		assert.NoError(t, err)
	})
}

// toBSOND is a helper function to convert a struct to a BSON document for mock responses.
func toBSOND(v any) primitive.D {
	data, err := bson.Marshal(v)
	if err != nil {
		panic(err)
	}
	var doc primitive.D
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		panic(err)
	}
	return doc
}