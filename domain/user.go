package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

// User represents the core user entity in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	Password 	 string             `bson:"password" json:"-"`
	Role         string             `bson:"role" json:"role"`
	Activated    bool               `bson:"activated" json:"activated"`
	Profile      UserProfile        `bson:"profile" json:"profile"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserProfile represents embedded profile data
type UserProfile struct {
	Bio               string `bson:"bio,omitempty" json:"bio" validate:"max=500"`
	ProfilePictureURL string `bson:"profile_picture_url,omitempty" json:"profile_picture_url"`
	ContactInfo       string `bson:"contact_information,omitempty" json:"contact_information" validate:"max=100"`
}