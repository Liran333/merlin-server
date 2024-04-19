package email

type Config struct {
	AuthCode string `json:"auth_code" required:"true"`
	From     string `json:"from"      required:"true"`
	Host     string `json:"host"      required:"true"`
	Port     int    `json:"port"      required:"true"`
}
