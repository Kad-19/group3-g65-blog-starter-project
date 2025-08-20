package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"g3-g65-bsp/domain"
)

func TestUserRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &UserRepository{collection: mt.Coll}
		user := &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
			Role:     "user",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &UserRepository{collection: mt.Coll}
		email := "test@example.com"
		expectedUser := &UserDTO{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Email:    email,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedUser)))

		user, err := repo.FindByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := &UserRepository{collection: mt.Coll}
		email := "test@example.com"

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		_, err := repo.FindByEmail(context.Background(), email)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &UserRepository{collection: mt.Coll}
		id := primitive.NewObjectID()
		expectedUser := &UserDTO{
			ID:       id,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedUser)))

		user, err := repo.FindByID(context.Background(), id.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := &UserRepository{collection: mt.Coll}
		id := primitive.NewObjectID()

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		_, err := repo.FindByID(context.Background(), id.Hex())
		assert.Error(t, err)
	})
}
