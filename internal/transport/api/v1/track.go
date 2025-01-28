package v1

import (
	"encoding/json"
	"io"
	"music_storage/internal/domain"
	"net/http"
	"net/url"
	"strconv"
)

func (h *Handler) initTrackRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /list", h.list)
	routes.HandleFunc("GET /text", h.text)
	routes.HandleFunc("DELETE /delete", h.delete)

	trackRoutes := http.NewServeMux()
	trackRoutes.Handle("/track/", http.StripPrefix("/track", routes))

	return trackRoutes
}

type ListResponse struct {
	Data  []domain.Track `json:"data"`
	Error string         `json:"error"`
}

// @Summary  Get list of song
// @Tags track
// @Description get tracks with filter by query get params. Get parameters is optional
// @ModuleID list
// @Accept json
// @Produce json
// @Param        id				path	int		false  "song id"
// @Param        group_name		path	string	false  "group name"
// @Param        song			path	string	false  "song name"
// @Param        text			path	string	false  "gong text"
// @Param        realise_date	path	string	false  "realise date"
// @Param        limit   		path	int		false  "limit"
// @Param        offset   		path	int		false  "offset"
// @Success 200 {object} ListResponse
// @Failure 400 {object} ListResponse
// @Failure 500 {object} ListResponse
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

type TextResponse struct {
	Text  []string `json:"text"`
	Error string   `json:"error"`
}

// @Summary Get song text
// @Tags track
// @Description get song text with id separated by choruses
// @ModuleID text
// @Accept json
// @Produce json
// @Param        id		path	int		true	"song id"
// @Success 200 {object} TextResponse
// @Failure 400,404 {object} TextResponse
// @Failure 500 {object} TextResponse
// @Failure default {object} TextResponse
// @Router /track/text [get]
func (h *Handler) text(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(TextResponse{Text: nil, Error: "error: empty or incorrent id"})
		w.Write(json)

		return
	}

	chorus, err := h.services.Track.Text(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(TextResponse{Text: nil, Error: err.Error()})
		w.Write(json)

		return
	}

	json, err := json.Marshal(TextResponse{Text: chorus, Error: ""})
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

// @Summary Delete song
// @Tags track
// @Description Delete song from storage by id
// @ModuleID delete
// @Accept json
// @Produce json
// @Param        id		path	int		true	"song id"
// @Success 200 {object} DeleteResponse
// @Failure 400 {object} DeleteResponse
// @Failure 500 {object} DeleteResponse
// @Router /track/delete [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
