/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-sdk/activityapp"
	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

func ConvertReqToCreateActivityToCmd(req *activityapp.ReqToCreateActivity) (app.CmdToAddActivity, error) {
	var cmd app.CmdToAddActivity

	timeInt, err := strconv.ParseInt(req.Time, 10, 64)
	if err != nil {
		return cmd, err
	}

	resourceIdInt, err := strconv.ParseInt(req.ResourceId, 10, 64)
	if err != nil {
		logrus.Errorf("failed to convert to int: %s", err)
		return cmd, err
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

func ConvertReqToDeleteActivityToCmd(req *activityapp.ReqToDeleteActivity) (app.CmdToAddActivity, error) {
	var cmd app.CmdToAddActivity

	resourceIdInt, err := strconv.ParseInt(req.ResourceId, 10, 64)
	if err != nil {
		logrus.Errorf("failed to convert to int: %s", err)
		return cmd, err
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
