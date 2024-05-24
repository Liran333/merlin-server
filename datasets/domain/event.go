/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain for datasets.
package domain

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// datasetCreatedEvent
type datasetCreatedEvent struct {
	Time        int64  `json:"time"`
	Owner       string `json:"owner"`
	DatasetId   string `json:"dataset_id"`
	DatasetName string `json:"dataset_name"`
	CreatedBy   string `json:"created_by"`
}

func (e *datasetCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewDatasetCreatedEvent return a datasetCreatedEvent
func NewDatasetCreatedEvent(d *Dataset) datasetCreatedEvent {
	return datasetCreatedEvent{
		Time:        utils.Now(),
		Owner:       d.Owner.Account(),
		DatasetId:   d.Id.Identity(),
		DatasetName: d.Name.MSDName(),
		CreatedBy:   d.CreatedBy.Account(),
	}
}

// datasesDeletedEvent
type datasetDeletedEvent struct {
	Time        int64  `json:"time"`
	Owner       string `json:"owner"`
	DatasetId   string `json:"dataset_id"`
	DatasetName string `json:"dataset_name"`
	DeletedBy   string `json:"deleted_by"`
}

func (e *datasetDeletedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewDatasetDeletedEvent return a datasetDeletedEvent
func NewDatasetDeletedEvent(user primitive.Account, d Dataset) datasetDeletedEvent {
	return datasetDeletedEvent{
		Time:        utils.Now(),
		Owner:       d.Owner.Account(),
		DatasetId:   d.Id.Identity(),
		DatasetName: d.Name.MSDName(),
		DeletedBy:   user.Account(),
	}
}

// datasetUpdatedEvent
type datasetUpdatedEvent struct {
	Time        int64  `json:"time"`
	Repo        string `json:"repo"`
	Owner       string `json:"owner"`
	DatasetId   string `json:"dataset_id"`
	DatasetName string `json:"dataset_name"`
	UpdatedBy   string `json:"updated_by"`
	IsPriToPub  bool   `json:"is_pri_to_pub"` // private to public
}

func (e *datasetUpdatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewDatasetUpdatedEvent return a datasetUpdatedEvent
func NewDatasetUpdatedEvent(d *Dataset, user primitive.Account, b bool) datasetUpdatedEvent {
	return datasetUpdatedEvent{
		Time:        utils.Now(),
		Repo:        d.Name.MSDName(),
		Owner:       d.Owner.Account(),
		DatasetId:   d.Id.Identity(),
		DatasetName: d.Name.MSDName(),
		UpdatedBy:   user.Account(),
		IsPriToPub:  b,
	}
}
