package db

import (
	"fmt"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting database: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error{
	if err:=db.AutoMigrate(&entity.User{},&entity.Token{},&entity.Income{}); err != nil {
		return fmt.Errorf("error autoMigrate database: %w", err)
	}
	return nil
}
