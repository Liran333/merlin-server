/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"strconv"

	"github.com/openmerlin/merlin-sdk/activityapp"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// const for parse
const (
	parseIntBase    = 10
	parseIntBitSize = 64
)

// ConvertReqToCreateActivityToCmd converts the request to the command
func ConvertReqToCreateActivityToCmd(req *activityapp.ReqToCreateActivity) (app.CmdToAddActivity, error) {
	var cmd app.CmdToAddActivity

	timeInt, err := strconv.ParseInt(req.Time, parseIntBase, parseIntBitSize)
	if err != nil {
		return cmd, xerrors.Errorf("failed to convert to int: %w", err)
	}

	resourceIdInt, err := strconv.ParseInt(req.ResourceId, parseIntBase, parseIntBitSize)
	if err != nil {
		return cmd, xerrors.Errorf("failed to convert to int: %w", err)
	}

	resource := domain.Resource{
		Type:  primitive.ObjType(req.ResourceType),
		Index: primitive.CreateIdentity(resourceIdInt),
	}

	cmd = app.CmdToAddActivity{
		Owner:    primitive.CreateAccount(req.Owner),
		Type:     domain.ActivityType(req.Type),
		Time:     timeInt,
		Resource: resource,
	}

	return cmd, nil
}

// ConvertReqToDeleteActivityToCmd converts a request to delete an activity to an app.CmdToAddActivity command.
func ConvertReqToDeleteActivityToCmd(
	user primitive.Account, req *activityapp.ReqToDeleteActivity) (
	app.CmdToAddActivity, error) {
	var cmd app.CmdToAddActivity

	resourceIdInt, err := strconv.ParseInt(req.ResourceId, parseIntBase, parseIntBitSize)
	if err != nil {
		return cmd, xerrors.Errorf("failed to convert to int: %w", err)
	}

	resource := domain.Resource{
		Type:  primitive.ObjType(req.ResourceType),
		Index: primitive.CreateIdentity(resourceIdInt),
	}

	cmd = app.CmdToAddActivity{
		Owner:    user,
		Resource: resource,
	}

	return cmd, nil
}

// ConvertInternalReqToDeleteActivityToCmd converts an internal request to delete an activity to an app.CmdToAddActivity command.
func ConvertInternalReqToDeleteActivityToCmd(req *activityapp.ReqToDeleteActivity) (app.CmdToAddActivity, error) {
	var cmd app.CmdToAddActivity

	resourceIdInt, err := strconv.ParseInt(req.ResourceId, parseIntBase, parseIntBitSize)
	if err != nil {
		return cmd, xerrors.Errorf("failed to convert to int: %w", err)
	}

	resource := domain.Resource{
		Type:  primitive.ObjType(req.ResourceType),
		Index: primitive.CreateIdentity(resourceIdInt),
	}

	cmd = app.CmdToAddActivity{
		Resource: resource,
	}

	return cmd, nil
}
