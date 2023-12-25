package primitive

type ObjType string

const (
	ObjTypeUser    ObjType = "user"
	ObjTypeOrg     ObjType = "organization"
	ObjTypeModel   ObjType = "model"
	ObjTypeDataset ObjType = "dataset"
	ObjTypeSpace   ObjType = "space"
	ObjTypeMember  ObjType = "member"
	ObjTypeInvite  ObjType = "invite"
)

type Action int

const (
	ActionRead Action = iota
	ActionWrite
	ActionDelete
	ActionCreate
)

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
	}
	return ""
}
