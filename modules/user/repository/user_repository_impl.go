package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/winartodev/apollo-be/infrastructure/database"
	"github.com/winartodev/apollo-be/modules/user/domain/entities"
	"github.com/winartodev/apollo-be/modules/user/domain/repository"
)

type UserRepositoryImpl struct {
	*database.Database
}

func NewUserRepository(db *database.Database) (repository.UserRepository, error) {
	return &UserRepositoryImpl{
		Database: db,
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

	query := fmt.Sprintf("%s WHERE usr.%s = $1", getUserQuery, field)

	err := ur.DB.QueryRowContext(ctx, query, value).Scan(
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
