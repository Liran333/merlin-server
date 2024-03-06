/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package sessionrepositoryadapter

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

func toSessionDO(v *domain.Session) sessionDO {
	return sessionDO{
		IP:        v.IP,
		User:      v.User.Account(),
		UserId:    v.UserId,
		IdToken:   v.IdToken,
		UserAgent: v.UserAgent.UserAgent(),
		CreatedAt: v.CreatedAt,
	}
}

type sessionDO struct {
	IP        string `json:"ip"`
	User      string `json:"user"`
	UserId    string `json:"user_id"`
	IdToken   string `json:"id_token"`
	UserAgent string `json:"user_agent"`
	CreatedAt int64  `json:"created_at"`
}

// MarshalBinary in order to store struct directly in redis
func (do *sessionDO) MarshalBinary() ([]byte, error) {
	return json.Marshal(do)
}

func (do *sessionDO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, do)
}

func (do *sessionDO) toSession(id primitive.RandomId) domain.Session {
	return domain.Session{
		Id:        id,
		IP:        do.IP,
		User:      primitive.CreateAccount(do.User),
		UserId:    do.UserId,
		IdToken:   do.IdToken,
		UserAgent: primitive.CreateUserAgent(do.UserAgent),
		CreatedAt: do.CreatedAt,
	}
}
