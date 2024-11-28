package v1

import (
	"encoding/json"
	"net/http"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
)

//	@Summary		Create New User
//	@Description	Creates nes User, return his ID
//	@Tags			Users
//	@Produce		json
//	@Param			user	body		CreateUserRequest	true	"User data"
//	@Success		201		{object}	CreateUserResponse
//	@Failure		422		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/v1/users [post]
func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	user, err := models.NewUser(
		req.Email,
		req.Username,
		req.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, err := a.userService.Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res CreateUserResponse

	res.ID = id

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Write(data)
}

//	@Summary		Create New User
//	@Description	Creates nes User, return his ID
//	@Tags			Users
//	@Produce		json
//	@Param			user	body		LoginUserRequest	true	"User login data"
//	@Success		200		{object}	LoginUserResponse
//	@Failure		422		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/v1/users/login [post]
func (a *API) LoginUser(w http.ResponseWriter, r *http.Request) {

	var req LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := a.userService.Login(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res LoginUserResponse

	res.Token = token.Token

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(data)

}
