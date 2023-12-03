package gitea

import (
	"fmt"
	"math/rand"
	"time"

	"code.gitea.io/sdk/gitea"
)

var _client Client

type Config struct {
	Url   string `json:"url"        required:"true"`
	Token string `json:"token"      required:"true"`
}

type Client struct {
	client *gitea.Client
}

func Init(cfg *Config) (err error) {
	if cfg == nil {
		return fmt.Errorf("cfg is nil")
	}
	_client.client, err = gitea.NewClient(cfg.Url, gitea.SetToken(cfg.Token))

	return
}

type UserCreateCmd struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetClient() *Client {
	return &_client
}

func genPasswd() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

func (c *Client) CreateUser(cmd *UserCreateCmd) (err error) {
	o := gitea.CreateUserOption{
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: genPasswd(),
	}
	_, _, err = c.client.AdminCreateUser(o)
	return
}
