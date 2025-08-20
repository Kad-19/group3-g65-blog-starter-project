package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"g3-g65-bsp/domain"
)

func TestTokenRepository_StoreRefreshToken(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &TokenRepository{collection: mt.Coll}
		token := &domain.RefreshToken{
			UserID:    primitive.NewObjectID().Hex(),
			Token:     "refresh-token",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.StoreRefreshToken(context.Background(), token)
		assert.NoError(t, err)
	})
}

func TestTokenRepository_FindRefreshToken(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &TokenRepository{collection: mt.Coll}
		token := "refresh-token"
		expectedToken := &RefreshTokenDTO{
			UserID:    primitive.NewObjectID(),
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedToken)))

		result, err := repo.FindRefreshToken(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken.UserID.Hex(), result.UserID)
	})
}
