/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/models/app"
)

type reqToResetLabel struct {
	Task        string   `yaml:"pipeline_tag"`
	Tags        []string `yaml:"tags"`
	License     string   `yaml:"license"`
	Frameworks  []string `yaml:"frameworks"`
	LibraryName string   `yaml:"library_name"`
}

func (req *reqToResetLabel) toCmd() app.CmdToResetLabels {
	cmd := app.CmdToResetLabels{}

	if config.tasks.Has(req.Task) {
		cmd.Task = req.Task
	}

	if config.libraryName.Has(req.LibraryName) {
		cmd.LibraryName = req.LibraryName
	}

	if len(req.Tags) > 0 {
		cmd.Others = sets.New[string](req.Tags...)
	}

	if len(req.License) > 0 {
		cmd.License = req.License
	}

	cmd.Frameworks = sets.New[string]()
	for _, framework := range req.Frameworks {
		if config.frameworks.Has(framework) {
			cmd.Frameworks.Insert(framework)
		}
	}

	return cmd
}
