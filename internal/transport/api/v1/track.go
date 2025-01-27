package v1

import (
	"encoding/json"
	"music_storage/internal/domain"
	"net/http"
	"net/url"
)

func (h *Handler) initTrackRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /list", h.trackList)

	trackRoutes := http.NewServeMux()
	trackRoutes.Handle("/track/", http.StripPrefix("/track", routes))

	return trackRoutes
}

// ??? Rename or remove this
type TrackResponse struct {
	Data []domain.Track `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary  Get tracks with filter
// @Tags track
// @Description get tracks with filter by query get params. Get parameters is optional
// @ModuleID trackList
// @Accept json
// @Produce json
// @Success 200 {object} TrackResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /track/list [get]
func (h *Handler) trackList(w http.ResponseWriter, r *http.Request) {
	// todo. json response when error
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ErrorResponse{Error: "error: incorrect query params"})
		w.Write(json)

		return
	}

	tracks, err := h.services.Track.List(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ErrorResponse{Error: "error: getting tracks"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(TrackResponse{Data: tracks})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: creating json"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
