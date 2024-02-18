package app

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxCountPerOwner int `json:"max_count_per_owner"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 1000
	}
}
