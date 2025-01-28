package router

import (
	"music_storage/docs"
	"music_storage/internal/service"
	v1 "music_storage/internal/transport/api/v1"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *http.ServeMux {
	v1 := v1.NewHandler(h.services)

	apiRoutes := http.NewServeMux()
	apiRoutes.Handle("/api/", http.StripPrefix("/api", v1.Init()))

	apiRoutes.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(docs.Spec())))

	return apiRoutes
}
