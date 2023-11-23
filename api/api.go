package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	data "github.com/Gurd0/helgelandsbruaBackend/api/dataSet"
	"github.com/Gurd0/helgelandsbruaBackend/api/knn"
)

type predictRespons struct {
	Predict string `json:"predict"`
}

func GetUpdateData(w http.ResponseWriter, r *http.Request) {
	data.GetJson()
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
		Predict: knn.Predict(forcastWind, forcastDir, forcastGust),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	/*knn.Predict(forcastWind, forcastDir, forcastGust)
	d1 := Message{knn.Predict(1, 1, 2)}
	fmt.Println(d1)
	render.JSON(w, r, d1)
	return*/
}
