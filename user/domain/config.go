package domain

type Config struct {
	AvatarURL    []string `json:"avatar_url"        required:"true"`
	MaxBioLength int      `json:"max_bio_length"    required:"true"`
}

var _config Config

func Init(cfg *Config) {
	_config = *cfg
}

func MaxBioLength() int {
	return _config.MaxBioLength
}

func DefaultAvatarURL() string {
	return _config.AvatarURL[0]
}
