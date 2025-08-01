package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository(db *mongo.Database) *TokenRepository {
	return &TokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID primitive.ObjectID, tokenHash string, expiresAt time.Time) error {
	_, err := r.collection.InsertOne(ctx, bson.M{
		"user_id":    userID,
		"token_hash": tokenHash,
		"expires_at": expiresAt,
	})
	return err
}

func (r *TokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (primitive.ObjectID, error) {
	var result struct {
		UserID primitive.ObjectID `bson:"user_id"`
	}
	err := r.collection.FindOne(ctx, bson.M{"token_hash": tokenHash}).Decode(&result)
	return result.UserID, err
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"token_hash": tokenHash})
	return err
}

func (r *TokenRepository) DeleteAllForUser(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}