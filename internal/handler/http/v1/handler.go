package v1

import (
	"net/http"

	"github.com/OksidGen/market-couriers_salary/internal/usecase"
	"github.com/gin-gonic/gin"
)

type (
	CourierUseCase interface {
		TelegramParser(*http.Request) map[string]interface{}
		SendMessage(map[string]interface{})

		CaseStart(map[string]interface{})
		CaseAuth(map[string]interface{})
		CaseCheck(map[string]interface{})
		CaseCalculate(map[string]interface{})
		CaseHelp(map[string]interface{})
		CaseWeek(map[string]interface{}, string)

		GetToken(string, int64) error
	}

	Handlers struct {
		TelegramHandler *TelegramHandler
		YandexHandler   *YandexHandler
	}
)

func New(uc *usecase.UseCases) *Handlers {
	return &Handlers{
		TelegramHandler: NewTelegramHandler(uc.CourierUseCase),
		YandexHandler:   NewYandexHandler(uc.CourierUseCase),
	}
}

func (h *Handlers) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.TelegramHandler.initTelegramHandlerPath(v1)
		h.YandexHandler.initYandexHandlerPath(v1)
	}
}
