package v1

import (
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
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
		logging.String("op", "ConnectChat"),
	)

	log.Debug("decoding request")

	chatID := c.Param("id")
	if chatID == "" {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "chat id is empty")
	}

	uid := c.Request().Context().Value("user_id")
	if uid == nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "user id is empty")
	}
	username := c.Request().Context().Value("username")
	if username == nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "username is empty")
	}

	user := &models.User{
		ID:       uid.(string),
		Username: username.(string),
	}

	if err := a.chatService.ConnectByID(c.Response().Writer, c.Request(), chatID, user); err != nil {
		return err
	}

	return nil
}
