package data

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/update", GetUpdateData)
	return r
}

func GetUpdateData(w http.ResponseWriter, r *http.Request) {
	GetJson()
	w.Header().Set("Content-Type", "application/json")

}
