package v1

import (
	"net/http"
)

func (h *Handler) initTrackRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /list", h.trackList)

	trackRoutes := http.NewServeMux()
	trackRoutes.Handle("/track/", http.StripPrefix("/track", routes))

	return trackRoutes
}

func (h *Handler) trackList(w http.ResponseWriter, r *http.Request) {
	h.services.Track.List()
	w.Write([]byte("track list......."))
}
