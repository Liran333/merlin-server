/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package csrftokenrepositoryadapter

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain"
)

func toCSRFTokenDO(v *domain.CSRFToken) csrfTokenDO {
	return csrfTokenDO{
		Expiry:  v.Expiry,
		HasUsed: v.HasUsed,
		LoginId: v.LoginId,
	}
}

type csrfTokenDO struct {
	Expiry  int64          `json:"expiry"`
	HasUsed bool           `json:"has_used"`
	LoginId primitive.UUID `json:"login_id"`
}

// MarshalBinary in order to store struct directly in redis
func (do *csrfTokenDO) MarshalBinary() ([]byte, error) {
	return json.Marshal(do)
}

// UnmarshalBinary unmarshals the binary data into the csrfTokenDO struct.
func (do *csrfTokenDO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, do)
}

func (do *csrfTokenDO) toCSRFToken(tid primitive.UUID) domain.CSRFToken {
	return domain.CSRFToken{
		Id:      tid,
		Expiry:  do.Expiry,
		HasUsed: do.HasUsed,
		LoginId: do.LoginId,
	}
}
