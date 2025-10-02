package repository

import (
	"context"

	entity "github.com/winartodev/apollo-be/modules/user/domain/entities"
)

type UserRepository interface {
	GetUserByIDDB(ctx context.Context, id int64) (user *entity.User, err error)
	GetUserByEmailDB(ctx context.Context, email string) (user *entity.User, err error)
	GetUserByUsernameDB(ctx context.Context, username string) (user *entity.User, err error)
}
