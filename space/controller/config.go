package controller

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxCountPerPage int `json:"max_count_per_page"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerPage <= 0 {
		cfg.MaxCountPerPage = 100
	}
}
