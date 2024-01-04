package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxSessionNum           int   `json:"max_session_num"`
	CSRFTokenTimeout        int64 `json:"csrf_token_timeout"`
	CSRFTokenTimeoutToReset int64 `json:"csrf_token_timeout_to_reset"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxSessionNum <= 0 {
		cfg.MaxSessionNum = 3
	}

	if cfg.CSRFTokenTimeout <= 0 {
		cfg.CSRFTokenTimeout = 8 * 3600
	}

	if cfg.CSRFTokenTimeoutToReset <= 0 {
		cfg.CSRFTokenTimeoutToReset = 3
	}
}
