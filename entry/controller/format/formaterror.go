package format

import (
	"errors"

	"github.com/evrintobing17/PersonalDiary/domain"
)

// Validate is used to check if the entry has a correct input format
func Validate(e *domain.Entry) error {
	if e.Title == "" {
		return errors.New("Required Title")
	}
	if e.Description == "" {
		return errors.New("Required Description")
	}
	if e.OwnerID < 1 {
		return errors.New("Required Owner of the Post")
	}
	return nil
}
