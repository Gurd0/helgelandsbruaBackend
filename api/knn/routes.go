package knn

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type predictRespons struct {
	Predict string `json:"predict"`
}

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/knn", PostPredict)
	return r
}

func PostPredict(w http.ResponseWriter, r *http.Request) {
	forcastWind, err := strconv.ParseFloat(r.URL.Query().Get("forcastWind"), 64)
	if err != nil {
		//TODO error handle

	}
	forcastDir, err := strconv.ParseFloat(r.URL.Query().Get("forcastDir"), 64)
	if err != nil {
		//TODO error handle

	}
	forcastGust, err := strconv.ParseFloat(r.URL.Query().Get("forcastGust"), 64)
	if err != nil {
		//TODO error handle

	}
	data := predictRespons{
		Predict: Predict(forcastWind, forcastDir, forcastGust),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
