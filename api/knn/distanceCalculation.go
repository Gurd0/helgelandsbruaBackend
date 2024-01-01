package knn

import "math"

//  length of the line segment between two points
func EuclideanDistance(a, b DataPoint) float64 {
	sum := 0.0
	for i := 0; i < len(a.Features); i++ {
		diff := a.Features[i] - b.Features[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

//the sum of the absolute differences between the coordinates.
func ManhattanDistance(a, b DataPoint) float64 {
	sum := 0.0
	for i := 0; i < len(a.Features); i++ {
		sum += math.Abs(a.Features[i] - b.Features[i])
	}
	return sum
}

//The maximum absolute difference between coordinates.
func ChebyshevDistance(a, b DataPoint) float64 {
	maxDiff := 0.0
	for i := 0; i < len(a.Features); i++ {
		diff := math.Abs(a.Features[i] - b.Features[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}

//Measures the cosine of the angle between two vectors.
func CosineSimilarity(a, b DataPoint) float64 {
	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0
	for i := 0; i < len(a.Features); i++ {
		dotProduct += a.Features[i] * b.Features[i]
		magnitudeA += math.Pow(a.Features[i], 2)
		magnitudeB += math.Pow(b.Features[i], 2)
	}
	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}
