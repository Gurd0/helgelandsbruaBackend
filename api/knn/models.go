package knn

type PredictInput struct {
	Wind     float64 `json:"wind"`
	WindDir  float64 `json:"windDir"`
	WindGust float64 `json:"windGust"`
}
type PredictInputList struct {
	Wind     []PredictInput `json:"wind"`
	Gust     []PredictInput `json:"gust"`
	Settings Settings       `json:"settings"`
}
type PredictRespons struct {
	Wind float64 `json:"wind"`
}
type PredictResponsList struct {
	Wind []string `json:"wind"`
	Gust []string `json:"gust"`
}
type Settings struct {
	K                 int    `json:"k"`
	WeightedDistances bool   `json:"weightedDistances"`
	DistanceMethod    string `json:"distanceMethod"`
	//if weightedDistances
	Sigma           float64 `json:"sigma"`
	WeigthingMethod string  `json:"weigthingMethod"`
}
