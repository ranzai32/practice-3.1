package middleware

import (
    "log"
    "net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Method: %s, URL: %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
    })
}

func AuthMiddleware(apiKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.Header.Get("X-API-KEY")
            if key != apiKey {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte(`{"error": "Unauthorized"}`))
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}