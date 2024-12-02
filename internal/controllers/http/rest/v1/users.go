package v1

import (
	"errors"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	userrepo "github.com/cHeLoVe4uK/EM_Project/internal/repo/userRepo"
	"github.com/labstack/echo/v4"
	"github.com/meraiku/logging"
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

	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "CreateUser"),
	)

	var req CreateUserRequest

	log.Debug("binding request")

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	log = logging.WithAttrs(
		c.Request().Context(),
		logging.String("email", req.Email),
		logging.String("username", req.Username),
	)

	log.Debug("creating user model")

	user, err := models.NewUser(
		req.Email,
		req.Username,
		req.Password,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	id, err := a.userService.Create(ctx, user)
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
	log := logging.WithAttrs(
		c.Request().Context(),
		logging.String("op", "LoginUser"),
	)

	var req LoginUserRequest

	log.Debug("binding request")

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	log = logging.WithAttrs(
		c.Request().Context(),
		logging.String("email", req.Email),
	)

	log.Debug("creating user model")

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := logging.ContextWithLogger(c.Request().Context(), log)

	token, err := a.userService.Login(ctx, user)
	if err != nil {
		return err
	}

	var res LoginUserResponse

	res.Token = token.Token

	return c.JSON(http.StatusOK, res)
}
