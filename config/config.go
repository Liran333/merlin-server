/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package config provides functionality for managing application configuration.
package config

import (
	"os"

	redislib "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/activity"
	"github.com/openmerlin/merlin-server/coderepo"
	common "github.com/openmerlin/merlin-server/common/config"
	internal "github.com/openmerlin/merlin-server/common/controller/middleware/internalservice"
	"github.com/openmerlin/merlin-server/common/controller/middleware/ratelimiter"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitaccess"
	gitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/kafka"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/common/infrastructure/securestorage"
	"github.com/openmerlin/merlin-server/computility"
	"github.com/openmerlin/merlin-server/models"
	orgdomain "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/domain/permission"
	"github.com/openmerlin/merlin-server/organization/domain/privilege"
	"github.com/openmerlin/merlin-server/other"
	"github.com/openmerlin/merlin-server/session"
	"github.com/openmerlin/merlin-server/space"
	spaceapp "github.com/openmerlin/merlin-server/spaceapp"
	userdomain "github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

// LoadConfig loads the configuration file from the specified path and deletes the file if needed
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

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	ReadHeaderTimeout   int  `json:"read_header_timeout"`
	NeedTokenForEachAPI bool `json:"need_token_for_each_api"`

	Git          gitea.Config         `json:"gitea"`
	Org          orgdomain.Config     `json:"organization"`
	User         userdomain.Config    `json:"user"`
	Vault        securestorage.Config `json:"vault"`
	Redis        redislib.Config      `json:"redis"`
	Kafka        kafka.Config         `json:"kafka"`
	Model        models.Config        `json:"model"`
	Space        space.Config         `json:"space"`
	Activity     activity.Config      `json:"activity"`
	Session      session.Config       `json:"session"`
	SpaceApp     spaceapp.Config      `json:"space_app"`
	CodeRepo     coderepo.Config      `json:"coderepo"`
	Internal     internal.Config      `json:"internal"`
	Primitive    primitive.Config     `json:"primitive"`
	GitAccess    gitaccess.Config     `json:"git_access"`
	Postgresql   postgresql.Config    `json:"postgresql"`
	Permission   permission.Config    `json:"permission"`
	RateLimiter  ratelimiter.Config   `json:"ratelimit"`
	Computility  computility.Config   `json:"computility"`
	OtherConfig  other.Config         `json:"other_config"`
	PrivilegeOrg privilege.Config     `json:"privilege_org"`
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() error {
	if err := primitive.Init(&cfg.Primitive); err != nil {
		return err
	}

	cfg.Model.Init()

	cfg.Space.Init()

	cfg.SpaceApp.Init()

	cfg.CodeRepo.Init()

	cfg.Activity.Init()

	internal.Init(&cfg.Internal)

	return nil
}

// InitSession initializes the session associated with the configuration.
func (cfg *Config) InitSession() error {
	return cfg.Session.Init()
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
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
		&cfg.GitAccess,
		&cfg.Primitive,
		&cfg.Postgresql,
		&cfg.Vault,
		&cfg.OtherConfig,
	}
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = 10
	}
}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	return utils.CheckConfig(cfg, "")
}
