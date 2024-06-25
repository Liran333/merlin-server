/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for the application.
package app

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

var config Config

// Init initializes the application with the provided configuration.
func Init(cfg *Config) {
	if len(cfg.AvatarIds) > 0 {
		cfg.avatarIdsSet = sets.New[string]()
		cfg.avatarIdsSet.Insert(cfg.AvatarIds...)
	} else {
		logrus.Fatal("avatar ids is empty")
	}

	config = *cfg
}

// Config is a struct that holds the configuration for max count per owner.
type Config struct {
	AvatarIds             []string         `json:"avatar_ids" required:"true"`
	ObsPath               string           `json:"obs_path" required:"true"`
	ObsBucket             string           `json:"obs_bucket" required:"true"`
	CdnEndpoint           string           `json:"cdn_endpoint" required:"true"`
	MaxCountPerUser       int              `json:"max_count_per_user"`
	MaxCountPerOrg        int              `json:"max_count_per_org"`
	MaxCountSpaceSecret   int              `json:"max_count_space_secret"`
	MaxCountSpaceVariable int              `json:"max_count_space_variable"`
	RecommendSpaces       []RecommendIndex `json:"recommend_spaces"`
	BoutiqueSpaces        []BoutiqueIndex  `json:"boutique_spaces"`

	avatarIdsSet sets.Set[string]
}

type RecommendIndex struct {
	Owner    string `json:"owner" required:"true"`
	Reponame string `json:"reponame" required:"true"`
}

type BoutiqueIndex struct {
	Owner    string `json:"owner" required:"true"`
	Reponame string `json:"reponame" required:"true"`
}

// SetDefault sets the default values for the Config struct.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerUser <= 0 {
		cfg.MaxCountPerUser = 50
	}

	if cfg.MaxCountPerOrg <= 0 {
		cfg.MaxCountPerOrg = 200
	}

	if cfg.MaxCountSpaceVariable <= 0 {
		cfg.MaxCountSpaceVariable = 100
	}

	if cfg.MaxCountSpaceSecret <= 0 {
		cfg.MaxCountSpaceSecret = 100
	}
}
