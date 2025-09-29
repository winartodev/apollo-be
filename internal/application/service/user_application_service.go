package service

import (
	"context"

	"github.com/winartodev/apollo-be/internal/domain/entities"
	"github.com/winartodev/apollo-be/internal/domain/repository"
)

// UserApplicationService handles cross-cutting user operations across modules
type UserApplicationService interface {
	// User existence and validation operations
	UserExists(ctx context.Context, username string) (bool, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)

	// User retrieval operations
	GetUserByID(ctx context.Context, id int64) (*entities.SharedUser, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.SharedUser, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.SharedUser, error)

	// User status operations
	ActivateUser(ctx context.Context, userID int64) error
	DeactivateUser(ctx context.Context, userID int64) error
	VerifyUserEmail(ctx context.Context, userID int64) error
	VerifyUserPhone(ctx context.Context, userID int64) error
}

type userApplicationService struct {
	userRepo repository.UserRepository
}

func NewUserApplicationService(userRepo repository.UserRepository) UserApplicationService {
	return &userApplicationService{
		userRepo: userRepo,
	}
}

func (s *userApplicationService) UserExists(ctx context.Context, username string) (bool, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return user != nil && !user.IsDeleted(), nil
}

func (s *userApplicationService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil && !user.IsDeleted(), nil
}

func (s *userApplicationService) GetUserByID(ctx context.Context, id int64) (*entities.SharedUser, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userApplicationService) GetUserByUsername(ctx context.Context, username string) (*entities.SharedUser, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *userApplicationService) GetUserByEmail(ctx context.Context, email string) (*entities.SharedUser, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userApplicationService) ActivateUser(ctx context.Context, userID int64) error {
	return s.userRepo.UpdateStatus(ctx, userID, true)
}

func (s *userApplicationService) DeactivateUser(ctx context.Context, userID int64) error {
	return s.userRepo.UpdateStatus(ctx, userID, false)
}

func (s *userApplicationService) VerifyUserEmail(ctx context.Context, userID int64) error {
	return s.userRepo.UpdateEmailVerification(ctx, userID, true)
}

func (s *userApplicationService) VerifyUserPhone(ctx context.Context, userID int64) error {
	return s.userRepo.UpdatePhoneVerification(ctx, userID, true)
}
