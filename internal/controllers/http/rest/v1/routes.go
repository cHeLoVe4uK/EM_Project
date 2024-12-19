package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (a *API) routes(e *echo.Echo) {

	e.GET("api/swagger/*", echoSwagger.WrapHandler)

	e.Use(a.loggingMiddleware)
	e.Use(a.corsMiddleware)
	e.Use(a.recoverMiddleware)

	api := e.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	ws := api.Group("/chats")

	ws.Use(a.wsAuthMiddleware)

	ws.GET("/:id/connect", a.ConnectChat)

	chats := api.Group("/chats")

	chats.Use(a.authMiddleware)

	chats.GET("", a.GetAllChats)
	chats.GET("/active", a.GetAllActiveChats)
	chats.POST("", a.CreateChat)

	chats.GET("/:id/messages", a.GetChatMessages)
	chats.PATCH("/:chat_id/messages/:msg_id", a.UpdateMessage)
	chats.DELETE("/:chat_id/messages/:msg_id", a.DeleteMessage)

	users := api.Group("/users")

	users.POST("", a.CreateUser)

	users.POST("/login", a.LoginUser)

}
