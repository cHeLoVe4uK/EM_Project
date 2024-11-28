package v1

import (
	"encoding/json"
	"github.com/cHeLoVe4uK/EM_Project/internal/schemas"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) initUserHandler(r *httprouter.Router) {
	r.POST("/api/v1/user/login", h.loginUser)
	r.POST("/api/v1/user/register", h.registerUser)
	r.POST("/api/v1/user/logout", h.authenticated(h.logoutUser))
}

// loginUser godoc
// @Tags         User API
// @Summary      Login
// @Description  Login
// @Accept       json
// @Produce      json
// @Param Audio body schemas.RequestLoginUser true "Email and password"
// @Success      200
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /user/login [post]
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userLogin := schemas.RequestLoginUser{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call userService, return jwt

	http.SetCookie(w, &http.Cookie{
		Name:  "_session",
		Value: "example-token",
		Path:  "/",
	})
	writeResponse(w, 200, nil)
}

// registerUser godoc
// @Tags         User API
// @Summary      Register
// @Description  Register
// @Accept       json
// @Produce      json
// @Param Audio body schemas.RequestRegisterUser true "Username, Email, password"
// @Success      200
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /user/register [post]
func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userRegister := schemas.RequestRegisterUser{}

	// parse body
	err := json.NewDecoder(r.Body).Decode(&userRegister)
	if err != nil {
		writeResponseErr(w, 400, err, "error on parse body")
		return
	}

	// todo call userService, return jwt

	http.SetCookie(w, &http.Cookie{
		Name:  "_session",
		Value: "example-token",
		Path:  "/",
	})
	writeResponse(w, 200, nil)
}

// logoutUser godoc
// @Tags         User API
// @Summary      Logout
// @Description  Logout
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  ErrResponse
// @Failure      500  {object}	ErrResponse
// @Router       /user/logout [post]
func (h *Handler) logoutUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	// clean _session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "_session",
		Value: "",
		Path:  "/",
	})

	writeResponse(w, 200, nil)
}
