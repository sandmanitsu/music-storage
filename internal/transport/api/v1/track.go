package v1

import (
	"encoding/json"
	"io"
	"music_storage/internal/domain"
	"music_storage/internal/service"
	"net/http"
	"net/url"
	"strconv"
)

const (
	listMsg     = "request: getting list of song"
	songTextMsg = "request: getting gong text"
	deleteMsg   = "request: delete song"
	updateMsg   = "request: update song"
	addMsg      = "request: add song"
)

func (h *Handler) initTrackRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /list", h.list)
	routes.HandleFunc("GET /text", h.text)
	routes.HandleFunc("DELETE /delete", h.delete)
	routes.HandleFunc("POST /update", h.update)
	routes.HandleFunc("POST /add", h.add)

	trackRoutes := http.NewServeMux()
	trackRoutes.Handle("/track/", http.StripPrefix("/track", routes))

	return trackRoutes
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
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
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ListResponse{Data: nil, Error: "error: incorrect query params"})
		w.Write(json)

		return
	}

	tracks, err := h.services.Track.List(params)
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(ListResponse{Data: nil, Error: "error: getting tracks"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(ListResponse{Data: tracks, Error: ""})
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())
		h.creatingJsonErr(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	h.logger.InfoAPI(listMsg, http.StatusOK, r.URL.String(), "")
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
		h.logger.InfoAPI(songTextMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(TextResponse{Text: nil, Error: "error: empty or incorrent id"})
		w.Write(json)

		return
	}

	chorus, err := h.services.Track.Text(id)
	if err != nil {
		h.logger.InfoAPI(songTextMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(TextResponse{Text: nil, Error: err.Error()})
		w.Write(json)

		return
	}

	json, err := json.Marshal(TextResponse{Text: chorus, Error: ""})
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())
		h.creatingJsonErr(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	h.logger.InfoAPI(songTextMsg, http.StatusOK, r.URL.String(), "")
}

type DeleteRequest struct {
	ID int `json:"id"`
}

// @Summary Delete song
// @Tags track
// @Description Delete song from storage by id
// @ModuleID delete
// @Accept json
// @Produce json
// @Param        id		path	int		true	"song id"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /track/delete [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.InfoAPI(deleteMsg, http.StatusOK, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: reading request body"})
		w.Write(json)

		return
	}

	var input DeleteRequest
	err = json.Unmarshal(body, &input)
	if err != nil {
		h.logger.InfoAPI(deleteMsg, http.StatusOK, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: unmarshaling json"})
		w.Write(json)

		return
	}

	err = h.services.Track.Delete(input.ID)
	if err != nil {
		h.logger.InfoAPI(deleteMsg, http.StatusOK, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: db delete error"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(Response{Status: "success", Error: ""})
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())
		h.creatingJsonErr(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	h.logger.InfoAPI(deleteMsg, http.StatusOK, r.URL.String(), "")
}

// @Summary Update song
// @Tags track
// @Description Update song by id
// @ModuleID update
// @Accept json
// @Produce json
// @param params body service.TrackInput true "id is required param"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /track/update [post]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: reading request body"})
		w.Write(json)

		return
	}

	var input service.TrackInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: bad params or empty id"})
		w.Write(json)

		return
	}

	err = h.services.Track.Update(input)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: update data"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(Response{Status: "success", Error: ""})
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())
		h.creatingJsonErr(w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	h.logger.InfoAPI(updateMsg, http.StatusOK, r.URL.String(), "")
}

// @Summary add song
// @Tags track
// @Description adding song to storage
// @ModuleID add
// @Accept json
// @Produce json
// @param params body service.TrackAddInput true "group and song name is required"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /track/add [post]
func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: reading request body"})
		w.Write(json)

		return
	}

	var input service.TrackAddInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: bad params or empty id"})
		w.Write(json)

		return
	}

	err = h.services.Track.Add(input)
	if err != nil {
		h.logger.InfoAPI(updateMsg, http.StatusBadRequest, r.URL.String(), err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(Response{Status: "failed", Error: "error: insert data"})
		w.Write(json)

		return
	}

	json, err := json.Marshal(Response{Status: "success", Error: ""})
	if err != nil {
		h.logger.InfoAPI(listMsg, http.StatusBadRequest, r.URL.String(), err.Error())
		h.creatingJsonErr(w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	h.logger.InfoAPI(updateMsg, http.StatusOK, r.URL.String(), "")
}
