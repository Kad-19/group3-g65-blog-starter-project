package repository

import (
	"context"
	"fmt"

	// <<<<<<< temesgen_user_manag
	"errors"
	"g3-g65-bsp/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (mr *UserRepository) UpdateUser(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error {
	filter := bson.M{"email": Email}
	update := bson.M{
		"$set": bson.M{
			"bio":                 bio,
			"profile_picture_url": imagePath,
			"contact_information": contactInfo,
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

func (mr *UserRepository) UpdateUserRole(ctx context.Context, role string, Email string) error {
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

func (mr *UserRepository) UpdateActiveStatus(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"activated": true,
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

func (mr *UserRepository) UpdateUserPassword(ctx context.Context, email string, newPasswordHash string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"password": newPasswordHash,
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

// >>>>>>> main
