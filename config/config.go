package config

import (
	"os"

	redislib "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/coderepo"
	common "github.com/openmerlin/merlin-server/common/config"
	internal "github.com/openmerlin/merlin-server/common/controller/middleware/internalservice"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	gitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/kafka"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/models"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/permission"
	"github.com/openmerlin/merlin-server/session"
	"github.com/openmerlin/merlin-server/space"
	spaceapp "github.com/openmerlin/merlin-server/spaceapp"
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

	common.SetDefault(cfg)

	return common.Validate(cfg)
}

type Config struct {
	ReadHeaderTimeout int `json:"read_header_timeout"`

	Git        gitea.Config      `json:"gitea"`
	Org        orgdomain.Config  `json:"organization"`
	User       userdomain.Config `json:"user"`
	Redis      redislib.Config   `json:"redis"`
	Kafka      kafka.Config      `json:"kafka"`
	Model      models.Config     `json:"model"`
	Space      space.Config      `json:"space"`
	Session    session.Config    `json:"session"`
	SpaceApp   spaceapp.Config   `json:"space_app"`
	CodeRepo   coderepo.Config   `json:"coderepo"`
	Internal   internal.Config   `json:"internal"`
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

	cfg.CodeRepo.Init()

	internal.Init(&cfg.Internal)

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
		&cfg.Kafka,
		&cfg.Model,
		&cfg.Space,
		&cfg.Session,
		&cfg.SpaceApp,
		&cfg.CodeRepo,
		&cfg.Internal,
		&cfg.Primitive,
		&cfg.Postgresql,
	}
}

func (cfg *Config) SetDefault() {
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 10
	}
}

func (cfg *Config) Validate() error {
	return utils.CheckConfig(cfg, "")
}
