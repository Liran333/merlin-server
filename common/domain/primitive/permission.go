package primitive

import "errors"

type ObjType string

const (
	ObjTypeUser    ObjType = "user"
	ObjTypeOrg     ObjType = "organization"
	ObjTypeModel   ObjType = "model"
	ObjTypeDataset ObjType = "dataset"
	ObjTypeSpace   ObjType = "space"
	ObjTypeMember  ObjType = "member"
	ObjTypeInvite  ObjType = "invite"

	TokenPermWrite string = "write"
	TokenPermRead  string = "read"

	ActionRead Action = iota
	ActionWrite
	ActionDelete
	ActionCreate
)

type Action int

func (a Action) String() string {
	switch a {
	case ActionRead:
		return "read"
	case ActionWrite:
		return "write"
	case ActionDelete:
		return "delete"
	case ActionCreate:
		return "create"
	default:
		return ""
	}
}

func (a Action) IsModification() bool {
	return a == ActionDelete || a == ActionWrite
}

type tokenPerm string

func (r tokenPerm) TokenPerm() string {
	return string(r)
}

func (t tokenPerm) PermissionAllow(expect TokenPerm) bool {
	if expect.TokenPerm() == TokenPermRead {
		return true
	}

	if expect.TokenPerm() == TokenPermWrite {
		return t.TokenPerm() == TokenPermWrite
	}

	return false
}

type TokenPerm interface {
	TokenPerm() string
	PermissionAllow(expect TokenPerm) bool
}

func NewTokenPerm(v string) (TokenPerm, error) {
	if v != TokenPermWrite && v != TokenPermRead {
		return nil, errors.New("invalid permission")
	}

	return tokenPerm(v), nil
}

func NewReadPerm() TokenPerm {
	return tokenPerm(TokenPermRead)
}

func NewWritePerm() TokenPerm {
	return tokenPerm(TokenPermWrite)
}

func CreateTokenPerm(v string) TokenPerm {
	return tokenPerm(v)
}
