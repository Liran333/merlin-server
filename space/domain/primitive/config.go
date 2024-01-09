package primitive

import "strings"

var (
	allSDK      map[string]bool
	allHardware map[string]bool
)

type Config struct {
	SDK      []string `json:"sdk"      required:"true"`
	Hardware []string `json:"hardware" required:"true"`
}

func Init(cfg *Config) {
	allHardware = map[string]bool{}
	for _, v := range cfg.Hardware {
		allHardware[strings.ToLower(v)] = true
	}

	allSDK = map[string]bool{}
	for _, sv := range cfg.SDK {
		allSDK[strings.ToLower(sv)] = true
	}
}
