package config

import (
	"os"

	common "github.com/openmerlin/merlin-server/common/config"
	"github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/controller"
	gitea "github.com/openmerlin/merlin-server/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"
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

	Authing oidcimpl.Config      `json:"authing"      required:"true"`
	Mongodb Mongodb              `json:"mongodb"      required:"true"`
	Redis   Redis                `json:"redis"        required:"true"`
	API     controller.APIConfig `json:"api"          required:"true"`
	User    userdomain.Config    `json:"user"         required:"true"`
	Git     gitea.Config         `json:"gitea"        required:"true"`
	Org     orgdomain.Config     `json:"organization"          required:"true"`
}

func (cfg *Config) GetRedisConfig() redislib.Config {
	return redislib.Config{
		Address:  cfg.Redis.DB.Address,
		Password: cfg.Redis.DB.Password,
		DB:       cfg.Redis.DB.DB,
		Timeout:  cfg.Redis.DB.Timeout,
		DBCert:   cfg.Redis.DB.DBCert,
	}
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Org,
		&cfg.Authing,
		&cfg.Mongodb,
		&cfg.Redis.DB,
		&cfg.Git,
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
	Collections MongodbCollections `json:"collections"   required:"true"`
}

type Redis struct {
	DB redis.Config `json:"db" required:"true"`
}

type MongodbCollections struct {
	User         string `json:"user"                   required:"true"`
	Session      string `json:"session"                required:"true"`
	Organization string `json:"organization"           required:"true"`
	Member       string `json:"member"                 required:"true"`
}
