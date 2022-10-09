package entity

import (
	"fmt"
	"time"
)

type (
	User struct {
		ID        uint `gorm:"primaryKey"`
		TGID      int64
		FirstName string
		LastName  string
		UserName  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Token struct {
		ID           uint `gorm:"primaryKey"`
		TGID         int64
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	Income struct {
		ID           uint `gorm:"primaryKey"`
		TGID         int64
		CreatedAt    time.Time
		Income       int
		Count        int
		Client       int
		PvzPstLocker int
		Multi5       int
	}
)

func (income *Income) ToTableStirng() string {
	out := fmt.Sprintf(`
|        Дата       | %v |
|:------------------|-----------:|
| Кол-во пвз/пст    | %v
| Кол-во клиентов   | %v
| Посылок в пвз/пст | %v
| 5+ мультизаказов  | %v
|:------------------|-----------:|
|        Итого      |    %v ₽  |`,
		income.CreatedAt.Format("2006-01-02"),
		income.Count,
		income.Client,
		income.PvzPstLocker,
		income.Multi5,
		income.Income)
	return out
}
