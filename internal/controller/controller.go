package controller

import (
	"log/slog"

	"example/internal/controller/user"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Controller struct {
	user *user.Handler
	log  *slog.Logger
}

func New(user *user.Handler, log *slog.Logger) *Controller {
	return &Controller{
		user: user,
		log:  log,
	}
}

func (c *Controller) Setup() *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		c.user.Register(r)
	})

	return r
}
