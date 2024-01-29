package primitive

var (
	maxBranchNameLength int
)

func Init(cfg *Config) {
	maxBranchNameLength = cfg.MaxBranchNameLength
}

type Config struct {
	MaxBranchNameLength int `json:"max_branch_name_length" `
}

func (cfg *Config) SetDefault() {
	if cfg.MaxBranchNameLength <= 0 {
		cfg.MaxBranchNameLength = 100
	}
}
