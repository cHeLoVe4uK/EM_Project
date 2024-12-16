package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userService "github.com/cHeLoVe4uK/EM_Project/internal/services/user"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
)

// @Summary		Create New User
// @Description	Creates nes User, return his ID
// @Tags			Users
// @Produce		json
// @Param			user	body		CreateUserRequest	true	"User data"
// @Success		201		{object}	CreateUserResponse
// @Failure		422		{object}	HTTPError
// @Failure		500		{object}	object
// @Router			/api/v1/users [post]
func (a *API) CreateUser(c echo.Context) error {

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "CreateUser"),
	)

	var req CreateUserRequest

	log.Debug("binding request")

	if err := c.Bind(&req); err != nil {
		log.Warn("binding request", logging.Err(err))

		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	log = logging.WithAttrs(
		c.Request().Context(),
		logging.String("email", req.Email),
		logging.String("username", req.Username),
	)

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	log.Debug("creating user model")

	user, err := models.NewUser(
		req.Email,
		req.Username,
		req.Password,
	)
	if err != nil {
		log.Warn("creating user model", logging.Err(err))

		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	id, err := a.userService.Register(ctx, user)
	if err != nil {
		if errors.Is(err, userService.ErrUserExists) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return err
	}

	var res CreateUserResponse

	res.ID = id

	return c.JSON(http.StatusCreated, res)
}

// @Summary		Login User
// @Description	Login User, returns token
// @Tags			Users
// @Produce		json
// @Param			user	body		LoginUserRequest	true	"User login data"
// @Success		200		{object}	LoginUserResponse
// @Failure		422		{object}	HTTPError
// @Failure		500		{object}	object
// @Router			/api/v1/users/login [post]
func (a *API) LoginUser(c echo.Context) error {
	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "LoginUser"),
	)

	var req LoginUserRequest

	log.Debug("binding request")

	if err := c.Bind(&req); err != nil {
		log.Warn("binding request", logging.Err(err))

		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	log = logging.WithAttrs(
		c.Request().Context(),
		logging.String("email", req.Email),
	)

	log.Debug("creating user model")

	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		log.Warn("parse email", logging.Err(err))

		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Errorf("invalid email: %w", err))
	}

	user := models.User{
		Email:    email.String(),
		Password: req.Password,
	}

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	tokens, err := a.userService.Login(ctx, user)
	switch err {
	case nil:
		break
	case userService.ErrUserNotFound:
		return echo.NewHTTPError(http.StatusNotFound, err)
	case userService.ErrInvalidPassword:
		return echo.NewHTTPError(http.StatusBadRequest, err)
	default:
		return err
	}

	var res LoginUserResponse

	res.AccessToken = tokens.AccessToken
	res.RefreshToken = tokens.RefreshToken

	return c.JSON(http.StatusOK, res)
}
