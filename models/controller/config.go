/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

var config Config

// Init initializes the configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config represents the application configuration.
type Config struct {
	Tasks           []string `json:"tasks"               required:"true"`
	Frameworks      []string `json:"frameworks"          required:"true"`
	LibraryName     []string `json:"library_name"        required:"true"`
	Languages       []string `json:"languages"`
	Hardwares       []string `json:"hardwares"`
	MaxCountPerPage int      `json:"max_count_per_page"`

	tasks       sets.Set[string]
	frameworks  sets.Set[string]
	libraryName sets.Set[string]
	languages   sets.Set[string]
	hardwares   sets.Set[string]
}

// SetDefault sets the default values for the configuration.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerPage <= 0 {
		cfg.MaxCountPerPage = 100
	}
}

// Validate check values for Config whether they are valid.
func (cfg *Config) Validate() (err error) {
	if len(cfg.Tasks) > 0 {
		cfg.tasks = sets.New[string](cfg.Tasks...)
	}

	if len(cfg.Frameworks) > 0 {
		cfg.frameworks = sets.New[string](cfg.Frameworks...)
	}

	if len(cfg.Languages) > 0 {
		cfg.languages = sets.New[string](cfg.Languages...)
	}

	if len(cfg.Hardwares) > 0 {
		cfg.hardwares = sets.New[string](cfg.Hardwares...)
	}

	if len(cfg.LibraryName) > 0 {
		cfg.libraryName = sets.New[string](cfg.LibraryName...)
	}

	return
}
