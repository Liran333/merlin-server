package giteauser

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"code.gitea.io/sdk/gitea"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

type UserCreateCmd struct {
	Username string `json:"username"`
	Email    string `json:"email"`
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

func toUser(u *gitea.User) (user domain.User, err error) {
	user.Account, err = primitive.NewAccount(u.UserName)
	if err != nil {
		return
	}

	user.Email, err = domain.NewEmail(u.Email)
	if err != nil {
		return
	}

	user.PlatformId = u.ID

	return
}

func GetClient(c *gitea.Client) *UserClient {
	return &UserClient{
		client: c,
	}
}

// Client admin client
type UserClient struct {
	client *gitea.Client
}

func (c *UserClient) CreateUser(cmd *UserCreateCmd) (user domain.User, err error) {
	changePwd := false
	pwd := genPasswd()
	o := gitea.CreateUserOption{
		Username:           cmd.Username,
		Email:              cmd.Email,
		Password:           pwd,
		MustChangePassword: &changePwd,
	}

	u, _, err := c.client.AdminCreateUser(o)
	if err != nil {
		return
	}

	user, err = toUser(u)
	user.PlatformPwd = pwd

	return
}

func (c *UserClient) DeleteUser(name string) (err error) {
	_, err = c.client.AdminDeleteUser(name)

	return
}
