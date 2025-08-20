package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set up mock environment variables
	os.Setenv("MONGODB_DB", "test_db")
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("ACCESS_TOKEN_SECRET", "access_secret")
	os.Setenv("REFRESH_TOKEN_SECRET", "refresh_secret")
	os.Setenv("ACCESS_TOKEN_EXPIRY", "15m")
	os.Setenv("REFRESH_TOKEN_EXPIRY", "72h")
	os.Setenv("GOOGLE_OAUTH_CLIENT_ID", "google_id")
	os.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "google_secret")
	os.Setenv("OAUTH_STATE_STRING", "random_string")

	// Clean up environment variables after the test
	defer func() {
		os.Unsetenv("MONGODB_DB")
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("ACCESS_TOKEN_SECRET")
		os.Unsetenv("REFRESH_TOKEN_SECRET")
		os.Unsetenv("ACCESS_TOKEN_EXPIRY")
		os.Unsetenv("REFRESH_TOKEN_EXPIRY")
		os.Unsetenv("GOOGLE_OAUTH_CLIENT_ID")
		os.Unsetenv("GOOGLE_OAUTH_CLIENT_SECRET")
		os.Unsetenv("OAUTH_STATE_STRING")
	}()

	// Load the configuration
	LoadConfig()

	// Assert that the configuration is loaded correctly
	assert.NotNil(t, AppConfig)
	assert.Equal(t, "test_db", AppConfig.DbName)
	assert.Equal(t, "mongodb://localhost:27017", AppConfig.MongoURI)
	assert.Equal(t, "access_secret", AppConfig.AccessTokenSecret)
	assert.Equal(t, "refresh_secret", AppConfig.RefreshTokenSecret)
	assert.Equal(t, 15*time.Minute, AppConfig.AccessTokenExpiry)
	assert.Equal(t, 72*time.Hour, AppConfig.RefreshTokenExpiry)
	assert.Equal(t, "google_id", AppConfig.GoogleClientID)
	assert.Equal(t, "google_secret", AppConfig.GoogleClientSecret)
	assert.Equal(t, "random_string", AppConfig.OauthStateString)
}
