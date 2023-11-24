package data

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/data/update", GetUpdateData)
	return r
}

func GetUpdateData(w http.ResponseWriter, r *http.Request) {
	GetJson()
}
