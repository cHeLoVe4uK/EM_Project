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

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userLogin := schemas.RequestUserLogin{}

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

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userRegister := schemas.RequestUserRegister{}

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

func (h *Handler) logoutUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	// clean _session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "_session",
		Value: "",
		Path:  "/",
	})

	writeResponse(w, 200, nil)
}
