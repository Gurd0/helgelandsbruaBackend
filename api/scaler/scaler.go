package scaler

import (
	"fmt"
	"math"
)

type StandardScaler struct {
	Mean    []float64
	StdDev  []float64
	Trained bool
}

// DataPoint represents a data point in the dataset.
type DataPoint struct {
	Features []float64
	Label    string
}
type DataPointTransformer interface {
	Transform(data []DataPoint) []DataPoint
	Fit(data []DataPoint)
}

// Fit computes the mean and standard deviation of each feature from the provided dataset.
func (scaler *StandardScaler) Fit(data []DataPoint) {
	numFeatures := 1
	numSamples := len(data)

	// Initialize Mean and StdDev slices
	scaler.Mean = make([]float64, numFeatures)
	scaler.StdDev = make([]float64, numFeatures)

	// Calculate mean
	for _, point := range data {
		scaler.Mean[0] += point.Features[0]
		/*for j, value := range point.Features {
			scaler.Mean[j] += value
		} */
	}
	scaler.Mean[0] /= float64(numSamples)
	/*for j := range scaler.Mean {
		scaler.Mean[j] /= float64(numSamples)
	} */

	// Calculate standard deviation
	for _, point := range data {
		scaler.StdDev[0] += math.Pow(point.Features[0]-scaler.Mean[0], 2)
	}
	/*for _, point := range data {
		for j, value := range point.Features {
			scaler.StdDev[j] += math.Pow(value-scaler.Mean[j], 2)
		}
	} */

	/*for j := range scaler.StdDev {
		scaler.StdDev[j] = math.Sqrt(scaler.StdDev[j] / float64(numSamples))
	} */

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
		/*for j, value := range point.Features {
			transformed[i].Features[j] =
		} */
		transformed[i].Features[0] = point.Features[0] //((point.Features[0] - scaler.Mean[0]) / scaler.StdDev[0])
		transformed[i].Features[1] = point.Features[1]
		transformed[i].Features[2] = point.Features[2]
		//transformed[i].Features[0] = transformed[i].Features[0] / 2
	}
	return transformed
}

// ToDataPointSlice converts a slice of DataPoint to []scaler.DataPoint
func ToDataPointSlice(data []DataPoint) []DataPoint {
	result := make([]DataPoint, len(data))
	for i, dp := range data {
		result[i] = dp
	}
	return result
}
