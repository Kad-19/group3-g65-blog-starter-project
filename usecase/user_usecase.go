package usecase

import (
	"context"
	"errors"
	"fmt"
	"g3-g65-bsp/domain"
	"io"
)

type UserUsecase struct {
	userRepo      domain.UserRepository
	imageUploader domain.ImageUploader
}

func NewUserUsecase(ur domain.UserRepository, iu domain.ImageUploader) domain.UserUsecase {
	return &UserUsecase{
		userRepo:      ur,
		imageUploader: iu,
	}
}

func (upd *UserUsecase) Promote(ctx context.Context, Email string) error {
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

func (upd *UserUsecase) Demote(ctx context.Context, Email string) error {
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

func (upd *UserUsecase) ProfileUpdate(ctx context.Context, userid string, bio string, contactinfo string, file io.Reader) error {
	user, err := upd.userRepo.FindByID(ctx, userid)
	if err != nil {
		return err
	}

	imageURL, err := upd.imageUploader.UploadImage(ctx, file, "profile")
	if err != nil {
		return errors.New("failed to upload image")
	}

	if err := upd.userRepo.UpdateUserProfile(ctx, bio, contactinfo, imageURL, user.Email); err != nil {
		return err
	}
	return nil
}

func (upd *UserUsecase) GetAllUsers(ctx context.Context) 