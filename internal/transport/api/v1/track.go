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

func (h *Handler) trackList(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: parsing query params"))

		return
	}

	tracks, err := h.services.Track.List(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: getting tracks"))

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
