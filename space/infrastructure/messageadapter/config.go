package messageadapter

type Topics struct {
	SpaceDeleted string `json:"space_deleted" required:"true"`
	SpaceUpdated string `json:"space_updated" required:"true"`
}
