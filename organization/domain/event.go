/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides domain organization and configuration for the app service.
package domain

import "encoding/json"

// userJoinEvent
type userJoinEvent struct {
	OrgName   string `json:"org_name"`
	UserName  string `json:"user_name"`
	CreatedAt int64  `json:"created_at"`
}

// Message returns the JSON representation of the userJoinEvent.
func (e *userJoinEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewUserJoinEvent creates a new userJoinEvent instance with the given Organization.
func NewUserJoinEvent(a *Approve) userJoinEvent {
	return userJoinEvent{
		OrgName:   a.OrgName.Account(),
		UserName:  a.Username.Account(),
		CreatedAt: a.CreatedAt,
	}
}

// userRemoveEvent
type userRemoveEvent struct {
	OrgName  string `json:"org_name"`
	UserName string `json:"user_name"`
}

// Message returns the JSON representation of the userRemoveEvent.
func (e *userRemoveEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewUserRemoveEvent creates a new userRemoveEvent instance with the given Organization.
func NewUserRemoveEvent(a *OrgRemoveMemberCmd) userRemoveEvent {
	return userRemoveEvent{
		OrgName:  a.Org.Account(),
		UserName: a.Account.Account(),
	}
}

// orgDeleteEvent
type orgDeleteEvent struct {
	OrgName string `json:"org_name"`
}

// Message returns the JSON representation of the orgDeleteEvent.
func (e *orgDeleteEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewOrgDeleteEvent creates a new orgDeleteEvent instance with the given Organization.
func NewOrgDeleteEvent(a *OrgDeletedCmd) orgDeleteEvent {
	return orgDeleteEvent{
		OrgName: a.Name.Account(),
	}
}
