package v1

import (
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
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

		return next(c)
	}
}

func (a *API) recoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log := logging.WithAttrs(
					c.Request().Context(),
					logging.String("op", "RecoverMiddleware"),
				)

				log.Error("panic", logging.Any("error", err))
			}
		}()

		return next(c)
	}
}

func (a *API) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access, err := c.Cookie("access_token")
		if err != nil {
			_, err := c.Cookie("refresh_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			//tokens, err := a.authService.Refresh(c.Request().Context(), refresh.Value)
			//if err != nil {
			//return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			//}

			// How parse jwt????

			// Save in cookie

			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		tokens := models.Tokens{
			AccessToken: access.Value,
		}

		claims, err := a.authService.Authenticate(c.Request().Context(), tokens)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		return next(c)
	}
}
