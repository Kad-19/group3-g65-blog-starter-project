package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName             string
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// AppConfig is the global config instance
var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}

	dbName := os.Getenv("MONGODB_DB")

	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	accessExpiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	refreshExpiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")

	accessExpiry, err := time.ParseDuration(accessExpiryStr)
	if err != nil {
		log.Fatal("Invalid ACCESS_TOKEN_EXPIRY value")
	}

	refreshExpiry, err := time.ParseDuration(refreshExpiryStr)
	if err != nil {
		log.Fatal("Invalid REFRESH_TOKEN_EXPIRY value")
	}

	AppConfig = &Config{
		DbName 			:   dbName,
		AccessTokenSecret:  accessSecret,
		RefreshTokenSecret: refreshSecret,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
	}
}
