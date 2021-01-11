package repository

import "github.com/evrintobing17/PersonalDiary/domain"

// AuthRepository handles users authentication
type AuthRepository interface {
	Login(email, password string) (*domain.TokenDetail, error)
	Logout(uuid string) (int64, error)
	Refresh(refreshUUID string) (*domain.TokenDetail, error)
}
