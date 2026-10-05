package routes

import (
	"net/http"
	"synthori/ediary/m/handlers"
	"synthori/ediary/m/middleware"
	"synthori/ediary/m/models"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RateLimitMiddleware)
	r.Get("/", healthCheck)

	r.Route("/user", func(r chi.Router) {
		r.Get("/dummy", h.User.GetDummyUser)
		r.Post("/", h.User.CreateUser)
		r.Get("/{id}", h.User.GetUserByID)
		r.Get("/", h.User.GetUsers)
	})

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	handlers.WriteSuccess(w, 200, "hellow world", models.StatusCheck{
		Message:    "ok",
		StatusCode: 200,
	})
}
