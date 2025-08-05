// repository/activation_token_repo.go
package repository

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivationTokenDTO struct {
	Token     string    `bson:"token,omitempty"`
	Email     string    `bson:"email,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type ActivationTokenRepo struct {
	collection *mongo.Collection
}

func NewActivationTokenRepo(db *mongo.Database) *ActivationTokenRepo {
	return &ActivationTokenRepo{
		collection: db.Collection("activation_token"),
	}
}

func (at *ActivationTokenRepo) Create(ctx context.Context, token *domain.ActivationToken) error {
	dto := toDTO(token)
	_, err := at.collection.InsertOne(ctx, dto)
	return err
}

func (at *ActivationTokenRepo) GetByToken(ctx context.Context, token string) (*domain.ActivationToken, error) {
	var dto ActivationTokenDTO
	filter := bson.M{"token": token}

	err := at.collection.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return toDomain(&dto), nil
}

func (at *ActivationTokenRepo) Delete(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	_, err := at.collection.DeleteOne(ctx, filter)
	return err
}

func (at *ActivationTokenRepo) GetByEmail(ctx context.Context, email string) (*domain.ActivationToken, error) {
	var dto ActivationTokenDTO
	filter := bson.M{"email": email}

	err := at.collection.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return toDomain(&dto), nil
}

func toDTO(token *domain.ActivationToken) *ActivationTokenDTO {
	return &ActivationTokenDTO{
		Token:     token.Token,
		Email:     token.Email,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
	}
}

func toDomain(dto *ActivationTokenDTO) *domain.ActivationToken {
	return &domain.ActivationToken{
		Token:     dto.Token,
		Email:     dto.Email,
		CreatedAt: dto.CreatedAt,
		ExpiresAt: dto.ExpiresAt,
	}
}
