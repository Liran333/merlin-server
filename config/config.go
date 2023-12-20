package config

import (
	"os"

	common "github.com/openmerlin/merlin-server/common/config"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	gitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"
	modelctl "github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	userdomain "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
	redislib "github.com/opensourceways/redis-lib"
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

	API        controller.APIConfig `json:"api"`
	Git        gitea.Config         `json:"gitea"`
	Org        orgdomain.Config     `json:"organization"`
	User       userdomain.Config    `json:"user"`
	Model      modelConfig          `json:"model"`
	Redis      redis.Config         `json:"redis"`
	Mongodb    Mongodb              `json:"mongodb"`
	Authing    oidcimpl.Config      `json:"authing"`
	Primitive  primitive.Config     `json:"primitive"`
	Postgresql postgresql.Config    `json:"postgresql"`
}

func (cfg *Config) InitUserDomain() {
	userdomain.Init(&cfg.User)
}

func (cfg *Config) InitPrimitive() {
	primitive.Init(&cfg.Primitive)
}

func (cfg *Config) InitModel() {
	cfg.Model.initModel()
}

func (cfg *Config) GetRedisConfig() redislib.Config {
	return redislib.Config{
		DB:       cfg.Redis.DB,
		DBCert:   cfg.Redis.DBCert,
		Timeout:  cfg.Redis.Timeout,
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
	}
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.API,
		&cfg.Git,
		&cfg.Org,
		&cfg.User,
		&cfg.Model,
		&cfg.Redis,
		&cfg.Mongodb,
		&cfg.Authing,
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

type Mongodb struct {
	DBName      string             `json:"db_name"       required:"true"`
	DBConn      string             `json:"db_conn"       required:"true"`
	DBCert      string             `json:"db_cert"`
	Collections MongodbCollections `json:"collections"`
}

type MongodbCollections struct {
	User         string `json:"user"                   required:"true"`
	Session      string `json:"session"                required:"true"`
	Organization string `json:"organization"           required:"true"`
	Member       string `json:"member"                 required:"true"`
	Token        string `json:"token"                  required:"true"`
}

// modelConfig
type modelConfig struct {
	Tables     modelrepositoryadapter.Tables `json:"tables"`
	Controller modelctl.Config               `json:"controller"`
}

func (cfg *modelConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
		&cfg.Controller,
	}
}

func (cfg *modelConfig) initModel() {
	modelctl.Init(&cfg.Controller)
}
