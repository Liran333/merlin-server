package giteauser

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/openmerlin/go-sdk/gitea"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/sirupsen/logrus"
)

type UserCreateCmd struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserUpdateCmd = domain.UserCreateCmd

func genPasswd() (string, error) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!#$%&()*,./:;?@[]^_`{|}~+<=>"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	// gen 32 bytes password
	for i := 0; i < 32; i++ {
		randomInt, err := rand.Int(rand.Reader, bigInt)
		if err != nil {
			logrus.Errorf("internal error, rand.Int: %s", err.Error())

			return "", err
		}

		container += string(str[randomInt.Int64()])
	}
	return container, nil
}

func toUser(u *gitea.User) (user domain.User, err error) {
	user.Account, err = primitive.NewAccount(u.UserName)
	if err != nil {
		return
	}

	user.Email, err = primitive.NewEmail(u.Email)
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
	pwd, err := genPasswd()
	if err != nil {
		return
	}

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

func (c *UserClient) UpdateUser(cmd *UserUpdateCmd) (err error) {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	if cmd.Account == nil {
		return fmt.Errorf("account is nil")
	}

	d := gitea.EditUserOption{
		LoginName: cmd.Account.Account(),
	}

	if cmd.Email != nil {
		email := cmd.Email.Email()
		d.Email = &email
	}

	if cmd.Desc != nil {
		desc := cmd.Desc.MSDDesc()
		d.Description = &desc
	}

	if cmd.Fullname != nil {
		fullname := cmd.Fullname.MSDFullname()
		d.FullName = &fullname
	}

	_, err = c.client.AdminEditUser(cmd.Account.Account(), d)

	return
}
