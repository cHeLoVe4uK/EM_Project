package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	chatrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo"
)

// @Summary		Create chat
// @Description	Creates new chat, runs in background and returns chat ID
// @Tags			Chats
// @Accept			json
// @Produce		json
// @Param			chat	body		CreateChatRequest	true	"Chat name"
// @Success		200		{object}	CreateChatResponse
// @Failure		400		{object}	object
// @Failure		500		{object}	object
// @Router			/api/v1/chats [post]
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

// @Summary		Get all chats
// @Description	Prints all chats id and name
// @Tags			Chats
// @Produce		json
// @Success		200	{array}		Chat
// @Failure		400	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats [get]
func (a *API) GetAllChats(w http.ResponseWriter, r *http.Request) {

	chats, err := a.chatService.GetAllChats(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := make([]Chat, len(chats))

	for i, chat := range chats {
		out[i].ID = chat.ID
		out[i].Name = chat.Name
	}

	data, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(data)
}

// @Summary		Get all active chats
// @Description	Prints all active chats id and name
// @Tags			Chats
// @Produce		json
// @Success		200	{array}		Chat
// @Failure		400	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats/active [get]
func (a *API) GetAllActiveChats(w http.ResponseWriter, r *http.Request) {

	chats, err := a.chatService.GetActiveChats(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := make([]Chat, len(chats))

	for i, chat := range chats {
		out[i].ID = chat.ID
		out[i].Name = chat.Name
	}

	data, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(data)
}

// @Summary		Get chat message history
// @Description	Prints 100 messages from chat
// @Tags			Chats
// @Produce		json
// @Param			id	path		string	true	"Chat ID"
// @Success		200	{array}		Message
// @Failure		400	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats/{id}/messages [get]
func (a *API) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("id")

	msgs, err := a.chatService.GetMessages(r.Context(), chatID)
	if err != nil {
		if errors.Is(err, chatrepo.ErrChatNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out := make([]Message, len(msgs))

	for i, msg := range msgs {
		out[i].ID = msg.ID
		out[i].Author = msg.Author
		out[i].ChatID = msg.ChatID
		out[i].Content = msg.Content
		out[i].CreatedAt = msg.CreatedAt
		out[i].IsEdited = msg.IsEdited
	}

	data, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(data)
}
