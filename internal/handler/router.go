package handler

import (
	"net/http"

	"github.com/Oxygenss/yandex_go_final_project/internal/handler/middleware"
	"github.com/go-chi/chi"
)

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/signin", h.SignIn)

	router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.AuthMiddleware(h.cfg.Auth.Secret, next)
		})

		r.Post("/api/task", h.AddTask)
		r.Get("/api/task", h.GetTaskByID)
		r.Put("/api/task", h.EditTask)
		r.Delete("/api/task", h.DeleteTask)
		r.Post("/api/task/done", h.DoneTask)
		r.Get("/api/tasks", h.GetTasks)
	})

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	router.Handle("/*", http.StripPrefix("/", fs))

	return router
}
