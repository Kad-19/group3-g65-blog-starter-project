package repository

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type UnactiveUserRepo struct {
	collection *mongo.Collection
}

func NewUnactiveUserRepo(db *mongo.Database) domain.UnactiveUserRepo {
	return &UnactiveUserRepo{
		collection: db.Collection("unactivated_users"),
	}
}

func (at *UnactiveUserRepo) CreateUnactiveUser(ctx context.Context, user *domain.UnactivatedUser) error {
	_, err := at.collection.InsertOne(ctx, user)
	return err
}

func (at *UnactiveUserRepo) FindByEmailUnactive(ctx context.Context, email string) (*domain.UnactivatedUser, error) {
	var user domain.UnactivatedUser
	filter := bson.M{"email": email}
	err := at.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (at *UnactiveUserRepo) DeleteUnactiveUser(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := at.collection.DeleteOne(ctx, filter)
	return err
}

func (at *UnactiveUserRepo) UpdateActiveToken(ctx context.Context, email, token string, expiry time.Time) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"activation_token":         token,
			"activation_token_expiry":  expiry,
		},
	}
	_, err := at.collection.UpdateOne(ctx, filter, update)
	return err
}