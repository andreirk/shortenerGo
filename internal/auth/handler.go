package auth

import (
	"fmt"
	"go/adv-demo/config"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"log"
	"net/http"
)

type AuthHandlerDeps struct {
	*config.Config
}
type AuthHandler struct {
	*config.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/login", handler.Login())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(handler.Config.Auth.Secret)
		fmt.Println("Register")
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			log.Println(err)
			response.JsonResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Login with payload", payload)
		res := LoginResponse{
			AccessToken: "token 1234",
		}
		response.JsonResponse(w, res, http.StatusOK)
	}
}
