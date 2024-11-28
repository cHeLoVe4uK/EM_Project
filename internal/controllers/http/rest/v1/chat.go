package v1

import (
	"encoding/json"
	"github.com/cHeLoVe4uK/EM_Project/internal/schemas"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) initChatHandler(r *httprouter.Router) {
	r.POST("/api/v1/chat/connect", h.authenticated(h.connectToChat))
	r.POST("/api/v1/chat", h.authenticated(h.createChat))
	r.DELETE("/api/v1/chat", h.authenticated(h.deleteChat))
}

func (h *Handler) connectToChat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	connectToChat := schemas.RequestConnectToChat{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&connectToChat)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call chatService

	writeResponseErr(w, 501, nil, "connectToChat: Not implemented")
}

func (h *Handler) createChat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	connectToChat := schemas.RequestCreateChat{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&connectToChat)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call chatService, return chat_id

	// resp := &schemas.ResponseCreateChat{}
	// writeResponse(w, 200, resp)
	writeResponseErr(w, 501, nil, "createChat: Not implemented")
}

func (h *Handler) deleteChat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	deleteChat := schemas.RequestDeleteChat{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&deleteChat)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call chatService

	writeResponseErr(w, 501, nil, "deleteChat: Not implemented")
}
