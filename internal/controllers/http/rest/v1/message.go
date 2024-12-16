package v1

import (
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/cHeLoVe4uK/EM_Project/internal/services/message"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
)

// @Summary		Update message content
// @Description	Updates message content
// @Tags			Chats
// @Produce		json
// @Param			chat_id	path		string	true	"Chat ID"
// @Param			msg_id	path		string	true	"Message ID"
// @Success		200	{object}	object
// @Failure		400	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	object
// @Router			/api/v1/chats/{chat_id}/messages/{msg_id} [patch]
func (a *API) UpdateMessage(c echo.Context) error {
	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "UpdateMessage"),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	var req UpdateChatRequest

	if err := c.Bind(&req); err != nil {
		log.Warn("bind request", logging.Err(err))

		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	chatID := c.Param("chat_id")
	if chatID == "" {
		log.Warn("chat id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "chat id is empty")
	}

	msgID := c.Param("msg_id")
	if msgID == "" {
		log.Warn("message id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "message id is empty")
	}

	userID := c.Get("user_id")
	if userID == nil {
		log.Warn("user id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "user id is empty")
	}

	msg := models.Message{
		ID:       msgID,
		ChatID:   chatID,
		AuthorID: userID.(string),
		Content:  req.Content,
	}

	log = logging.WithAttrs(
		ctx,
		logging.String("message_id", msg.ID),
		logging.String("chat_id", msg.ChatID),
		logging.String("user_id", msg.AuthorID),
	)

	ctx = logging.ContextWithLogger(ctx, log)

	err := a.chatService.UpdateMessage(ctx, msg)
	switch err {
	case nil:
		return nil
	case message.ErrMessageNotFound:
		return echo.NewHTTPError(http.StatusNotFound, err)
	case message.ErrNotAllowed:
		return echo.NewHTTPError(http.StatusForbidden, err)
	default:
		return err
	}
}

// @Summary		Delete message
// @Description	Deletes message from chat
// @Tags			Chats
// @Produce		json
// @Param			chat_id	path		string	true	"Chat ID"
// @Param			msg_id	path		string	true	"Message ID"
// @Success		200	{object}	object
// @Failure		400	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	object
// @Router			/api/v1/chats/{chat_id}/messages/{msg_id} [delete]
func (a *API) DeleteMessage(c echo.Context) error {
	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "DeleteMessage"),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	chatID := c.Param("chat_id")
	if chatID == "" {
		log.Warn("chat id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "chat id is empty")
	}

	msgID := c.Param("msg_id")
	if msgID == "" {
		log.Warn("message id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "message id is empty")
	}

	userID := c.Get("user_id")
	if userID == nil {
		log.Warn("user id is empty")

		return echo.NewHTTPError(http.StatusBadRequest, "user id is empty")
	}

	msg := models.Message{
		ID:       msgID,
		ChatID:   chatID,
		AuthorID: userID.(string),
	}

	log = logging.WithAttrs(
		ctx,
		logging.String("message_id", msg.ID),
		logging.String("chat_id", msg.ChatID),
		logging.String("user_id", msg.AuthorID),
	)

	ctx = logging.ContextWithLogger(ctx, log)

	err := a.chatService.DeleteMessage(ctx, msg)
	switch err {
	case nil:
		return nil
	case message.ErrMessageNotFound:
		return echo.NewHTTPError(http.StatusNotFound, err)
	case message.ErrNotAllowed:
		return echo.NewHTTPError(http.StatusForbidden, err)
	default:
		return err
	}
}
