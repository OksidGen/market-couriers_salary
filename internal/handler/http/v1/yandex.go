package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type YandexHandler struct {
	CourierUseCase CourierUseCase
}

const botLink = "https://t.me/MarketSalaryBot"

func NewYandexHandler(CourierUseCase CourierUseCase) *YandexHandler {
	return &YandexHandler{
		CourierUseCase: CourierUseCase,
	}
}

func (h *YandexHandler) initYandexHandlerPath(api *gin.RouterGroup) {
	tg := api.Group("/yandex")
	{
		tg.GET("/code", h.YandexRedirectHandler)
	}
}

func (h *YandexHandler) YandexRedirectHandler(c *gin.Context) {
	code := c.Query("code")
	tgId, err := strconv.ParseInt(c.Query("state"), 10, 64)
	if err != nil {
		log.Printf("error parse tgId from yandex `state` : %v", err)
	}
	if err = h.CourierUseCase.GetToken(code, tgId); err != nil {
		c.String(http.StatusInternalServerError, "Внутренняя ошибка, пожалуйста сообщите мне о ней, используя телеграмм-бота - раздел `Помощь`")
	}
	c.Redirect(http.StatusMovedPermanently, botLink)
}
