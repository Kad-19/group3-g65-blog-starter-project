package image

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService() *CloudinaryService {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cls, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Cloudinary client: %v", err))
	}
	return &CloudinaryService{
		client: cls,
	}
}

func (cs *CloudinaryService) UploadImage(ctx context.Context, file io.Reader, folderName string) (string, error) {
	uploadparams := uploader.UploadParams{
		Folder: folderName,
	}

	resp, err := cs.client.Upload.Upload(ctx, file, uploadparams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	return resp.SecureURL, nil
}
