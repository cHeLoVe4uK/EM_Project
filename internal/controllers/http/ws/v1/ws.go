package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
)

// @Summary		Upgrade http connection
// @Description	Upgrades http connection to websocket
// @Tags			Chats
// @Produce		json
// @Param			chat body		JoinChatRequest	true	"Chat data"
// @Failure		422		{object}	object
// @Failure		500		{object}	object
// @Router			/api/v1/chats/connect [get]
func (a *API) ConnectChat(w http.ResponseWriter, r *http.Request) {

	var req JoinChatRequest

	slog.Debug("decoding request")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Get username from cookie

	username := gofakeit.Username()

	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	slog.With(
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("chat_id", req.ChatID),
	)

	slog.Debug("connecting to chat")

	if err := a.chatService.ConnectByID(w, r, req.ChatID, user); err != nil {
		slog.Error("connecting to chat", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info(
		"user connected",
	)
}
