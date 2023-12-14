package session

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
	"github.com/openmerlin/merlin-server/infrastructure/repositories"
)

const (
	fieldName    = "name"
	fieldCount   = "count"
	fieldAccount = "account"
)

func sessionDocFilter(account string) bson.M {
	return bson.M{
		fieldAccount: account,
	}
}

func NewSessionStore(store mongodbClient) sessionStore {
	return sessionStore{store}
}

type sessionStore struct {
	cli mongodbClient
}

type dLogin struct {
	Account string `bson:"account"   json:"account"`
	Info    string `bson:"info"      json:"info"`
	Email   string `bson:"email"    json:"email"`
	UserId  string `bson:"user_id"   json:"user_id"`
}

func (se sessionStore) Insert(do SessionDO) error {
	doc, err := se.toLoginDoc(&do)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := se.cli.ReplaceDoc(
			ctx,
			sessionDocFilter(do.Account), doc,
		)

		return err
	}

	return primitive.WithContext(f)
}

func (se sessionStore) Get(account string) (do SessionDO, err error) {
	var v dLogin

	f := func(ctx context.Context) error {
		return se.cli.GetDoc(
			ctx, sessionDocFilter(account), nil, &v,
		)
	}

	if err = primitive.WithContext(f); err == nil {
		se.toSessionDO(&v, &do)

		return
	}

	if isDocNotExists(err) {
		err = repositories.NewErrorDataNotExists(err)
	}

	return
}

func (se sessionStore) toLoginDoc(do *SessionDO) (bson.M, error) {
	docObj := dLogin{
		Account: do.Account,
		Info:    do.Info,
		Email:   do.Email,
		UserId:  do.UserId,
	}

	return mongo.GenDoc(docObj)
}

func (se sessionStore) toSessionDO(u *dLogin, do *SessionDO) {
	*do = SessionDO{
		Account: u.Account,
		Info:    u.Info,
		Email:   u.Email,
		UserId:  u.UserId,
	}
}
