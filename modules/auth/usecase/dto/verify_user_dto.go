package dto

import (
	domainEntity "github.com/winartodev/apollo-be/internal/domain/entities"
)

type VerifyUserDto struct {
	User        *domainEntity.SharedUser
	Suggestions []string
}
