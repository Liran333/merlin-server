package config

import (
	"os"

	redislib "github.com/opensourceways/redis-lib"

	common "github.com/openmerlin/merlin-server/common/config"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	gitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/models"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/permission"
	"github.com/openmerlin/merlin-server/session"
	"github.com/openmerlin/merlin-server/space"
	userdomain "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

func LoadConfig(path string, cfg *Config, remove bool) error {
	if remove {
		defer os.Remove(path)
	}

	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}

type Config struct {
	ReadHeaderTimeout int `json:"read_header_timeout"`

	Git        gitea.Config      `json:"gitea"`
	Org        orgdomain.Config  `json:"organization"`
	User       userdomain.Config `json:"user"`
	Redis      redislib.Config   `json:"redis"`
	Model      models.Config     `json:"model"`
	Space      space.Config      `json:"space"`
	Session    session.Config    `json:"session"`
	Primitive  primitive.Config  `json:"primitive"`
	Postgresql postgresql.Config `json:"postgresql"`
	Permission permission.Config `json:"permission"`
}

func (cfg *Config) Init() error {
	userdomain.Init(&cfg.User)

	if err := primitive.Init(&cfg.Primitive); err != nil {
		return err
	}

	cfg.Model.Init()

	cfg.Space.Init()

	return nil
}

func (cfg *Config) InitSession() error {
	return cfg.Session.Init()
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Git,
		&cfg.Org,
		&cfg.User,
		&cfg.Redis,
		&cfg.Model,
		&cfg.Space,
		&cfg.Session,
		&cfg.Primitive,
		&cfg.Postgresql,
	}
}

func (cfg *Config) setDefault() {
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 10
	}

	common.SetDefault(cfg)
}

func (cfg *Config) validate() error {
	if err := utils.CheckConfig(cfg, ""); err != nil {
		return err
	}

	return common.Validate(cfg)
}
