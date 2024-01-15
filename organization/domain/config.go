package domain

type Config struct {
	InviteExpiry int64  `json:"invite_expiry"`
	DefaultRole  string `json:"default_role"`
	Tables       tables `json:"tables"`
}

type tables struct {
	Member string `json:"member" required:"true"`
	Invite string `json:"invite" required:"true"`
}
