package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Gurd0/helgelandsbruaBackend/api/_pkg/data"
	"github.com/Gurd0/helgelandsbruaBackend/api/_pkg/knn"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func Routes() *chi.Mux {
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
	predictObj := knn.PredictInput{Wind: forcastWind, WindDir: forcastDir, WindGust: forcastGust}
	data := knn.PredictRespons{
		Wind: knn.Predict(predictObj),
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

	var windList knn.PredictInputList
	err := json.NewDecoder(r.Body).Decode(&windList)
	if err != nil {
		//TODO error
	}
	predictionWind := knn.PredictList(windList.Wind)
	predictionGust := knn.PredictList(windList.Gust)
	data := knn.PredictResponsList{
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
