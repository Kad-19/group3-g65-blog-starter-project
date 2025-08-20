package image

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudinaryService_UploadImage(t *testing.T) {
	// This is an integration test and requires Cloudinary credentials to be set in the environment
	if os.Getenv("CLOUD_NAME") == "" || os.Getenv("API_KEY") == "" || os.Getenv("API_SECERT") == "" {
		t.Skip("Skipping Cloudinary integration test: credentials not set")
	}

	service := NewCloudinaryService()
	ctx := context.Background()
	file := strings.NewReader("fake image data")
	folderName := "test-uploads"

	url, err := service.UploadImage(ctx, file, folderName)

	assert.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "res.cloudinary.com")
}
