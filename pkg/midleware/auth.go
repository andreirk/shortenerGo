package midleware

import (
	"log"
	"net/http"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.Split(authHeader, " ")[1] // Remove "Bearer:"
		log.Println(authHeader, "Token:", token)
		if token != "" {
			next.ServeHTTP(w, r)
		}

	})
}
