package usecase

import (
	"net/http"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
	"github.com/OksidGen/market-couriers_salary/internal/infrastructure/repository"
	"github.com/OksidGen/market-couriers_salary/internal/infrastructure/webapi"
)

type (
	UserRepo interface {
		Create(*entity.User) error
		Check(int64) (bool, error)
	}
	TokenRepo interface {
		Create(*entity.Token) error
		Find(int64) (string, error)
	}
	IncomeRepo interface {
		Create(*entity.Income) error
		Select(int64, string) []entity.Income
	}

	TelegramWebApi interface {
		ParseRequest(*http.Request) (map[string]interface{}, error)
		SendMessage(map[string]interface{})
	}
	YandexWebApi interface {
		GetURLForYandexAuth(int64) string
		CheckCompletedTask(string) string
		CalculateIncome(string) map[string]int
		SendCodeForToken(string) []byte
	}

	UseCases struct {
		CourierUseCase *CourierUseCase
	}
)

func New(repo *repository.Repositories, webapi *webapi.WebApi) *UseCases {
	return &UseCases{
		CourierUseCase: NewCourierUseCase(
			repo.UserRepo,
			repo.TokenRepo,
			repo.IncomeRepo,
			webapi.TelegramWebApi,
			webapi.YandexWebApi,
		),
	}
}
