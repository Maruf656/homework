package models

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type HistoryElement struct {
	NumberOne float64 `json:"num_one"`
	NumberTwo float64 `json:"num_two"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}
