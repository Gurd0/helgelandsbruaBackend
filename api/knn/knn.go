package knn

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/Gurd0/helgelandsbruaBackend/api/scaler"
)

// DataPoint represents a data point in the dataset.
type DataPoint = scaler.DataPoint

var data []DataPoint
var trainingData []DataPoint
var testData []DataPoint
var GlobalScaler scaler.DataPointTransformer = &scaler.StandardScaler{}

func init() {
	UpdateDataInKNN()
}

var DefaultSetting = Settings{
	K:                 120,
	DistanceMethod:    "Euclidean",
	WeightedDistances: true,
	Sigma:             0.001,
}

type ExampleOptions struct {
	distanceMethod string
}

// distanceSwitch calculates the distance between two DataPoints based on the specified distance method.
// If the specified distance method is not recognized, default to Euclidean distance.
func distanceSwitch(queryPoint, trainingPoint DataPoint, options ExampleOptions) float64 {
	switch os := options.distanceMethod; os {
	case "Euclidean":
		return EuclideanDistance(queryPoint, trainingPoint)
	case "Manhattan":
		return ManhattanDistance(queryPoint, trainingPoint)
	case "Chebyshev":
		return ChebyshevDistance(queryPoint, trainingPoint)
	case "CosineSimilarity":
		return CosineSimilarity(queryPoint, trainingPoint)
	default:
		return EuclideanDistance(queryPoint, trainingPoint)
	}
}

// KNN implements the k-Nearest Neighbors algorithm.
func KNN(queryPoints []scaler.DataPoint, settings Settings) []string {
	if settings.K == 0 && settings.WeightedDistances == false && settings.DistanceMethod == "" && settings.Sigma == 0 {
		settings = DefaultSetting
	}
	predictions := make([]string, len(queryPoints))

	//transforms input data to match training data.
	queryPoints = GlobalScaler.Transform(queryPoints)

	// Parallelizeed distance calculation
	distances := nonParallelDistanceCalculation(queryPoints, trainingData)
	// Process each query point
	for i, queryPointDistances := range distances {
		// Sort distances for the current query point in ascending order.
		sort.Slice(queryPointDistances, func(j, k int) bool {
			return queryPointDistances[j].distance < queryPointDistances[k].distance
		})
		topKDistances := queryPointDistances[:settings.K]
		//TODO, move this to a different file and make a switch like for distance calculation.
		if settings.WeightedDistances {
			//add weighted distances
			labelCounts := make(map[string]float64)
			totalWeight := 0.0
			sigma := settings.Sigma // Adjust this value as needed

			for j := 0; j < settings.K; j++ {
				distance := topKDistances[j].distance     // Rescale the distance
				weight := math.Exp(-(distance / (sigma))) // Use Gaussian function of the rescaled distance as the weight
				totalWeight += weight
			}

			for j := 0; j < settings.K; j++ {
				distance := topKDistances[j].distance
				weight := math.Exp(-(distance / (sigma))) // Use Gaussian function of the rescaled distance as the weight
				normalizedWeight := weight / totalWeight  // Normalize the weight
				label := trainingData[topKDistances[j].trainingIndex].Label
				labelCounts[label] += normalizedWeight
			}
			// Find the label with the maximum weighted count.
			var maxWeight = math.Inf(-1)
			var predictedLabel string
			var labelsWithMaxWeight []string
			for label, weight := range labelCounts {
				if weight > maxWeight {
					maxWeight = weight
					predictedLabel = label
					labelsWithMaxWeight = []string{label}
				} else if weight == maxWeight {
					labelsWithMaxWeight = append(labelsWithMaxWeight, label)
				}
			}

			// Store the predicted label for the current query point.
			predictions[i] = predictedLabel
		} else {
			predictions[i] = trainingData[topKDistances[0].trainingIndex].Label
		}
	}
	return predictions
}
func nonParallelDistanceCalculation(queryPoints []DataPoint, trainingData []DataPoint) [][]struct {
	queryIndex    int
	trainingIndex int
	distance      float64
} {
	// Initialize distances matrix
	distances := make([][]struct {
		queryIndex    int
		trainingIndex int
		distance      float64
	}, len(queryPoints))

	// Sequential calculation of distances
	for i, queryPoint := range queryPoints {
		distances[i] = make([]struct {
			queryIndex    int
			trainingIndex int
			distance      float64
		}, len(trainingData))

		for j, trainingPoint := range trainingData {
			distance := distanceSwitch(queryPoint, trainingPoint, ExampleOptions{""})
			distances[i][j] = struct {
				queryIndex    int
				trainingIndex int
				distance      float64
			}{i, j, distance}
		}
	}

	return distances
}

// TODO, not sure but slower than non parallel
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
			distance := ChebyshevDistance(work.queryPoint, work.trainingPoint)

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
		//features[0] = features[0] / 10

		dataPoint := DataPoint{
			Features: features,
			Label:    label,
		}
		d = append(d, dataPoint)
	}
	//queryPointsForGlobalScaler := scaler.ToDataPointSlice(d)
	GlobalScaler.Fit(d)
	trainingData = GlobalScaler.Transform(d)
}

func Predict(obj PredictInput) float64 {
	var queryPointList []DataPoint
	//scaledX, scaledY := CircularScale(obj.WindDir)
	queryPoint := DataPoint{Features: []float64{obj.Wind, obj.WindDir}, Label: ""}
	queryPointList = append(queryPointList, queryPoint)
	predictedLabel := KNN(queryPointList, Settings{})
	predictedLabelFloat, _ := strconv.ParseFloat(predictedLabel[0], 64)
	res := predictedLabelFloat * obj.Wind
	return res
}
func PredictList(inputList []PredictInput, setting Settings) []string {
	var datapointList []DataPoint
	fmt.Println(setting)
	for _, obj := range inputList {
		queryPoint := DataPoint{Features: []float64{
			obj.Wind,
			obj.WindDir}, Label: ""}
		datapointList = append(datapointList, queryPoint)
	}
	res := KNN(datapointList, setting)
	for index, obj := range res {
		objFloat, _ := strconv.ParseFloat(obj, 64)
		r := objFloat * inputList[index].Wind
		resString := strconv.FormatFloat(r, 'f', 2, 64)
		res[index] = resString
	}

	return res
}

// TODO, did not manage to get good res scaling
func CircularScale(angle float64) (scaledX, scaledY float64) {
	radian := angle * (math.Pi / 180.0)
	scaledX = math.Cos(radian)
	scaledY = math.Sin(radian)
	return scaledX * 10, scaledY * 100
}
