package main

import (
	"net/http"

	"github.com/Gurd0/helgelandsbruaBackend/api/data"
	"github.com/Gurd0/helgelandsbruaBackend/api/knn"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/knn", knn.Routes())
	r.Mount("/data", data.Routes())

	http.ListenAndServe(":3000", r)
}
