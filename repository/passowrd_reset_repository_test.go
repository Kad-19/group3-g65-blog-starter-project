package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"g3-g65-bsp/domain"
)

func TestPasswordReset_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &PasswordReset{collection: mt.Coll}
		token := &domain.PasswordResetToken{
			Email:     "test@gmail.com",
			Token:     "reset-token",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Create(context.Background(), token)
		assert.NoError(t, err)
	})
}

func TestPasswordReset_GetByToken(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := &PasswordReset{collection: mt.Coll}
		token := "reset-token"
		expectedToken := &PasswordResetTokenResponseDTO{
			Token:     token,
			Email:     "test@gmail.com",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, toBSOND(expectedToken)))

		result, err := repo.GetByToken(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken.Email, result.Email)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := &PasswordReset{collection: mt.Coll}
		token := "reset-token"

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		_, err := repo.GetByToken(context.Background(), token)
		assert.Error(t, err)
	})
}
