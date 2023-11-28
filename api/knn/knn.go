package knn

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

// DataPoint represents a data point in the dataset.
type DataPoint struct {
	Features []float64
	Label    string
}

var trainingData []DataPoint
var k = 10

func init() {
	UpdateDataInKNN()
}

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
	//use min max on input
	queryPoint.Features = minMaxScaling(queryPoint.Features, 0, 365)
	//combine the points to reduce dimensions
	queryPoint.Features = []float64{queryPoint.Features[0] + queryPoint.Features[1]}
	fmt.Println(queryPoint.Features)
	// Parallelizeed distance calculation
	distances := parallelDistanceCalculation(queryPoint, trainingData, 4)

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
func parallelDistanceCalculation(queryPoint DataPoint, trainingData []DataPoint, numWorkers int) []struct {
	index    int
	distance float64
} {
	var wg sync.WaitGroup
	distancesMutex := sync.Mutex{}
	distances := make([]struct {
		index    int
		distance float64
	}, len(trainingData))

	// Channel for sending work to workers
	workCh := make(chan int, len(trainingData))

	// Channel for receiving results from workers
	resultCh := make(chan struct {
		index    int
		distance float64
	}, len(trainingData))

	// Function to launch a worker
	worker := func() {
		defer wg.Done()
		for i := range workCh {
			dataPoint := trainingData[i]
			distance := Distance(queryPoint, dataPoint)

			resultCh <- struct {
				index    int
				distance float64
			}{i, distance}
		}
	}

	// Start worker goroutines
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker()
	}

	// Send work to workers
	go func() {
		defer close(workCh)
		for i := 0; i < len(trainingData); i++ {
			workCh <- i
		}
	}()

	// Close resultCh when all workers are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results from workers
	go func() {
		for result := range resultCh {
			distancesMutex.Lock()
			distances[result.index] = result
			distancesMutex.Unlock()
		}
	}()

	// Wait for all workers to finish
	wg.Wait()

	return distances
}
func UpdateDataInKNN() {
	fmt.Println("update run")
	//TODO, må no fins ein måte å få læst filæ
	resp, err := http.Get("https://helgelandsbrua-backend.vercel.app/data.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.Status)
		return
	}

	reader := csv.NewReader(resp.Body)
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

		features = minMaxScaling(features, 0, 365)
		features = []float64{features[0] + features[1]}

		dataPoint := DataPoint{
			Features: features,
			Label:    label,
		}
		trainingData = append(trainingData, dataPoint)
	}
}

// min max scaling for points to reduce size diff.
func minMaxScaling(data []float64, min, max float64) []float64 {
	scaledData := make([]float64, len(data))

	for i, value := range data {
		scaledData[i] = (value - min) / (max - min)
	}
	//fmt.Println(scaledData)
	return scaledData
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
