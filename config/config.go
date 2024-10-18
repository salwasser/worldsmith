package config

type Config struct {
	Window struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window"`
}
