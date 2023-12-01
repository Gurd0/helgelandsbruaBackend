package data

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WindData struct {
	ID             int       `json:"id"`
	Wind           string    `json:"wind"`
	Forcastwind    string    `json:"forcastwind"`
	Forcastwinddir string    `json:"forcastwinddir"`
	Forcastgust    string    `json:"forcastgust"`
	Timestamp      time.Time `json:"timestamp"`
}

func GetJson() {
	requestURL := "https://helgelandsbrua.vercel.app/api/history"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		//TODO return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//TODO return err
	}
	data := []WindData{}
	err = json.Unmarshal([]byte(resBody), &data)
	if err == nil {
		convertJSONToCSV(data, "data.csv")
	}
}
func convertJSONToCSV(source []WindData, destination string) error {

	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for _, r := range source {
		var csvRow []string
		/*windFloat, _ := strconv.ParseFloat(r.Wind, 64)*/

		forcastWindDirFloat, _ := strconv.ParseFloat(r.Forcastwind, 64)
		scaledX, scaledY := CircularScale(forcastWindDirFloat)

		csvRow = append(csvRow, r.Forcastwind, strconv.FormatFloat(scaledX, 'f', -1, 64), strconv.FormatFloat(scaledY, 'f', -1, 64), r.Wind)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

// CircularScale scales a given angle in degrees using sine and cosine transformations.
func CircularScale(angle float64) (scaledX, scaledY float64) {
	radian := angle * (math.Pi / 180.0)
	scaledX = math.Cos(radian)
	scaledY = math.Sin(radian)
	return scaledX, scaledY
}
