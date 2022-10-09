package http

import (
	"net/http"

	v1 "github.com/OksidGen/market-couriers_salary/internal/handler/http/v1"
	"github.com/OksidGen/market-couriers_salary/internal/usecase"
	"github.com/gin-gonic/gin"
)

type (
	Handlers struct {
		v1 *v1.Handlers
	}
)

func New(uc *usecase.UseCases) *Handlers {
	return &Handlers{
		v1: v1.New(uc),
	}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.Default()

	// router.Use(
	// 	gin.Recovery(),
	// 	gin.Logger(),
	// )

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusTeapot, "418 - I`m a teapot")
	})
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initApi(router)

	return router
}

func (h *Handlers) initApi(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.v1.Init(api)
	}
}
