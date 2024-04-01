/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package spacerepositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain"
	spaceprimitive "github.com/openmerlin/merlin-server/space/domain/primitive"
	"github.com/openmerlin/merlin-server/space/domain/repository"
)

var (
	spaceEnvSecretTableName = ""
)

const (
	variableTypeName = "variable"
	secretTypeName   = "secret"

	fieldType    = "type"
	filedSpaceId = "space_id"
)

func toSpaceVariableDO(m *domain.SpaceVariable) spaceEnvSecretDO {
	resDesc := ""
	if m.Desc != nil {
		resDesc = m.Desc.MSDDesc()
	}
	return spaceEnvSecretDO{
		SpaceId:   m.SpaceId.Integer(),
		Desc:      resDesc,
		Name:      m.Name.MSDName(),
		Value:     m.Value.ENVValue(),
		Type:      variableTypeName,
		UpdatedAt: m.UpdatedAt,
	}
}

func toSpaceSecretDO(m *domain.SpaceSecret) spaceEnvSecretDO {
	resDesc := ""
	if m.Desc != nil {
		resDesc = m.Desc.MSDDesc()
	}
	return spaceEnvSecretDO{
		SpaceId:   m.SpaceId.Integer(),
		Desc:      resDesc,
		Name:      m.Name.MSDName(),
		Type:      secretTypeName,
		UpdatedAt: m.UpdatedAt,
	}
}

type spaceEnvSecretDO struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`
	SpaceId   int64  `gorm:"column:space_id;index:space_env_index,unique,priority:1"`
	Name      string `gorm:"column:name;index:space_env_index,unique,priority:2"`
	Desc      string `gorm:"column:desc"`
	Value     string `gorm:"column:value"`
	Type      string `gorm:"column:type"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

// TableName returns the table name of spaceDO.
func (do *spaceEnvSecretDO) TableName() string {
	return spaceEnvSecretTableName
}

func (do *spaceEnvSecretDO) toSpaceVariable() domain.SpaceVariable {
	return domain.SpaceVariable{
		Id:        primitive.CreateIdentity(do.Id),
		SpaceId:   primitive.CreateIdentity(do.SpaceId),
		Name:      primitive.CreateMSDName(do.Name),
		Desc:      primitive.CreateMSDDesc(do.Desc),
		Value:     spaceprimitive.CreateENVValue(do.Value),
		UpdatedAt: do.UpdatedAt,
	}
}

func (do *spaceEnvSecretDO) toSpaceSecret() domain.SpaceSecret {
	return domain.SpaceSecret{
		Id:        primitive.CreateIdentity(do.Id),
		SpaceId:   primitive.CreateIdentity(do.SpaceId),
		Name:      primitive.CreateMSDName(do.Name),
		Desc:      primitive.CreateMSDDesc(do.Desc),
		Value:     spaceprimitive.CreateENVValue(do.Value),
		UpdatedAt: do.UpdatedAt,
	}
}

func (do *spaceEnvSecretDO) toSpaceVariableSecretSummary() repository.SpaceVariableSecretSummary {
	return repository.SpaceVariableSecretSummary{
		Id:        primitive.CreateIdentity(do.Id).Identity(),
		Name:      do.Name,
		Value:     do.Value,
		Type:      do.Type,
		Desc:      do.Desc,
		UpdatedAt: do.UpdatedAt,
	}
}
