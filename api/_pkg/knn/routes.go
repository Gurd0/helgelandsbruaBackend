package knn

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/knn", GetPredict)
	r.Post("/knn", PostPredict)
	return r
}

func GetPredict(w http.ResponseWriter, r *http.Request) {
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
	predictObj := PredictInput{forcastWind, forcastDir, forcastGust}
	data := PredictRespons{
		Wind: Predict(predictObj),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func PostPredict(w http.ResponseWriter, r *http.Request) {

	var windList PredictInputList
	err := json.NewDecoder(r.Body).Decode(&windList)
	if err != nil {
		//TODO error
	}
	predictionWind := PredictList(windList.Wind)
	predictionGust := PredictList(windList.Gust)
	data := PredictResponsList{
		Wind: predictionWind,
		Gust: predictionGust,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		//TODO error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
