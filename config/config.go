package config

import (
	"os"

	common "github.com/openmerlin/merlin-server/common/config"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	gitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/models"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/permission"
	"github.com/openmerlin/merlin-server/session"
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
	Redis      redislib.Config      `json:"redis"`
	Model      models.Config        `json:"model"`
	Session    session.Config       `json:"session"`
	Mongodb    Mongodb              `json:"mongodb"`
	Primitive  primitive.Config     `json:"primitive"`
	Postgresql postgresql.Config    `json:"postgresql"`
	Permission permission.Config    `json:"permission"`
}

func (cfg *Config) Init() {
	userdomain.Init(&cfg.User)

	primitive.Init(&cfg.Primitive)

	cfg.Model.Init()
}

func (cfg *Config) InitSession() error {
	return cfg.Session.Init()
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.API,
		&cfg.Git,
		&cfg.Org,
		&cfg.User,
		&cfg.Redis,
		&cfg.Model,
		&cfg.Session,
		&cfg.Mongodb,
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
	User          string `json:"user"                   required:"true"`
	Session       string `json:"session"                required:"true"`
	Organization  string `json:"organization"           required:"true"`
	Member        string `json:"member"                 required:"true"`
	Token         string `json:"token"                  required:"true"`
	Invitation    string `json:"invitation"             required:"true"`
	MemberRequest string `json:"member_request"         required:"true"`
}
