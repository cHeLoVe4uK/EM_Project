package v1

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// authenticated
// set user_id, username into context
func (h *Handler) authenticated(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// parse cookie
		_, err := h.getJWT(r.Header)
		if err != nil {
			writeResponseErr(w, 500, err, "error on parse cookie")
			return
		}

		// todo authService call, parse jwt payload, set payload to ctx

		// sets user_id, username in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", "some_user_id")
		ctx = context.WithValue(ctx, "username", "some_username")
		r = r.WithContext(ctx)

		// next handler...
		handler(w, r, ps)
	}
}

// getJWT
// parse cookie header and return _session value
func (h *Handler) getJWT(header http.Header) (string, error) {
	cookieRaw := header.Get("Cookie")

	cookies, err := http.ParseCookie(cookieRaw)
	if err != nil {
		return "", err
	}

	for _, c := range cookies {
		if c.Name == "_session" {
			return c.Value, nil
		}
	}

	return "", nil
}
