package handler

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gurd0/helgelandsbruaBackend/api/_pkg/knn"
)

var content embed.FS

func Knn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Access-Control-Allow-Methods", "GET,OPTIONS,PATCH,DELETE,POST,PUT")
	knn.UpdateDataInKNN()
	switch r.Method {
	case "GET":
		fmt.Println("Get")
		getPredict(w, r)
	case "POST":
		postPredict(w, r)
	default:
		fmt.Println("It's something else.")
	}
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
func postPredict(w http.ResponseWriter, r *http.Request) {

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
