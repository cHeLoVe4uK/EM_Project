package v1

import (
	"net/http"
	"strings"

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
		log := logging.WithAttrs(
			c.Request().Context(),
			logging.String("op", "AuthMiddleware"),
		)

		log.Debug("auth middleware executed")

		tokenHeader := c.Request().Header.Get("Authorization")

		log.Debug("token header", logging.String("token", tokenHeader))

		reqToken := strings.Split(tokenHeader, "Bearer ")
		if len(reqToken) != 2 {
			log.Debug("token is empty")

			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		token := reqToken[1]

		tokens := models.Tokens{
			AccessToken: token,
		}

		log.Debug("got token, trying to authenticate")

		claims, err := a.authService.Authenticate(c.Request().Context(), tokens)
		if err != nil {

			log.Warn("authenticate user", logging.Any("error", err))

			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		log.Debug(
			"user authenticated",
			logging.String("user_id", claims.UserID),
			logging.String("username", claims.Username),
		)

		return next(c)
	}
}
