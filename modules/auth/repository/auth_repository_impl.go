package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/winartodev/apollo-be/core/helper"
	"github.com/winartodev/apollo-be/modules/auth/domain/entities"
	authDomainError "github.com/winartodev/apollo-be/modules/auth/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

type AuthRepositoryImpl struct {
	*helper.DatabaseUtil
}

func NewAuthRepository(database *helper.DatabaseUtil) (repository.AuthRepository, error) {
	return &AuthRepositoryImpl{
		DatabaseUtil: database,
	}, nil
}

func (ar *AuthRepositoryImpl) RegisterNewUserDB(ctx context.Context, data entities.User) (id *int64, err error) {
	stmt, err := ar.DB.PrepareContext(ctx, registerUserQuery)
	if err != nil {
		return nil, fmt.Errorf(helper.ErrFailedPrepareStatement, err)
	}

	defer ar.DatabaseUtil.CloseStatement(stmt, &err)

	createdAt := time.Now().Unix()
	var lastInsertID int64
	err = stmt.QueryRowContext(ctx,
		data.Username,
		data.Email,
		data.PhoneNumber,
		data.FirstName,
		data.LastName,
		data.Password,
		createdAt,
	).Scan(&lastInsertID)
	if err != nil {
		return nil, fmt.Errorf(authDomainError.ErrFailedCreateUser, err)
	}

	return &lastInsertID, nil
}

func (ar *AuthRepositoryImpl) UpdateRefreshTokenDB(ctx context.Context, id int64, token *string) (err error) {
	stmt, err := ar.DB.PrepareContext(ctx, updateRefreshTokenQuery)
	if err != nil {
		return fmt.Errorf(helper.ErrFailedPrepareStatement, err)
	}

	defer ar.DatabaseUtil.CloseStatement(stmt, &err)

	updatedAt := time.Now().Unix()
	_, err = stmt.ExecContext(
		ctx,
		id,
		&token,
		updatedAt,
	)
	if err != nil {
		return fmt.Errorf(authDomainError.ErrFailedUpdateRefreshToken, err)
	}

	return nil
}

func (ar *AuthRepositoryImpl) GetUserDataDB(ctx context.Context, username string) (data *entities.User, err error) {
	result := &entities.User{}
	err = ar.DB.QueryRowContext(ctx, getUserData, username, username).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf(authDomainError.ErrFailedGetUserData, err)
	}

	return result, nil
}
