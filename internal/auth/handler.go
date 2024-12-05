package auth

import (
	"fmt"
	"go/adv-demo/config"
	"go/adv-demo/pkg/jwt"
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
			Message:         "Please login",
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
		email, err := handler.AuthService.Login(payload.Email, payload.Password)
		if err != nil {
			response.JsonResponse(w, nil, http.StatusUnauthorized)
			return
		}
		jwtToken, err := jwt.NewJwt(handler.Config.Auth.Secret).Sign(jwt.JwtData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		resp := LoginResponse{
			AccessToken: jwtToken,
		}
		response.JsonResponse(w, resp, http.StatusOK)
	}
}
