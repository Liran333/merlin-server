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
	node             *snowflake.Node
	msdConfig        MSDConfig
	allLicenses      map[string]bool
	randomIdLength   int
	passwordInstance *passwordImpl
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

	randomIdLength = cfg.RandomIdLength
	passwordInstance = newPasswordImpl(cfg.PasswordConfig)

	if len(msdConfig.ReservedAccounts) > 0 {
		msdConfig.reservedAccounts = sets.New[string]()
		msdConfig.reservedAccounts.Insert(msdConfig.ReservedAccounts...)
	}
	return
}

// Config represents the main configuration structure.
type Config struct {
	MSDConfig

	Licenses       []string       `json:"licenses" required:"true"`
	RandomIdLength int            `json:"random_id_length"`
	PasswordConfig PasswordConfig `json:"password_config"`
}

// SetDefault sets default values for Config if they are not provided.
func (cfg *Config) SetDefault() {
	if cfg.RandomIdLength <= 0 {
		cfg.RandomIdLength = 24
	}
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.MSDConfig,
		&cfg.PasswordConfig,
	}
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
func (cfg *MSDConfig) SetDefault() {
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
}

// PasswordConfig represents the configuration for password.
type PasswordConfig struct {
	MinLength                int `json:"min_length"`
	MaxLength                int `json:"max_length"`
	MinNumOfCharKind         int `json:"min_num_of_char_kind"`
	MinNumOfConsecutiveChars int `json:"min_num_of_consecutive_chars"`
}

// SetDefault sets default values for PasswordConfig if they are not provided.
func (cfg *PasswordConfig) SetDefault() {
	if cfg.MinLength <= 0 {
		cfg.MinLength = 8
	}

	if cfg.MaxLength <= 0 {
		cfg.MaxLength = 20
	}

	if cfg.MinNumOfCharKind <= 0 {
		cfg.MinNumOfCharKind = 3
	}

	if cfg.MinNumOfConsecutiveChars <= 0 {
		cfg.MinNumOfConsecutiveChars = 2
	}
}
