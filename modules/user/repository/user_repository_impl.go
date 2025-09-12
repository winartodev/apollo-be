package repository

import (
	"context"
	"fmt"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/user/domain/entities"
	"github.com/winartodev/apollo-be/modules/user/domain/repository"
)

type UserRepositoryImpl struct {
	*helper.DatabaseUtil
}

func NewUserRepository(database *helper.DatabaseUtil) (repository.UserRepository, error) {
	return &UserRepositoryImpl{
		DatabaseUtil: database,
	}, nil
}

func (ur *UserRepositoryImpl) GetUserByIDDB(ctx context.Context, id int64) (user *entities.User, err error) {
	user = &entities.User{}

	err = ur.DB.QueryRowContext(ctx, getUserByID, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) GetUserByEmailDB(ctx context.Context, email string) (user *entities.User, err error) {
	return ur.getUserByField(ctx, "email", email)
}

func (ur *UserRepositoryImpl) GetUserByUsernameDB(ctx context.Context, username string) (user *entities.User, err error) {
	return ur.getUserByField(ctx, "username", username)
}

func (ur *UserRepositoryImpl) getUserByField(ctx context.Context, field, value string) (*entities.User, error) {
	user := &entities.User{}

	query := fmt.Sprintf("%s WHERE usr.%s = $1", checkUserExists, field)

	err := ur.DB.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
