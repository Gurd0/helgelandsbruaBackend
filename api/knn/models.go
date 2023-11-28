package knn

type PredictInput struct {
	Wind     float64 `json:"wind"`
	WindDir  float64 `json:"windDir"`
	WindGust float64 `json:"windGust"`
}
type PredictInputList struct {
	Wind []PredictInput `json:"wind"`
	Gust []PredictInput `json:"gust"`
}
type PredictRespons struct {
	Wind float64 `json:"wind"`
}
type PredictResponsList struct {
	Wind []string `json:"wind"`
	Gust []string `json:"gust"`
}
