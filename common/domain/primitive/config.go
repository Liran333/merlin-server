package primitive

import "strings"

var (
	msdConfig   MSDConfig
	allLicenses map[string]bool
)

func Init(cfg *Config) {
	msdConfig = cfg.MSDConfig

	m := map[string]bool{}
	for _, v := range cfg.Licenses {
		m[strings.ToLower(v)] = true
	}

	allLicenses = m
}

// Config
type Config struct {
	MSDConfig

	Licenses []string `json:"licenses" required:"true"`
}

// MSDConfig
type MSDConfig struct {
	MaxNameLength     int `json:"max_name_length"`
	MinNameLength     int `json:"min_name_length"`
	MaxDescLength     int `json:"max_desc_length"`
	MaxFullnameLength int `json:"max_fullname_length"`
}

func (cfg *MSDConfig) SetDefault() {
	if cfg.MaxNameLength <= 0 {
		cfg.MaxNameLength = 50
	}

	if cfg.MinNameLength <= 0 {
		cfg.MinNameLength = 5
	}

	if cfg.MaxDescLength <= 0 {
		cfg.MaxDescLength = 200
	}

	if cfg.MaxFullnameLength <= 0 {
		cfg.MaxFullnameLength = 200
	}
}
