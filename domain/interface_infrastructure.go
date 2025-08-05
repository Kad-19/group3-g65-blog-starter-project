package domain

import (
	"context"
	"io"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file io.Reader, folderName string) (string, error)
}

type EmailProvider interface {
	SendPasswordResetEmail(toEmail, resetLink string) error
	SendActivationEmail(toEmail, activationLink string) error
}

type AIService interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}
