package controller

import (
	"fmt"
	"g3-g65-bsp/config" 
	"g3-g65-bsp/domain"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthController handles Google OAuth2 authentication requests.
type OAuthController struct {
	usecase domain.OAuthUsecase 
}

// NewOAuthController creates a new instance of OAuthController.
func NewOAuthController(uc domain.OAuthUsecase) *OAuthController {
	return &OAuthController{usecase: uc}
}

// getGoogleOauthConfig initializes and returns the Google OAuth2 configuration.
func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     config.AppConfig.GoogleClientID,
		ClientSecret: config.AppConfig.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// HandleGoogleLogin redirects the user to Google's OAuth consent screen.
func (oc *OAuthController) HandleGoogleLogin(c *gin.Context) {
	googleOauthConfig := getGoogleOauthConfig()
	oauthStateString := config.AppConfig.OauthStateString 

	// Generate the authorization URL and redirect the user
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback processes the callback from Google after user authentication.
func (oc *OAuthController) HandleGoogleCallback(c *gin.Context) {
	googleOauthConfig := getGoogleOauthConfig()
	oauthStateString := config.AppConfig.OauthStateString

	// Validate OAuth state to prevent CSRF attacks
	if c.Query("state") != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	// Extract the authorization code from the query parameters
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Call the OAuthLogin usecase to handle token exchange, user info retrieval,
	accessToken, refreshToken, accessExpirySeconds, user, err := oc.usecase.OAuthLogin(c.Request.Context(), *googleOauthConfig, code)
	if err != nil {
		// Log the error for debugging purposes (optional, but recommended)
		fmt.Printf("Error during OAuthLogin: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// On successful login, return the tokens and user information
	c.JSON(http.StatusOK, gin.H{
		"refresh_token": refreshToken,
		"access_token":  accessToken,
		"expiry_in":     accessExpirySeconds, 
		"user":          user,
	})
}