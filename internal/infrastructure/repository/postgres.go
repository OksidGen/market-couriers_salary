package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo   *UserRepo
	TokenRepo  *TokenRepo
	IncomeRepo *IncomeRepo
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:   NewUserRepo(db),
		TokenRepo:  NewTokenRepo(db),
		IncomeRepo: NewIncomeRepo(db),
	}
}
