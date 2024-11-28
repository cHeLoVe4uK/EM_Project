package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *API) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/chats", a.GetChats)
	r.Post("/api/v1/chats", a.CreateChat)

	r.Get("/ws", a.WebSocket)

	return r
}
