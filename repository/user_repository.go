package repository

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserDTO represents the user data transfer object
type UserDTO struct {
	ID           primitive.ObjectID		`bson:"_id,omitempty"`
	Username     string             	`bson:"username"`
	Email        string             	`bson:"email"`
	Password 	 string             	`bson:"password"`
	Role         string             	`bson:"role"`
	Activated    bool               	`bson:"activated"`
	Profile      UserProfileDTO      	`bson:"profile"`
	CreatedAt    time.Time          	`bson:"created_at"`
	UpdatedAt    time.Time          	`bson:"updated_at"`
}

// UserProfileDTO represents the user profile data transfer object
type UserProfileDTO struct {
	Bio               string `bson:"bio,omitempty"`
	ProfilePictureURL string `bson:"profile_picture_url,omitempty"`
	ContactInfo       string `bson:"contact_information,omitempty"`
}

// ConvertToDomain converts UserDTO to domain.User
func (dto *UserDTO) ConvertToUserDomain() *domain.User {
	return &domain.User{
		ID:        dto.ID.Hex(),
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      dto.Role,
		Activated: dto.Activated,
		Profile: domain.UserProfile{
			Bio:               dto.Profile.Bio,
			ProfilePictureURL: dto.Profile.ProfilePictureURL,
			ContactInfo:       dto.Profile.ContactInfo,
		},
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// ConvertToDTO converts domain.User to UserDTO
func ConvertToUserDTO(u *domain.User) *UserDTO {
	id, _ := primitive.ObjectIDFromHex(u.ID)

	return &UserDTO{
		ID:        id,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Activated: u.Activated,
		Profile: UserProfileDTO{
			Bio:               u.Profile.Bio,
			ProfilePictureURL: u.Profile.ProfilePictureURL,
			ContactInfo:       u.Profile.ContactInfo,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.collection.InsertOne(ctx, ConvertToUserDTO(user))
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user UserDTO
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return user.ConvertToUserDomain(), err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user UserDTO
	idObj, _ := primitive.ObjectIDFromHex(id)
	err := r.collection.FindOne(ctx, bson.M{"_id": idObj}).Decode(&user)
	return user.ConvertToUserDomain(), err
}

func (mr *UserRepository) UpdateUserProfile(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error {
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
			return errors.New("user not found")
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
			return errors.New("user not found")
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
			return errors.New("user not found")
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
			return errors.New("user not found")
		}
		return nil
	}
}