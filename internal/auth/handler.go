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
	*AuthService
}
type AuthHandler struct {
	*config.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/login", handler.Login())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		log.Println("Register with payload", payload)
		handler.AuthService.Register(payload.Email, payload.Name, payload.Password)
		res := RegisterResponse{
			RegisterSuccess: true,
		}
		response.JsonResponse(w, res, http.StatusOK)
		fmt.Println("Register")
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		fmt.Println("Login with payload", payload)
		res := LoginResponse{
			AccessToken: "token 1234",
		}
		response.JsonResponse(w, res, http.StatusOK)
	}
}
