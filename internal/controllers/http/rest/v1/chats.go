package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	chatrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/chatRepo"
	"github.com/labstack/echo/v4"
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

	var req CreateChatRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	chat := models.Chat{
		Name: req.Name,
	}

	slog.Debug("creating chat", slog.String("chat_name", chat.Name))

	chatId, err := a.chatService.CreateChat(c.Request().Context(), chat)
	if err != nil {

		slog.Error("creating chat", slog.Any("error", err))

		return err
	}

	slog.Debug("chat created", slog.String("chat_id", chatId))

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

	chats, err := a.chatService.GetAllChats(c.Request().Context())
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

	chats, err := a.chatService.GetActiveChats(c.Request().Context())
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

	msgs, err := a.chatService.GetMessages(c.Request().Context(), chatID)
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
