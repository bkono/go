package http

import (
	"log"
	"net/http"
)

// log middleware
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log middleware")
		next.ServeHTTP(w, r)
	})
}
