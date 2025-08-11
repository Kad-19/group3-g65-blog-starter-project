package domain

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents the core user entity in the system
type User struct {
	ID           string
	Username     string
	Email        string
	Password     string
	Role         string
	Activated    bool
	Profile      UserProfile
	CreatedAt    time.Time
	UpdatedAt    time.Time 
}

// UserProfile represents embedded profile data
type UserProfile struct {
	Bio               string
	ProfilePictureURL string 
	ContactInfo       string 
}

// UnactivatedUser represents a user who has not yet activated their account
type UnactivatedUser struct {
	ID                    string 
	Username              string
	Email                 string
	Password              string 
	Activated			  bool 
	ActivationToken       string  
	ActivationTokenExpiry *time.Time  
	CreatedAt             time.Time 
	UpdatedAt             time.Time 
}