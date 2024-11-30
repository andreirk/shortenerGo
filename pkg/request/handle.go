package request

import (
	"go/adv-demo/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		response.JsonResponse(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = CheckIfValid(body)
	if err != nil {
		response.JsonResponse(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}
