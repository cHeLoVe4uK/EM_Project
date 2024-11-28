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

// connectToChat godoc
// @Tags         Chat API
// @Summary      Connect to chat
// @Description  Connect to chat
// @Accept       json
// @Produce      json
// @Param Audio body schemas.RequestConnectToChat true "Chat ID"
// @Success      200
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /chat/connect [post]
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

// createChat godoc
// @Tags         Chat API
// @Summary      Create chat
// @Description  Create chat
// @Accept       json
// @Produce      json
// @Param Audio body schemas.RequestCreateChat true "Chat name"
// @Success      200  {object}  schemas.ResponseCreateChat
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /chat [post]
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

// deleteChat godoc
// @Tags         Chat API
// @Summary      Delete chat
// @Description  Delete chat
// @Accept       json
// @Produce      json
// @Param Audio body schemas.RequestDeleteChat true "Chat ID"
// @Success      200
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /chat [delete]
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
