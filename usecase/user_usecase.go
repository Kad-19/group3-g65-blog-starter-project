package usecase

import (
	"context"
	"errors"
	"fmt"
	"g3-g65-bsp/domain"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOperations struct {
	userRepo      domain.UserRepository
	imageUploader domain.ImageUploader
}

func NewUserUseCase(ur domain.UserRepository, iu domain.ImageUploader) *UserOperations {
	return &UserOperations{
		userRepo:      ur,
		imageUploader: iu,
	}
}

func (upd *UserOperations) Promote(ctx context.Context, Email string) error {
	user, err := upd.userRepo.FindByEmail(ctx, Email)
	if err != nil {
		return err
	}
	if user.Role == string(domain.RoleAdmin) {
		return fmt.Errorf("the user is already an admin")
	}

	if ok := upd.userRepo.UpdateUserRole(ctx, string(domain.RoleAdmin), Email); ok != nil {
		return ok
	}
	return nil
}

func (upd *UserOperations) Demote(ctx context.Context, Email string) error {
	user, err := upd.userRepo.FindByEmail(ctx, Email)
	if err != nil {
		return err
	}
	if user.Role == string(domain.RoleUser) {
		return fmt.Errorf("the user is already a user")
	}
	if ok := upd.userRepo.UpdateUserRole(ctx, string(domain.RoleUser), Email); ok != nil {
		return ok
	}

	return nil
}

func (upd *UserOperations) ProfileUpdate(ctx context.Context, userid primitive.ObjectID, bio string, contactinfo string, file io.Reader) error {
	user, err := upd.userRepo.FindByID(ctx, userid)
	if err != nil {
		return err
	}

	imageURL, err := upd.imageUploader.UploadImage(ctx, file, "uploads/profile")
	if err != nil {
		return errors.New("failed to upload image")
	}

	if err := upd.userRepo.UpdateUser(ctx, bio, contactinfo, imageURL, user.Email); err != nil {
		return err
	}
	return nil
}
