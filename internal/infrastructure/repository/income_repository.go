package repository

import (
	"fmt"
	"time"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
	"gorm.io/gorm"
)

type IncomeRepo struct {
	db *gorm.DB
}

func NewIncomeRepo(db *gorm.DB) *IncomeRepo {
	return &IncomeRepo{
		db: db,
	}
}

func (rep *IncomeRepo) Create(income *entity.Income) error {
	start, end := rangeWeek("now")
	if result := rep.db.Where("tg_id = ? AND created_at BETWEEN ? AND ?", income.TGID, start, end).FirstOrCreate(income); result.Error != nil && result.RowsAffected != 1 {
		return fmt.Errorf("error creating income - %w", result.Error)
	}
	return nil
}

func (rep *IncomeRepo) Select(tgid int64, week string) []entity.Income {
	incomes := []entity.Income{}
	startWeek, endWeek := rangeWeek(week)
	res := rep.db.Where("tg_id = ? AND created_at BETWEEN ? AND ?", tgid, startWeek, endWeek).Find(&incomes)
	if res.Error != nil {
		fmt.Printf("error select income: %v", res.Error)
		return nil
	}
	return incomes
}

func rangeWeek(week string) (time.Time, time.Time) {
	y, m, d := time.Now().Date()
	weekday := int(time.Now().Weekday()) - 1
	if weekday < 0 {
		weekday = 7
	}
	startWeek := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	endWeek := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	switch week {
	case "last":
		startWeek = startWeek.AddDate(0, 0, -weekday-6)
		endWeek = endWeek.AddDate(0, 0, -weekday)
	case "current":
		startWeek = startWeek.AddDate(0, 0, -weekday+1)
	}
	return startWeek, endWeek
}
