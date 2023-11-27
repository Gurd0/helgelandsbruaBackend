package handler

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	fmt.Println(p)
	fmt.Fprint(w, "Welcome to the home page!")
}
