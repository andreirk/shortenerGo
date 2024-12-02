package link

import (
	"fmt"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"log"
	"net/http"
)

type LinkHandler struct {
}

type LinkDeps struct {
}

func NewLinkHandler(router *http.ServeMux, deps LinkDeps) *LinkHandler {
	handler := &LinkHandler{}
	router.HandleFunc("POST /link", handler.Create())
	router.HandleFunc("PATCH /link/{id}", handler.Update())
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{alias}", handler.GoTo())

	return handler
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[CreateLinkRequest](&w, r)
		if err != nil {
			return
		}
		log.Println("Create link request", payload)
		res := CreateLinkResponse{
			Success: true,
		}
		response.JsonResponse(w, res, http.StatusOK)

	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println("Delete link request", id)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
