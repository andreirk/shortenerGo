package midleware

import (
	"context"
	"go/adv-demo/config"
	"go/adv-demo/pkg/jwt"
	"log"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func CheckAuthed(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.Split(authHeader, " ")[1] // Remove "Bearer:"
		data, isValid := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		//log.Println(data, isValid)
		r.Context()
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		log.Println(authHeader, "Token:", token)
		next.ServeHTTP(w, req)
	})
}
