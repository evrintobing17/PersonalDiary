package repository

import (
	"github.com/evrintobing17/PersonalDiary/domain"
)

// UserRepository represents the users repository
type UserRepository interface {
	CreateUser(*domain.User) (*domain.User, error)
	UpdateUser(uint64, *domain.User) (*domain.User, error)
	DeleteUser(uid uint64) (int64, error)
	GetUserByID(uid uint64) (*domain.User, error)
	GetAllUsers() (*[]domain.User, error)
}
