package v1

import (
	"encoding/json"
	"github.com/cHeLoVe4uK/EM_Project/internal/schemas"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) initMsgHandler(r *httprouter.Router) {
	r.DELETE("/api/v1/message", h.authenticated(h.deleteMsg))
	r.PATCH("/api/v1/message", h.authenticated(h.updateMsg))
}

func (h *Handler) deleteMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deleteMsg := schemas.RequestDeleteMsg{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&deleteMsg)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call msgService

	writeResponseErr(w, 501, nil, "deleteMsg: Not implemented")
}

func (h *Handler) updateMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	updateMsg := schemas.RequestUpdateMsg{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&updateMsg)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call msgService

	writeResponseErr(w, 501, nil, "updateMsg: Not implemented")
}
