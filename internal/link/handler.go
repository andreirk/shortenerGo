package link

import (
	"fmt"
	"go/adv-demo/pkg/midleware"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
}

type LinkDeps struct {
	LinkRepository *LinkRepository
}

func NewLinkHandler(router *http.ServeMux, deps LinkDeps) *LinkHandler {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", midleware.IsAuthed(handler.Update()))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /link/{hash}", handler.GoTo())

	return handler
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := request.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return
		}
		log.Println("Create link request", payload)
		// Business logic start
		link := NewLink(payload.Url, 0)
		for range 5 {
			existingLink, _ := handler.LinkRepository.GetByHash(link.Hash)
			if existingLink == nil {
				break
			}
			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// BL end

		response.JsonResponse(w, createdLink, http.StatusOK)

	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{
				ID: uint(id),
			},
			Url:  body.Url,
			Hash: body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.JsonResponse(w, link, http.StatusOK)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.FindById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.JsonResponse(w, nil, http.StatusOK)

		fmt.Println("Delete link request", id)

	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		foundLink, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, foundLink.Url, http.StatusTemporaryRedirect)
	}
}
