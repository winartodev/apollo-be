package repository

import (
	"context"

	"github.com/winartodev/apollo-be/internal/domain/entities"
)

type AuthRepository interface {
	RegisterNewUserDB(ctx context.Context, data entities.SharedUser) (id *int64, err error)
	UpdateRefreshTokenDB(ctx context.Context, id int64, token *string) (err error)
	GetUserDataDB(ctx context.Context, username string) (data *entities.SharedUser, err error)
	UpdatePasswordDB(ctx context.Context, id int64, hashedPassword string) (err error)
}
