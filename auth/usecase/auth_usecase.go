package usecase

import "github.com/evrintobing17/PersonalDiary/domain"

// AuthUseCase represents the users authentication usecase
type AuthUseCase interface {
	Login(email, password string) (*domain.TokenDetail, error)
	Logout(uuid string) (int64, error)
	Refresh(refreshUUID string) (*domain.TokenDetail, error)
}
