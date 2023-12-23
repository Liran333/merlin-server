package controller

import (
	"fmt"

	"github.com/openmerlin/merlin-server/organization/domain"
)

type orgBasicInfoUpdateRequest struct {
	FullName    string `json:"full_name"`
	Website     string `json:"website"`
	AvatarId    string `json:"avatar_id"`
	Description string `json:"description"`
}

func (req *orgBasicInfoUpdateRequest) toCmd() (
	cmd domain.OrgUpdatedBasicInfoCmd,
	err error,
) {
	empty := true
	if req.FullName != "" {
		cmd.FullName = req.FullName
		empty = false
	}

	if req.AvatarId != "" {
		cmd.AvatarId = req.AvatarId
		empty = false
	}

	if req.Website != "" {
		cmd.Website = req.Website
		empty = false
	}

	if req.Description != "" {
		cmd.Description = req.Description
		empty = false
	}

	if empty {
		err = fmt.Errorf("edit org param can't be all empty")
	}

	return
}

type orgListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Invitee  string `form:"invitee"`
	Owner    string `form:"owner"`
	Website  string `form:"website"`
	Fullname string `form:"fullname"`
}

type orgCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Website     string `json:"website"`
	FullName    string `json:"full_name"`
	AvatarId    string `json:"avatar_id"`
	Description string `json:"description"`
}

func (req *orgCreateRequest) toCmd() (cmd domain.OrgCreatedCmd, err error) {
	if cmd.Name == "" {
		err = fmt.Errorf("org name can't be empty")
		return
	}

	if cmd.FullName == "" {
		err = fmt.Errorf("org full name can't be empty")
		return
	}

	err = cmd.Validate()

	return
}

type orgMemberAddRequest struct {
	User string `json:"user" binding:"required"`
}

type orgMemberRemoveRequest struct {
	User string `json:"user" binding:"required"`
}

type OrgMemberEditRequest struct {
	Role string `json:"role" binding:"required"`
	User string `json:"user" binding:"required"`
}

type OrgInviteMemberRequest struct {
	Role string `json:"role" binding:"required"`
	User string `json:"user" binding:"required"`
}

type OrgRevokeInviteRequest struct {
	User string `json:"user" binding:"required"`
}
