package repository

import (
	"errors"
	"fmt"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
	"gorm.io/gorm"
)

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (rep *TokenRepo) Create(token *entity.Token) error {
	if result := rep.db.FirstOrCreate(token); result.Error != nil && result.RowsAffected != 1 {
		return fmt.Errorf("error creating token - %w", result.Error)
	}
	return nil
}
func (rep *TokenRepo) Find(tgid int64) (string, error) {
	token := entity.Token{TGID: tgid}
	if res := rep.db.First(&token); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("error find token: %w", res.Error)
		}
		return "", nil
	}
	return token.AccessToken, nil
}
