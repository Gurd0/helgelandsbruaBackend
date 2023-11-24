package main

import (
	"fmt"
	"net/http"

	"github.com/Gurd0/helgelandsbruaBackend/data"
	"github.com/Gurd0/helgelandsbruaBackend/knn"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/knn", knn.Routes())
	r.Mount("/data", data.Routes())

	fmt.Printf("running on port :3000")
	http.ListenAndServe(":3000", r)

}
