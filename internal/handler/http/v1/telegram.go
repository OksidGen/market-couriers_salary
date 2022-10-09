package v1

import (
	"github.com/gin-gonic/gin"
)

type TelegramHandler struct {
	CourierUseCase CourierUseCase
}

func NewTelegramHandler(CourierUseCase CourierUseCase) *TelegramHandler {
	return &TelegramHandler{
		CourierUseCase: CourierUseCase,
	}
}

func (h *TelegramHandler) initTelegramHandlerPath(api *gin.RouterGroup) {
	tg := api.Group("/telegram")
	{
		tg.POST("/webhook", h.TelegramWebhookHandler)
	}
}

func (h *TelegramHandler) TelegramWebhookHandler(c *gin.Context) {
	upd := h.CourierUseCase.TelegramParser(c.Request)
	msg := map[string]interface{}{
		"chat_id": upd["FromId"],
		"text":    "Неизвестная команда",
	}
	switch upd["Text"] {
	case "Помощь":
		h.CourierUseCase.CaseHelp(msg)
	case "Прошлая неделя":
		h.CourierUseCase.CaseWeek(msg, "last")
	case "Текущая неделя":
		h.CourierUseCase.CaseWeek(msg, "current")
	case "Расчет смены":
		h.CourierUseCase.CaseCalculate(msg)
	case "Проверка":
		h.CourierUseCase.CaseCheck(msg)
	case "Авторизация":
		msg["FirstName"] = upd["FirstName"]
		msg["LastName"] = upd["LastName"]
		msg["UserName"] = upd["UserName"]
		h.CourierUseCase.CaseAuth(msg)
	case "/start":
		h.CourierUseCase.CaseStart(msg)
	}
	delete(msg, "FirstName")
	delete(msg, "LastName")
	delete(msg, "Username")
	h.CourierUseCase.SendMessage(msg)
}
