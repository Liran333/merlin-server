/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	node           *snowflake.Node
	msdConfig      MSDConfig
	allLicenses    map[string]bool
	passwordLength int
	randomIdLength int
)

// Init initializes the configuration with the given Config struct.
func Init(cfg *Config) (err error) {
	msdConfig = cfg.MSDConfig

	m := map[string]bool{}
	for _, v := range cfg.Licenses {
		m[strings.ToLower(v)] = true
	}

	allLicenses = m

	// TODO: node id should be same with replica id
	node, err = snowflake.NewNode(1)

	passwordLength = cfg.PasswordLength
	randomIdLength = cfg.RandomIdLength

	if len(msdConfig.ReservedAccounts) > 0 {
		msdConfig.reservedAccounts = sets.New[string]()
		msdConfig.reservedAccounts.Insert(msdConfig.ReservedAccounts...)
	}
	return
}

// Config represents the main configuration structure.
type Config struct {
	MSDConfig

	Licenses       []string `json:"licenses" required:"true"`
	RandomIdLength int      `json:"random_id_length"`
	PasswordLength int      `json:"password_length"`
}

// MSDConfig represents the configuration for MSD.
type MSDConfig struct {
	MaxNameLength     int      `json:"max_name_length"`
	MinNameLength     int      `json:"min_name_length"`
	MaxDescLength     int      `json:"max_desc_length"`
	MaxFullnameLength int      `json:"max_fullname_length"`
	ReservedAccounts  []string `json:"reserved_accounts" required:"true"`
	reservedAccounts  sets.Set[string]
}

// SetDefault sets default values for MSDConfig if they are not provided.
func (cfg *Config) SetDefault() {
	if cfg.MaxNameLength <= 0 {
		cfg.MaxNameLength = 50
	}

	if cfg.MinNameLength <= 0 {
		cfg.MinNameLength = 5
	}

	if cfg.MaxDescLength <= 0 {
		cfg.MaxDescLength = 200
	}

	if cfg.MaxFullnameLength <= 0 {
		cfg.MaxFullnameLength = 200
	}

	if cfg.RandomIdLength <= 0 {
		cfg.RandomIdLength = 24
	}

	if cfg.PasswordLength <= 0 {
		cfg.PasswordLength = 32
	}
}
