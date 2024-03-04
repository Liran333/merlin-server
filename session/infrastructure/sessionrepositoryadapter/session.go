/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package sessionrepositoryadapter provides an adapter for the session repository,
// implementing the necessary methods to interact with the session data storage.
package sessionrepositoryadapter

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	"github.com/openmerlin/merlin-server/session/domain"
)

const sessionKeyPrefix = "session_"

type dao interface {
	Get(key string, val interface{}) error
	Expire(key string, expire time.Duration) error
	SetWithExpiry(key string, val interface{}, expiry time.Duration) error
	IsKeyNotExists(err error) bool
}

// NewSessionAdapter creates a new session adapter with the provided DAO and returns a pointer to it.
func NewSessionAdapter(d dao) *sessionAdapter {
	return &sessionAdapter{
		dao: d,
	}
}

type sessionAdapter struct {
	dao dao
}

func (adapter *sessionAdapter) Add(t *domain.Session) error {
	v := toSessionDO(t)

	// must pass *sessionDO, because it implements the interface of MarshalBinary
	return adapter.dao.SetWithExpiry(adapter.generateKey(t.Id), &v, t.LifeTime())
}

func (adapter *sessionAdapter) Save(t *domain.Session) error {
	return adapter.Add(t)
}

func (adapter *sessionAdapter) Find(id primitive.RandomId) (domain.Session, error) {
	var do sessionDO

	// note the *sessionDO implements interface of UnmarshalBinary
	if err := adapter.dao.Get(adapter.generateKey(id), &do); err != nil {
		if adapter.dao.IsKeyNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return domain.Session{}, err
	}

	return do.toSession(id), nil
}

func (adapter *sessionAdapter) Delete(id primitive.RandomId) error {
	return adapter.dao.Expire(adapter.generateKey(id), 0)
}

func (adapter *sessionAdapter) generateKey(id primitive.RandomId) string {
	return sessionKeyPrefix + id.RandomId()
}
