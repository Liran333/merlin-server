/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain for models.
package domain

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// modelCreatedEvent
type modelCreatedEvent struct {
	Time      int64  `json:"time"`
	Owner     string `json:"owner"`
	ModelId   string `json:"model_id"`
	CreatedBy string `json:"created_by"`
}

func (e *modelCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewModelCreatedEvent return a modelCreatedEvent
func NewModelCreatedEvent(m *Model) modelCreatedEvent {
	return modelCreatedEvent{
		Time:      utils.Now(),
		Owner:     m.Owner.Account(),
		ModelId:   m.Id.Identity(),
		CreatedBy: m.CreatedBy.Account(),
	}
}

// modelDeletedEvent
type modelDeletedEvent struct {
	ModelId   string `json:"model_id"`
	DeletedBy string `json:"deleted_by"`
}

func (e *modelDeletedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewModelDeletedEvent return a modelDeletedEvent
func NewModelDeletedEvent(user primitive.Account, modelId primitive.Identity) modelDeletedEvent {
	return modelDeletedEvent{
		ModelId:   modelId.Identity(),
		DeletedBy: user.Account(),
	}
}

// modelUpdatedEvent
type modelUpdatedEvent struct {
	Time       int64  `json:"time"`
	Repo       string `json:"repo"`
	Owner      string `json:"owner"`
	ModelId    string `json:"model_id"`
	UpdatedBy  string `json:"updated_by"`
	IsPriToPub bool   `json:"is_pri_to_pub"` // private to public
}

func (e *modelUpdatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewModelUpdatedEvent return a modelUpdatedEvent
func NewModelUpdatedEvent(m *Model, user primitive.Account, b bool) modelUpdatedEvent {
	return modelUpdatedEvent{
		Time:       utils.Now(),
		Repo:       m.Name.MSDName(),
		Owner:      m.Owner.Account(),
		ModelId:    m.Id.Identity(),
		UpdatedBy:  user.Account(),
		IsPriToPub: b,
	}
}

// modelDisableEvent
type modelDisableEvent struct {
	Time      int64  `json:"time"`
	Repo      string `json:"repo"`
	Owner     string `json:"owner"`
	ModelId   string `json:"model_id"`
	UpdatedBy string `json:"updated_by"`
}

// Message serializes the modelDisableEvent into a JSON byte array.
func (e *modelDisableEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewModelDisableEvent creates a new modelDisableEvent instance with the given Space.
func NewModelDisableEvent(m *Model, user primitive.Account) modelDisableEvent {
	return modelDisableEvent{
		Time:      utils.Now(),
		Repo:      m.Name.MSDName(),
		Owner:     m.Owner.Account(),
		ModelId:   m.Id.Identity(),
		UpdatedBy: user.Account(),
	}
}
