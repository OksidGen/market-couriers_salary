package repository

import (
	"errors"
	"fmt"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (rep *UserRepo) Create(user *entity.User) error {
	if result := rep.db.Create(user); result.Error != nil {
		return fmt.Errorf("error creating user - %w", result.Error)
	}
	return nil
}

func (rep *UserRepo) Check(id int64) (bool, error) {
	user := entity.User{TGID: id}
	if result := rep.db.First(&user); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("error checking user: %w", result.Error)
		}
		return false, nil
	}
	return true, nil
}
