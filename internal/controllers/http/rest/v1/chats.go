package v1

import (
	"errors"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	chatrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
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
func (a *API) CreateChat(c echo.Context) error {

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "CreateChat"),
	)

	var req CreateChatRequest

	log.Debug("binding request")

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	chat, err := models.NewChat(req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	chatId, err := a.chatService.CreateChat(c.Request().Context(), chat)
	if err != nil {
		return err
	}

	log.Debug("chat created", logging.String("chat_id", chatId))

	resp := CreateChatResponse{
		ID: chatId,
	}

	return c.JSON(http.StatusOK, resp)
}

// @Summary		Get all chats
// @Description	Prints all chats id and name
// @Tags			Chats
// @Produce		json
// @Success		200	{array}		Chat
// @Failure		400	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats [get]
func (a *API) GetAllChats(c echo.Context) error {

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "GetAllChats"),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	chats, err := a.chatService.GetAllChats(ctx)
	if err != nil {
		return err
	}

	out := make([]Chat, len(chats))

	for i, chat := range chats {
		out[i].ID = chat.ID
		out[i].Name = chat.Name
	}

	return c.JSON(http.StatusOK, out)
}

// @Summary		Get all active chats
// @Description	Prints all active chats id and name
// @Tags			Chats
// @Produce		json
// @Success		200	{array}		Chat
// @Failure		400	{object}	object
// @Failure		500	{object}	object
// @Router			/api/v1/chats/active [get]
func (a *API) GetAllActiveChats(c echo.Context) error {

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "GetAllActiveChats"),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	chats, err := a.chatService.GetActiveChats(ctx)
	if err != nil {
		return err
	}

	out := make([]Chat, len(chats))

	for i, chat := range chats {
		out[i].ID = chat.ID
		out[i].Name = chat.Name
	}

	return c.JSON(http.StatusOK, out)
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
func (a *API) GetChatMessages(c echo.Context) error {

	chatID := c.Param("id")

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "GetChatMessages"),
		logging.String("chat_id", chatID),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	msgs, err := a.chatService.GetMessages(ctx, chatID)
	if err != nil {
		if errors.Is(err, chatrepo.ErrChatNotFound) {
			return c.JSON(http.StatusNotFound, err)
		}

		return err
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

	return c.JSON(http.StatusOK, out)
}
