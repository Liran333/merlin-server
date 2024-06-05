/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"k8s.io/apimachinery/pkg/util/sets"

	sdk "github.com/openmerlin/merlin-sdk/models"

	"github.com/openmerlin/merlin-server/models/app"
)

type reqToResetLabel sdk.ReqToResetLabel

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

	cmd.Hardwares = sets.New[string]()
	for _, hardware := range req.Hardwares {
		if config.hardwares.Has(hardware) {
			cmd.Hardwares.Insert(hardware)
		}
	}

	if len(req.Languages) > 0 {
		cmd.Languages = sets.New[string](req.Languages...)
	}

	return cmd
}
