package v1

import (
	"fmt"
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

func (h *Handler) trackList(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: parsing query params"))

		return
	}

	fmt.Println(params)

	h.services.Track.List(params)
	w.Write([]byte("track list......."))
}
