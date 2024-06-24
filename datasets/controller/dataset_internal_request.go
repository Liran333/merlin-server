/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openmerlin/merlin-server/datasets/app"
)

type reqToResetDatasetLabel struct {
	Task     []string `yaml:"task"`
	Licenses []string `yaml:"license"`
	Size     string   `yaml:"size"`
	Language []string `yaml:"language"`
	Domain   []string `yaml:"domain"`
}

func (req *reqToResetDatasetLabel) toCmd() app.CmdToResetLabels {
	cmd := app.CmdToResetLabels{}

	if len(req.Task) > 0 {
		cmd.Task = sets.New[string](req.Task...)
	}

	if len(req.Licenses) > 0 {
		cmd.License = sets.New[string](req.Licenses...)
	}

	if len(req.Size) > 0 {
		cmd.Size = req.Size
	}

	if len(req.Language) > 0 {
		cmd.Language = sets.New[string](req.Language...)
	}

	if len(req.Domain) > 0 {
		cmd.Domain = sets.New[string](req.Domain...)
	}

	return cmd
}
