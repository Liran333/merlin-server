package postgresql

import (
	"fmt"
	"time"
)

type Config struct {
	Host    string    `json:"host"     required:"true"`
	User    string    `json:"user"     required:"true"`
	Pwd     string    `json:"pwd"      required:"true"`
	Name    string    `json:"name"     required:"true"`
	Port    int       `json:"port"     required:"true"`
	Life    int       `json:"life"     required:"true"` // the unit is minute
	MaxConn int       `json:"max_conn" required:"true"`
	MaxIdle int       `json:"max_idle" required:"true"`
	Dbcert  string    `json:"cert"`
	Code    errorCode `json:"error_code"`
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

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Code,
	}
}

func (cfg *Config) getLifeDuration() time.Duration {
	return time.Minute * time.Duration(cfg.Life)
}

func (p *Config) dsn() string {
	if p.Dbcert != "" {
		return fmt.Sprintf(
			"host=%v user=%v password=%v dbname=%v port=%v sslmode=verify-ca TimeZone=Asia/Shanghai sslrootcert=%v",
			p.Host, p.User, p.Pwd, p.Name, p.Port, p.Dbcert,
		)
	} else {
		return fmt.Sprintf(
			"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai",
			p.Host, p.User, p.Pwd, p.Name, p.Port,
		)
	}
}

type errorCode struct {
	UniqueConstraint string `json:"unique_constraint"`
}

func (cfg *errorCode) SetDefault() {
	if cfg.UniqueConstraint == "" {
		cfg.UniqueConstraint = "23505"
	}
}
