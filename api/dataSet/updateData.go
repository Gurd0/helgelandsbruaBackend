package data

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
		csvRow = append(csvRow, r.Forcastwind, r.Forcastwinddir, r.Forcastgust, r.Wind)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
