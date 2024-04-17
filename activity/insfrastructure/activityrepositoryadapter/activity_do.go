/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package activityrepositoryadapter provides an adapter for the model repository
package activityrepositoryadapter

import (
	"github.com/openmerlin/merlin-server/activity/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var (
	ActiviyTableName = ""
)

// TableName returns the table name of the model.
func (do *activityDO) TableName() string {
	return ActiviyTableName
}

const (
	fieldSpace         = "space"
	fieldModel         = "model"
	fieldTime          = "time"
	fieldLike          = "like"
	fieldTypeOwner     = "owner"
	fieldResourceType  = "resource_type"
	fieldResourceIndex = "resource_id"
	fieldType          = "type"
)

type activityDO struct {
	AutoID        uint   `gorm:"primaryKey;autoIncrement"`
	Owner         string `gorm:"column:owner"`
	Type          string `gorm:"column:type"`
	Time          int64  `gorm:"column:time"`
	ResourceIndex int64  `gorm:"column:resource_id"`
	ResourceType  string `gorm:"column:resource_type"`
}

func convertToActivityDomain(d activityDO) (domain.Activity, error) {
	return domain.Activity{
		Owner: primitive.CreateAccount(d.Owner),
		Type:  domain.ActivityType(d.Type),
		Time:  d.Time,
		Resource: domain.Resource{
			Type:  primitive.ObjType(d.ResourceType),
			Index: primitive.CreateIdentity(d.ResourceIndex),
		},
	}, nil
}
