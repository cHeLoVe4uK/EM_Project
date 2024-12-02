package v1

import (
	"log/slog"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
)

// @Summary		Upgrade http connection
// @Description	Upgrades http connection to websocket
// @Tags			Chats
// @Produce		json
// @Param			id	path		string	true	"Chat ID"
// @Failure		422	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats/{id}/connect [get]
func (a *API) ConnectChat(c echo.Context) error {

	log := logging.WithAttrs(
		c.Request().Context(),
		slog.String("op", "ConnectChat"),
	)

	log.Debug("decoding request")

	chatID := c.Param("id")

	if chatID == "" {
		log.Error("chat id is empty")

		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	username := gofakeit.Username()

	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	if err := a.chatService.ConnectByID(c.Response().Writer, c.Request(), chatID, user); err != nil {
		log.Error("connecting to chat", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, err)
	}

	return nil
}
