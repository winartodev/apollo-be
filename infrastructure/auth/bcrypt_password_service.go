package auth

import (
	"github.com/winartodev/apollo-be/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct {
}

func NewBcryptPasswordService() domain.PasswordService {
	return &BcryptPasswordService{}
}

func (b BcryptPasswordService) HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (b BcryptPasswordService) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
