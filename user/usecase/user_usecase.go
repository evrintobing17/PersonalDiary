package usecase

import "github.com/evrintobing17/PersonalDiary/domain"

// UserUseCase represents the users usecase
type UserUseCase interface {
	CreateUser(*domain.User) (*domain.User, error)
	UpdateUser(uint64, *domain.User) (*domain.User, error)
	DeleteUser(uint64) (int64, error)
	GetUserByID(uint64) (*domain.User, error)
	GetAllUsers() (*[]domain.User, error)
}
