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

type Subject struct {
	Identity string  `json:"id"`
	Type     ObjType `json:"type"`
}

type Object = Subject

type Action int

const (
	ActionRead Action = iota
	ActionWrite
	ActionDelete
	ActionCreate
)
