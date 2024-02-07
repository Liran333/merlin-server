package domain

type tables struct {
	User  string `json:"user" required:"true"`
	Token string `json:"token" required:"true"`
}

type Config struct {
	Tables tables `json:"tables"            required:"true"`
	Key    []byte `json:"key"               required:"true"`
}

var _config Config

func Init(cfg *Config) {
	_config = *cfg
}
