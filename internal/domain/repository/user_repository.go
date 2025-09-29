package repository

import (
	"context"

	"github.com/winartodev/apollo-be/internal/domain/entities"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *entities.SharedUser) (*entities.SharedUser, error)
	GetByID(ctx context.Context, id int64) (*entities.SharedUser, error)
	GetByUsername(ctx context.Context, username string) (*entities.SharedUser, error)
	GetByEmail(ctx context.Context, email string) (*entities.SharedUser, error)
	Update(ctx context.Context, user *entities.SharedUser) error
	Delete(ctx context.Context, id int64) error

	// Status update operations
	UpdateStatus(ctx context.Context, id int64, isActive bool) error
	UpdateEmailVerification(ctx context.Context, id int64, isVerified bool) error
	UpdatePhoneVerification(ctx context.Context, id int64, isVerified bool) error
	UpdateRefreshToken(ctx context.Context, id int64, token *string) error
	UpdateLastLogin(ctx context.Context, id int64) error

	// Query operations
	List(ctx context.Context, offset, limit int) ([]*entities.SharedUser, error)
	Count(ctx context.Context) (int64, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
