package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/winartodev/apollo-be/infrastructure/database"
	"github.com/winartodev/apollo-be/internal/domain/entities"
	domainError "github.com/winartodev/apollo-be/internal/domain/error"
	"github.com/winartodev/apollo-be/modules/auth/domain/repository"
)

type AuthRepositoryImpl struct {
	*database.Database
}

func NewAuthRepository(db *database.Database) (repository.AuthRepository, error) {
	return &AuthRepositoryImpl{
		Database: db,
	}, nil
}

func (ar *AuthRepositoryImpl) RegisterNewUserDB(ctx context.Context, data entities.SharedUser) (id *int64, err error) {
	stmt, err := ar.DB.PrepareContext(ctx, registerUserQuery)
	if err != nil {
		return nil, fmt.Errorf(database.ErrFailedPrepareStatement, err)
	}

	defer ar.Database.CloseStatement(stmt, &err)

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
		return nil, domainError.ErrFailedCreateUser
	}

	return &lastInsertID, nil
}

func (ar *AuthRepositoryImpl) UpdateRefreshTokenDB(ctx context.Context, id int64, token *string) (err error) {
	stmt, err := ar.DB.PrepareContext(ctx, updateRefreshTokenQuery)
	if err != nil {
		return fmt.Errorf(database.ErrFailedPrepareStatement, err)
	}

	defer ar.Database.CloseStatement(stmt, &err)

	updatedAt := time.Now().Unix()
	_, err = stmt.ExecContext(
		ctx,
		id,
		&token,
		updatedAt,
	)
	if err != nil {
		return domainError.ErrFailedUpdateRefreshToken
	}

	return nil
}

func (ar *AuthRepositoryImpl) GetUserDataDB(ctx context.Context, username string) (data *entities.SharedUser, err error) {
	result := &entities.SharedUser{}
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

		return nil, domainError.ErrFailedGetUserData
	}

	return result, nil
}
