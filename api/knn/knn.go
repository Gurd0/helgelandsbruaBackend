package knn

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
)

// DataPoint represents a data point in the dataset.
type DataPoint struct {
	Features []float64
	Label    string
}

var data []DataPoint
var trainingData []DataPoint
var testData []DataPoint
var scaler = StandardScaler{}
var k = 40

func init() {
	UpdateDataInKNN()
}

type StandardScaler struct {
	Mean    []float64
	StdDev  []float64
	Trained bool
}

// Fit computes the mean and standard deviation of each feature from the provided dataset.
func (scaler *StandardScaler) Fit(data []DataPoint) {
	numFeatures := len(data[0].Features)
	numSamples := len(data)

	// Initialize Mean and StdDev slices
	scaler.Mean = make([]float64, numFeatures)
	scaler.StdDev = make([]float64, numFeatures)

	// Calculate mean
	for _, point := range data {
		for j, value := range point.Features {
			scaler.Mean[j] += value
		}
	}

	for j := range scaler.Mean {
		scaler.Mean[j] /= float64(numSamples)
	}

	// Calculate standard deviation
	for _, point := range data {
		for j, value := range point.Features {
			scaler.StdDev[j] += math.Pow(value-scaler.Mean[j], 2)
		}
	}

	for j := range scaler.StdDev {
		scaler.StdDev[j] = math.Sqrt(scaler.StdDev[j] / float64(numSamples))
	}

	scaler.Trained = true
}

// Transform scales the input data based on the mean and standard deviation computed during fitting.
func (scaler *StandardScaler) Transform(data []DataPoint) []DataPoint {
	if !scaler.Trained {
		fmt.Println("Scaler has not been trained. Call Fit() first.")
		return nil
	}
	numSamples := len(data)
	numFeatures := len(data[0].Features)

	transformed := make([]DataPoint, numSamples)

	for i, point := range data {
		transformed[i] = DataPoint{
			Features: make([]float64, numFeatures),
			Label:    point.Label,
		}
		for j, value := range point.Features {
			transformed[i].Features[j] = (value - scaler.Mean[j]) / scaler.StdDev[j]
		}

		//transformed[i].Features[0] = transformed[i].Features[0] / 2
	}
	return transformed
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
func KNN(queryPoints []DataPoint, scaleData bool) []string {
	predictions := make([]string, len(queryPoints))
	if scaleData {
		//use min max on input
		//queryPoints = minMaxScaling(queryPoints, 0, 365)
		queryPoints = scaler.Transform(queryPoints)

	}
	// Parallelizeed distance calculation
	distances := parallelDistanceCalculation(queryPoints, trainingData, 4)
	// Process each query point
	for i, queryPointDistances := range distances {
		// Sort distances for the current query point in ascending order.
		sort.Slice(queryPointDistances, func(j, k int) bool {
			return queryPointDistances[i].distance < queryPointDistances[j].distance
		})

		// Count the occurrences of each label among the k nearest neighbors.
		labelCounts := make(map[string]int)
		for j := 0; j < k; j++ {
			label := trainingData[queryPointDistances[j].trainingIndex].Label
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

		// Store the predicted label for the current query point.
		predictions[i] = predictedLabel
	}

	return predictions
}
func parallelDistanceCalculation(queryPoints []DataPoint, trainingData []DataPoint, numWorkers int) [][]struct {
	queryIndex    int
	trainingIndex int
	distance      float64
} {
	var wg sync.WaitGroup
	distancesMutex := sync.Mutex{}

	// Create a two-dimensional slice for distances
	distances := make([][]struct {
		queryIndex    int
		trainingIndex int
		distance      float64
	}, len(queryPoints))

	// Channel for sending work to workers
	workCh := make(chan struct {
		queryIndex    int
		queryPoint    DataPoint
		trainingIndex int
		trainingPoint DataPoint
	}, len(queryPoints)*len(trainingData))

	// Channel for receiving results from workers
	resultCh := make(chan struct {
		queryIndex    int
		trainingIndex int
		distance      float64
	}, len(queryPoints)*len(trainingData))

	// Function to launch a worker
	worker := func() {
		defer wg.Done()
		for work := range workCh {
			distance := Distance(work.queryPoint, work.trainingPoint)

			resultCh <- struct {
				queryIndex    int
				trainingIndex int
				distance      float64
			}{work.queryIndex, work.trainingIndex, distance}
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
		for i, queryPoint := range queryPoints {
			// Initialize the inner slice
			distances[i] = make([]struct {
				queryIndex    int
				trainingIndex int
				distance      float64
			}, len(trainingData))

			for j, trainingPoint := range trainingData {
				workCh <- struct {
					queryIndex    int
					queryPoint    DataPoint
					trainingIndex int
					trainingPoint DataPoint
				}{i, queryPoint, j, trainingPoint}
			}
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
			distances[result.queryIndex][result.trainingIndex] = struct {
				queryIndex    int
				trainingIndex int
				distance      float64
			}{result.queryIndex, result.trainingIndex, result.distance}
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
	/*resp, err := http.Get("https://helgelandsbrua-backend.vercel.app/data.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code", resp.Status)
		return
	} */
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Extract data from CSV records
	var d []DataPoint
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
		//features[1] = features[1] / 100

		//features = minMaxScaling(features, 0, 360)
		//features = []float64{features[0] + features[1]}

		dataPoint := DataPoint{
			Features: features,
			Label:    label,
		}
		d = append(d, dataPoint)
	}
	scaler.Fit(d)
	SplitData(scaler.Transform(d), 0.05)
}

// SplitData randomly splits the data into training and test sets.
func SplitData(data []DataPoint, testPercentage float64) {
	rand.Seed(42)

	// Shuffle the data.
	shuffledData := make([]DataPoint, len(data))
	copy(shuffledData, data)
	rand.Shuffle(len(shuffledData), func(i, j int) {
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	})

	// Determine the split index based on the test percentage.
	splitIndex := int(float64(len(shuffledData)) * testPercentage)

	// Split the data.
	trainingData = shuffledData[splitIndex:]
	testData = shuffledData[:splitIndex]
	fmt.Println(testData)
}

func Predict(obj PredictInput) float64 {
	var queryPointList []DataPoint
	queryPoint := DataPoint{Features: []float64{obj.Wind, obj.WindDir}, Label: ""}
	queryPointList = append(queryPointList, queryPoint)
	predictedLabel := KNN(queryPointList, true)
	predictedLabelFloat, _ := strconv.ParseFloat(predictedLabel[0], 64)
	res := predictedLabelFloat
	return res
}
func PredictList(inputList []PredictInput) []string {
	var datapointList []DataPoint
	for _, obj := range inputList {
		queryPoint := DataPoint{Features: []float64{
			obj.Wind,
			obj.WindDir}, Label: ""}
		datapointList = append(datapointList, queryPoint)
	}
	res := KNN(datapointList, true)
	for index, obj := range res {
		objFloat, _ := strconv.ParseFloat(obj, 64)
		r := objFloat
		resString := strconv.FormatFloat(r, 'f', 2, 64)
		res[index] = resString
	}

	return res
}
