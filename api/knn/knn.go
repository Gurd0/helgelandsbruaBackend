package knn

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

// DataPoint represents a data point in the dataset.
type DataPoint struct {
	Features []float64
	Label    string
}

var trainingData []DataPoint
var k = 5

// Distance calculates the Euclidean distance between two data points.
func Distance(a, b DataPoint) float64 {
	sum := 0.0
	for i := 0; i < len(a.Features); i++ {
		diff := a.Features[i] - b.Features[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// KNN implements the k-Nearest Neighbors algorithm.
func KNN(queryPoint DataPoint) string {
	// Calculate distances from the query point to all training points.
	distances := make([]struct {
		index    int
		distance float64
	}, len(trainingData))

	for i, dataPoint := range trainingData {
		distances[i] = struct {
			index    int
			distance float64
		}{i, Distance(queryPoint, dataPoint)}
	}

	// Sort distances in ascending order.
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Count the occurrences of each label among the k nearest neighbors.
	labelCounts := make(map[string]int)
	for i := 0; i < k; i++ {
		label := trainingData[distances[i].index].Label
		labelCounts[label]++
	}

	// Find the label with the maximum count.
	maxCount := 0
	var predictedLabel string
	for label, count := range labelCounts {
		if count > maxCount {
			maxCount = count
			predictedLabel = label
		}
	}

	return predictedLabel
}

func init() {
	UpdateDataInKNN()
}
func UpdateDataInKNN() {
	// Open the CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Extract data from CSV records

	for _, record := range records {
		// Assuming the last column is the label, and the rest are features
		features := make([]float64, len(record)-1)
		for i := 0; i < len(record)-1; i++ {
			features[i], err = strconv.ParseFloat(record[i], 64)
			if err != nil {
				log.Fatal(err)
			}
		}

		label := record[len(record)-1]
		//test to see if reducing size will help speed
		features[0] = features[0] / 10
		features[1] = features[1] / 100
		features[2] = features[2] / 10
		dataPoint := DataPoint{
			Features: features,
			Label:    label,
		}
		trainingData = append(trainingData, dataPoint)
	}
}
func Predict(obj PredictInput) float64 {
	queryPoint := DataPoint{Features: []float64{obj.Wind, obj.WindDir}, Label: ""}
	predictedLabel := KNN(queryPoint)
	predictedLabelFloat, _ := strconv.ParseFloat(predictedLabel, 64)
	res := predictedLabelFloat * obj.Wind
	return res
}
func PredictList(inputList []PredictInput) []string {
	var predictedLabels []string
	for _, obj := range inputList {
		queryPoint := DataPoint{Features: []float64{
			math.Round((obj.Wind/10)*100) / 100,
			obj.WindDir / 100,
			math.Round((obj.WindGust/10)*100) / 100}, Label: ""}
		predictedLabel := KNN(queryPoint)

		predictedLabelFloat, _ := strconv.ParseFloat(predictedLabel, 64)
		res := predictedLabelFloat * obj.Wind
		res = math.Round(res*100) / 100
		predictedLabels = append(predictedLabels, strconv.FormatFloat(res, 'f', -1, 64))
	}
	return predictedLabels
}
