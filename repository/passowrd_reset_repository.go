package repository

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordResetTokenResponseDTO struct {
	Token     string    `bson:"token,omitempty"`
	Email     string    `bson:"user_id,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

type PasswordReset struct {
	collection *mongo.Collection
}

func NewPasswordReset(mc *mongo.Collection) *PasswordReset {
	return &PasswordReset{
		collection: mc,
	}
}

func (pr *PasswordReset) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	_, err := pr.collection.InsertOne(ctx, toDTOPass(token))
	return err
}

func (pr *PasswordReset) GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	var dto PasswordResetTokenResponseDTO
	filter := bson.M{"token": token}
	if err := pr.collection.FindOne(ctx, filter).Decode(&dto); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return toDomainPass(&dto), nil
}

func (pr *PasswordReset) Delete(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	_, err := pr.collection.DeleteOne(ctx, filter)
	return err
}

func toDTOPass(token *domain.PasswordResetToken) *PasswordResetTokenResponseDTO {
	return &PasswordResetTokenResponseDTO{
		Token:     token.Token,
		Email:     token.Email,
		ExpiresAt: token.ExpiresAt,
	}
}

func toDomainPass(dto *PasswordResetTokenResponseDTO) *domain.PasswordResetToken {
	return &domain.PasswordResetToken{
		Token:     dto.Token,
		Email:     dto.Email,
		ExpiresAt: dto.ExpiresAt,
	}
}
