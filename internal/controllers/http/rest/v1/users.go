package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo"
	"github.com/labstack/echo/v4"
)

// @Summary		Create New User
// @Description	Creates nes User, return his ID
// @Tags			Users
// @Produce		json
// @Param			user	body		CreateUserRequest	true	"User data"
// @Success		201		{object}	CreateUserResponse
// @Failure		422		{object}	object
// @Failure		500		{object}	object
// @Router			/api/v1/users [post]
func (a *API) CreateUser(c echo.Context) error {

	var req CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	user, err := models.NewUser(
		req.Email,
		req.Username,
		req.Password,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	id, err := a.userService.Create(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
		}
		return err
	}

	var res CreateUserResponse

	res.ID = id

	return c.JSON(http.StatusCreated, res)
}

// @Summary		Create New User
// @Description	Creates nes User, return his ID
// @Tags			Users
// @Produce		json
// @Param			user	body		LoginUserRequest	true	"User login data"
// @Success		200		{object}	LoginUserResponse
// @Failure		422		{object}	object
// @Failure		500		{object}	object
// @Router			/api/v1/users/login [post]
func (a *API) LoginUser(c echo.Context) error {
	log := slog.With(
		slog.String("op", "LoginUser"),
	)

	var req LoginUserRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	log.Debug("login user", slog.String("email", user.Email))

	token, err := a.userService.Login(c.Request().Context(), user)
	if err != nil {
		log.Error(
			"failed to login user",
			slog.Any("error", err),
		)
		return err
	}

	log.Debug("user logged in", slog.String("token", token.Token))

	var res LoginUserResponse

	res.Token = token.Token

	return c.JSON(http.StatusOK, res)
}
