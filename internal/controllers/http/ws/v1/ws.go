package v1

import (
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
// @Param			id path string true "Chat ID"
// @Failure		422		{object}	object
// @Failure		500		{object}	object
// @Router			/api/v1/chats/{id}/connect [get]
func (a *API) ConnectChat(w http.ResponseWriter, r *http.Request) {

	slog.Debug("decoding request")

	chatID := r.PathValue("id")

	if chatID == "" {
		slog.Error("chat id is empty")

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Ignore for now, just testing
	cookie, _ := r.Cookie("username")

	username := cookie.Value
	if username == "" {
		username = gofakeit.Username()
	}

	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	log := slog.With(
		slog.String("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("chat_id", chatID),
	)

	log.Debug("connecting to chat")

	if err := a.chatService.ConnectByID(w, r, chatID, user); err != nil {
		log.Error("connecting to chat", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info(
		"user connected",
	)
}
