package repository

import (
	"context"

	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
)

type AuthRepository interface {
	RegisterNewUserDB(ctx context.Context, data entities.User) (id *int64, err error)
	UpdateRefreshTokenDB(ctx context.Context, id int64, token *string) (err error)
	GetUserDataDB(ctx context.Context, username string) (data *entities.User, err error)
}
