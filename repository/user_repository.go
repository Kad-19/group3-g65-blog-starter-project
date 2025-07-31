package repository

import (
	"context"
	"fmt"
	"g3-g65-bsp/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ domain.UserRepository = (*MongoRepo)(nil)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(mc *mongo.Collection) *MongoRepo {
	return &MongoRepo{
		collection: mc,
	}
}

func (mr *MongoRepo) GetUserByEmail(ctx context.Context, Email string) (*domain.User, error) {
	user := &domain.User{}
	filter := bson.M{"email": Email}
	err := mr.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (mr *MongoRepo) UpdateUser(ctx context.Context, up domain.UserProfile, Email string) error {
	filter := bson.M{"email": Email}
	update := bson.M{
		"$set": bson.M{
			"bio":                 up.Bio,
			"profile_picture_url": up.ProfilePictureURL,
			"contact_information": up.ContactInfo,
		},
	}

	if res, err := mr.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	} else {
		if res.MatchedCount == 0 {
			return fmt.Errorf("no user found")
		}
		return nil
	}
}

func (mr *MongoRepo) UpdateUserRole(ctx context.Context, role string, Email string) error {
	filter := bson.M{"email": Email}
	update := bson.M{
		"$set": bson.M{
			"role": role,
		},
	}

	if res, err := mr.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	} else {
		if res.MatchedCount == 0 {
			return fmt.Errorf("no user found")
		}
		return nil
	}
}
