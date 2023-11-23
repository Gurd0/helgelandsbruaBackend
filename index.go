package main

import (
	"net/http"

	"github.com/Gurd0/helgelandsbruaBackend/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/knn", func(r chi.Router) {
			r.Get("/", api.PostPredict)
		})

		r.Route("/data", func(r chi.Router) {
			r.Get("/update", api.GetUpdateData)
		})

	})

	http.ListenAndServe(":3000", r)
}
