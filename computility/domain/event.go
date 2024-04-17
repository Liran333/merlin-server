/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides domain organization and configuration for the app service.
package domain

import "encoding/json"

type computeRecallEvent struct {
	UserName    string `json:"user_name"`
	QuotaCount  int    `json:"quota_count"`
	ComputeType string `json:"compute_type"`
}

type computeRecallEventList struct {
	RecallList []computeRecallEvent `json:"recall_list"`
}

// Message returns the JSON representation of the computeRecallEventList.
func (e *computeRecallEventList) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewcomputeRecallEvent creates a new computeRecallEventList instance with the given RecallInfoList.
func NewcomputeRecallEvent(w *RecallInfoList) computeRecallEventList {
	list := []computeRecallEvent{}
	for _, v := range w.InfoList {
		list = append(list, computeRecallEvent{
			UserName:    v.UserName.Account(),
			QuotaCount:  v.QuotaCount,
			ComputeType: v.ComputeType.ComputilityType(),
		})
	}

	return computeRecallEventList{
		RecallList: list,
	}
}
