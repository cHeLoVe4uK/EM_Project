package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *API) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(corsMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler())

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/chats", func(r chi.Router) {

			r.Get("/", a.GetAllChats)
			r.Post("/", a.CreateChat)

			r.Get("/{id}/connect", a.ConnectChat)

		})

		r.Route("/users", func(r chi.Router) {

			r.Post("/", a.CreateUser)

			r.Post("/login", a.LoginUser)

		})

	})

	return r
}
