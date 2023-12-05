package gitea

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

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
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890'-!\"#$%&()*,./:;?@[]^_`{|}~+<=>"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	// gen 32 bytes password
	for i := 0; i < 32; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
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
