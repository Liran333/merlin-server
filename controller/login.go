package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type oldUserTokenPayload struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

func (pl *oldUserTokenPayload) DomainAccount() domain.Account {
	return primitive.CreateAccount(pl.Account)
}

func (pl *oldUserTokenPayload) isNotMe(a domain.Account) bool {
	return pl.Account != a.Account()
}

func (pl *oldUserTokenPayload) isMyself(a domain.Account) bool {
	return pl.Account == a.Account()
}

func (pl *oldUserTokenPayload) hasEmail() bool {
	return pl.Email != ""
}
