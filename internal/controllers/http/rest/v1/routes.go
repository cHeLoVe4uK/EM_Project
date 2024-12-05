package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (a *API) routes(e *echo.Echo) {

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(a.loggingMiddleware)
	e.Use(a.corsMiddleware)
	e.Use(a.recoverMiddleware)

	api := e.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	chats := api.Group("/chats")

	chats.Use(a.authMiddleware)

	chats.GET("", a.GetAllChats)
	chats.GET("/active", a.GetAllActiveChats)
	chats.POST("", a.CreateChat)

	chats.GET("/:id/connect", a.ConnectChat)

	chats.GET("/:id/messages", a.GetChatMessages)

	users := api.Group("/users")

	users.POST("", a.CreateUser)

	users.POST("/login", a.LoginUser)

}
