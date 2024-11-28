package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

//	@Summary		Create chat
//	@Description	Creates new chat, runs in background and returns chat ID
//	@Tags			Chats
//	@Accept			json
//	@Produce		json
//	@Param			chat	body		CreateChatRequest	true	"Chat name"
//	@Success		200		{object}	CreateChatResponse
//	@Failure		400		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/v1/chats [post]
func (a *API) CreateChat(w http.ResponseWriter, r *http.Request) {

	var req CreateChatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	chat := models.Chat{
		Name: req.Name,
	}

	slog.Debug("creating chat", slog.String("chat_name", chat.Name))

	chatId, err := a.chatService.CreateChat(r.Context(), chat)
	if err != nil {

		slog.Error("creating chat", slog.Any("error", err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	slog.Debug("chat created", slog.String("chat_id", chatId))

	resp := CreateChatResponse{
		ID: chatId,
	}

	data, err := json.Marshal(resp)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Write(data)
}
