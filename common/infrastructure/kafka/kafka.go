package kafka

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/kafka-lib/mq"
)

const (
	deaultVersion = "2.1.0"
)

var Exit = kfklib.Exit

// Config
type Config struct {
	kfklib.Config
}

func (cfg *Config) SetDefault() {
	if cfg.Version == "" {
		cfg.Version = deaultVersion
	}
}

// Init
func Init(cfg *Config, log mq.Logger, removeCfg bool) error {
	return kfklib.Init(&cfg.Config, log, nil, "", removeCfg)
}
