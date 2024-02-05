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
	if cfg == nil {
		return
	}

	allHardware = map[string]bool{}
	if cfg.Hardware != nil {
		for _, v := range cfg.Hardware {
			allHardware[strings.ToLower(v)] = true
		}
	}

	allSDK = map[string]bool{}
	if cfg.SDK != nil {
		for _, sv := range cfg.SDK {
			allSDK[strings.ToLower(sv)] = true
		}
	}
}
