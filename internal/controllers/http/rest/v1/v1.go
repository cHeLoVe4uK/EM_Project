package v1

import "github.com/julienschmidt/httprouter"

type Handler struct {
	// TODO service, logger
}

type Deps struct {
}

func NewHandler(d *Deps) *Handler {
	return &Handler{
		// TODO
	}
}

func (h *Handler) Init(r *httprouter.Router) {
	h.initUserHandler(r)
}
