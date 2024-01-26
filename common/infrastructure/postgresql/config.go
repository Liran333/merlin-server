package postgresql

import (
	"fmt"
	"time"
)

type Config struct {
	Host    string `json:"host"     required:"true"`
	User    string `json:"user"     required:"true"`
	Pwd     string `json:"pwd"      required:"true"`
	Name    string `json:"name"     required:"true"`
	Port    int    `json:"port"     required:"true"`
	Life    int    `json:"life"     required:"true"` // the unit is minute
	MaxConn int    `json:"max_conn" required:"true"`
	MaxIdle int    `json:"max_idle" required:"true"`
	Dbcert  string `json:"cert"`
}

func (p *Config) SetDefault() {
	if p.MaxConn <= 0 {
		p.MaxConn = 500
	}

	if p.MaxIdle <= 0 {
		p.MaxIdle = 250
	}

	if p.Life <= 0 {
		p.Life = 2
	}
}

func (cfg *Config) getLifeDuration() time.Duration {
	return time.Minute * time.Duration(cfg.Life)
}

func (p *Config) dsn() string {
	sslmode := "disable"
	if p.Dbcert != "" {
		sslmode = "require"
	}
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=Asia/Shanghai",
		p.Host, p.User, p.Pwd, p.Name, p.Port, sslmode,
	)
}
