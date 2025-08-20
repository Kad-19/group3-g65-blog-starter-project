package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"g3-g65-bsp/domain"
)

func TestUnactiveUserRepo_CreateUnactiveUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &UnactiveUserRepo{collection: mt.Coll}
		user := &domain.UnactivatedUser{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.CreateUnactiveUser(context.Background(), user)
		assert.NoError(t, err)
	})
}

func TestUnactiveUserRepo_FindByEmailUnactive(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &UnactiveUserRepo{collection: mt.Coll}
		email := "test@example.com"
		expectedUser := &UnactivatedUserDTO{
			ID:       primitive.NewObjectID(),
			Username: "testuser",
			Email:    email,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedUser)))

		user, err := repo.FindByEmailUnactive(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
	})
}
