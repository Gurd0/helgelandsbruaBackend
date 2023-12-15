package knn

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

var content embed.FS

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getPredict)
	r.Post("/", postPredict)
	return r
}
func getPredict(w http.ResponseWriter, r *http.Request) {
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
	predictObj := PredictInput{Wind: forcastWind, WindDir: forcastDir, WindGust: forcastGust}
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
func postPredict(w http.ResponseWriter, r *http.Request) {
	var windList PredictInputList
	err := json.NewDecoder(r.Body).Decode(&windList)
	if err != nil {
		fmt.Println(err)
	}
	predictionWind := PredictList(windList.Wind, windList.Settings)
	predictionGust := predictionWind

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
