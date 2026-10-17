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
	// r.Use(middleware.RateLimitMiddleware)
	r.Get("/health", healthCheck)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.User.AuthUser)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/dummy", h.User.GetDummyUser)
		r.Post("/", h.User.CreateUser)
		r.Get("/{id}", h.User.GetUserByID)
		r.Get("/", h.User.GetUsers)
	})

	r.Route("/diary", func(r chi.Router) {
		r.Get("/{id}", h.Diary.GetDayByID)
		r.Get("/week", h.Diary.GetAllWeek)
	})

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	setCookieHandler(w, r)
	handlers.WriteSuccess(w, 200, "hellow world", models.StatusCheck{
		Message:    "ok",
		StatusCode: 200,
	})

}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "user_id",
		Value:    "1",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
