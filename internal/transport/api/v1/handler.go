package v1

import (
	"music_storage/internal/logger"
	"music_storage/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
	logger   *logger.Logger
}

func NewHandler(logger *logger.Logger, services *service.Service) *Handler {
	return &Handler{services: services, logger: logger}
}

func (h *Handler) Init() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", h.initTrackRoutes()))

	return router
}

func (h *Handler) creatingJsonErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\": \"creating json\"}"))
}
