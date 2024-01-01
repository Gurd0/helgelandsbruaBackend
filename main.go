package main

import (
	"fmt"
	"net/http"

	"github.com/Gurd0/helgelandsbruaBackend/api/data"
	"github.com/Gurd0/helgelandsbruaBackend/api/knn"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
		MaxAge:           300, // Maximum age for caching
	})
	r.Use(corsMiddleware.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})
	//mount knn routes
	r.Mount("/knn", knn.Routes())
	r.Mount("/data", data.Routes())
	fmt.Println("started on port 3000")
	http.ListenAndServe(":8080", r)

}
