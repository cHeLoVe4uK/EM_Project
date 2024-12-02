package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
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
		reqID := uuid.New().String()

		log := logging.WithAttrs(
			c.Request().Context(),
			logging.String("request_id", reqID),
		)

		ctx := logging.ContextWithLogger(c.Request().Context(), log)

		log.Info(
			"request",
			logging.String("method", c.Request().Method),
			logging.String("path", c.Request().URL.Path),
		)

		r := c.Request().WithContext(ctx)

		c.SetRequest(r)

		err := next(c)
		if err != nil {
			log.Error("request error", logging.Err(err))
			return err
		}

		return nil
	}
}
