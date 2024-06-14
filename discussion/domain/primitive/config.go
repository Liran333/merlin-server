package primitive

var cfg *Config

type Config struct {
	MaxTitleLength   int `json:"max_title_length"`
	MaxContentLength int `json:"max_content_length"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxTitleLength <= 0 {
		cfg.MaxTitleLength = 200
	}

	if cfg.MaxContentLength <= 0 {
		cfg.MaxContentLength = 10000
	}
}

func InitConfig(c *Config) {
	cfg = c
}
