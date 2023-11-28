package main

import (
	"fmt"
	"net/http"

	"github.com/Gurd0/helgelandsbruaBackend/api/knn"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})
	//mount knn routes
	r.Mount("/knn", knn.Routes())
	fmt.Println("started on port 3000")
	http.ListenAndServe(":3000", r)
}
