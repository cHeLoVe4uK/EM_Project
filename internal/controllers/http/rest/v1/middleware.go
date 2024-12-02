package v1

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *API) corsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request().Method == http.MethodOptions {
			return c.JSON(http.StatusOK, nil)
		}

		return next(c)
	}
}

func (a *API) loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		slog.Info("request", slog.String("method", c.Request().Method), slog.String("path", c.Request().URL.Path))

		err := next(c)
		if err != nil {
			slog.Error("request error", slog.Any("error", err))
			return err
		}

		return nil
	}
}
