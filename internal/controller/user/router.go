package user

import "github.com/go-chi/chi/v5"

func (h *Handler) Register(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
		r.Post("/upsert", h.upsert)
	})
}
