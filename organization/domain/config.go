package domain

type Config struct {
	MaxCountPerOwner int64  `json:"max_count_per_owner"`
	InviteExpiry     int64  `json:"invite_expiry"`
	DefaultRole      string `json:"default_role"`
	Tables           tables `json:"tables"`
}

type tables struct {
	Member string `json:"member" required:"true"`
	Invite string `json:"invite" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 10
	}
}
