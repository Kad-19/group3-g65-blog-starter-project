package domain

import (
	"context"
	"io"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file io.Reader, folderName string) (string, error)
}
