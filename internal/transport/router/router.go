package router

import (
	"music_storage/internal/service"
	v1 "music_storage/internal/transport/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	hanlderV1 := v1.NewHandler(h.services)
	api := router.Group("/api")
	{
		hanlderV1.Init(api)
	}
}
