package router

import (
	"net/http"
	"practice4/practice-4/internal/handler"
	"practice4/practice-4/internal/middleware"
)

func NewRouter(h *handler.UserHandler, authkey string) http.Handler {
	authedMux := http.NewServeMux()
	authedMux.HandleFunc("GET /users", h.GetAll)
	authedMux.HandleFunc("GET /users/{id}", h.GetByID)
	authedMux.HandleFunc("POST /users", h.Create)
	authedMux.HandleFunc("PUT /users/{id}", h.Update)
	authedMux.HandleFunc("DELETE /users/{id}", h.Delete)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.Handle("/", middleware.AuthMiddleware(authkey)(authedMux))

	return middleware.LoggingMiddleware(mux)
}
