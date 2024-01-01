package knn

import "math"

func TriangularKernel(distance float64, sigma float64) float64 {
	weight := 1 - distance/sigma
	if weight < 0 {
		weight = 0
	}
	return weight
}
func EpanechnikovKernel(distance float64, sigma float64) float64 {
	weight := 1 - (distance/sigma)*(distance/sigma)
	return weight
}
func GaussianKernel(distance float64, sigma float64) float64 {
	exponent := -(distance * distance) / (2 * sigma * sigma)
	return math.Exp(exponent)
}
