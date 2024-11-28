package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *API) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/chats", func(r chi.Router) {

			r.Get("/", a.GetChats)
			r.Post("/", a.CreateChat)

		})

		r.Route("/users", func(r chi.Router) {

			r.Post("/", a.CreateUser)

			r.Post("/login", a.LoginUser)

		})

	})

	r.Get("/ws", a.WebSocket)

	return r
}
