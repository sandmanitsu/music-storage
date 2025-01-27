package v1

import (
	"music_storage/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", h.initTrackRoutes()))

	return router
}
