package primitive

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/sirupsen/logrus"
)

type Password interface {
	Password() string
}

func NewPassword() (Password, error) {
	str, err := genRandomString(passwordLength)

	return password(str), err
}

func CreatePassword(v string) Password {
	return password(v)
}

type password string

func (r password) Password() string {
	return string(r)
}

func genRandomString(len int) (string, error) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!#$%&()*,./:;?@[]^_`{|}~+<=>"

	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, err := rand.Int(rand.Reader, bigInt)
		if err != nil {
			logrus.Errorf("internal error, rand.Int: %s", err.Error())

			return "", err
		}

		container += string(str[randomInt.Int64()])
	}

	return container, nil
}
