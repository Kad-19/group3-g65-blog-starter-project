package usecase

import (
	"context"
	"fmt"
	"g3-g65-bsp/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(ur domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: ur,
	}
}

func (upd *UserUseCase) Promote(ctx context.Context, Email string) error {
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

func (upd *UserUseCase) Demote(ctx context.Context, Email string) error {
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

func (upd *UserUseCase) ProfileUpdate(ctx context.Context, up *domain.UserProfile, Email string) error {
	if up.ProfilePictureURL == "" && up.Bio == "" && up.ContactInfo == "" {
		return fmt.Errorf("all fields are required")
	}
	if err := upd.userRepo.UpdateUser(ctx, *up, Email); err != nil {
		return err
	}
	return nil
}
