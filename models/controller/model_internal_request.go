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
	Task       string   `yaml:"pipeline_tag"`
	Tags       []string `yaml:"tags"`
	License    string   `yaml:"license"`
	Frameworks []string `yaml:"frameworks"`
}

func (req *reqToResetLabel) toCmd() app.CmdToResetLabels {
	cmd := app.CmdToResetLabels{
		Task: req.Task,
	}

	if len(req.Tags) > 0 {
		cmd.Others = sets.New[string](req.Tags...)
	}

	if len(req.License) > 0 {
		cmd.License = req.License
	}

	if len(req.Frameworks) > 0 {
		cmd.Frameworks = sets.New[string](req.Frameworks...)
	}

	return cmd
}
