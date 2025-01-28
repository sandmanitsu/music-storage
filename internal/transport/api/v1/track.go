package v1

import (
	"encoding/json"
	"io"
	"music_storage/internal/domain"
	"net/http"
	"net/url"
)

func (h *Handler) initTrackRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /list", h.list)
	routes.HandleFunc("DELETE /delete", h.delete)

	trackRoutes := http.NewServeMux()
	trackRoutes.Handle("/track/", http.StripPrefix("/track", routes))

	return trackRoutes
}

type ListResponse struct {
	Data  []domain.Track `json:"data"`
	Error string         `json:"error"`
}

// @Summary  Get tracks with filter
// @Tags track
// @Description get tracks with filter by query get params. Get parameters is optional
// @ModuleID trackList
// @Accept json
// @Produce json
// @Success 200 {object} TrackResponse
// @Failure 400,404 {object} ListResponse
// @Failure 500 {object} ListResponse
// @Failure default {object} ListResponse
// @Router /track/list [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ListResponse{Data: nil, Error: "error: incorrect query params"})
		w.Write(json)

		return
	}

	tracks, err := h.services.Track.List(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ListResponse{Data: nil, Error: "error: getting tracks"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(ListResponse{Data: tracks, Error: ""})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"creating json\"}"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

type DeleteRequest struct {
	ID int `json:"id"`
}

type DeleteResponse struct {
	Status string `json:"status"`
	Error  string
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(DeleteResponse{Status: "failed", Error: "error: reading request body"})
		w.Write(json)

		return
	}

	var input DeleteRequest
	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(DeleteResponse{Status: "failed", Error: "error: unmarshaling json"})
		w.Write(json)

		return
	}

	err = h.services.Track.Delete(input.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(DeleteResponse{Status: "failed", Error: "error: db delete error"})
		w.Write(json)

		return
	}

	w.WriteHeader(http.StatusOK)
	json, _ := json.Marshal(DeleteResponse{Status: "success", Error: ""})
	w.Write(json)
}
