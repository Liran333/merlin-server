package main

import (
	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
)

type Config struct {
	Kafka      kafka.Config      `json:"kafka"`
	Postgresql postgresql.Config `json:"postgresql"`
	Model      modelConfig       `json:"model"`
	Gitea      gitea.Config      `json:"gitea"`
	Topics     Topics            `json:"topics"`
	Primitive  primitive.Config  `json:"primitive"`
	UserAgent  string            `json:"user_agent"`
}

type Topics struct {
	MerlinHookEvent string `json:"merlin_hook_event"  required:"true"`
}

type modelConfig struct {
	Tables modelrepositoryadapter.Tables `json:"tables"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return nil, err
	}

	cfg.SetDefault()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Kafka,
		&cfg.Postgresql,
		&cfg.Model,
		&cfg.Gitea,
		&cfg.Primitive,
	}
}

func (cfg *Config) SetDefault() {
	if cfg.UserAgent == "" {
		cfg.UserAgent = "Gitea-Hook-Delivery"
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(configValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) InitPrimitive() {
	primitive.Init(&cfg.Primitive)
}
