/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

const (
	fieldSpaceId = "space_id"
	fieldModelId = "model_id"
)

var (
	spaceModelRelationTableName = ""
)

type modelSpaceRelationDO struct {
	Id      int64 `gorm:"primarykey"`
	SpaceId int64 `gorm:"column:space_id;index:space_id_index"`
	ModelId int64 `gorm:"column:model_id;index:model_id_index"`
}

// TableName returns the table name of spaceDO.
func (do *modelSpaceRelationDO) TableName() string {
	return spaceModelRelationTableName
}
