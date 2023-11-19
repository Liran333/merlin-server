package app

import (
	"testing"

	"github.com/openmerlin/merlin-server/user/domain"
)

func TestCreateCmdValidate(t *testing.T) {
	validAcc, _ := domain.NewAccount("xxx")
	validEmail, _ := domain.NewEmail("yyy")
	tests := []UserCreateCmd{
		{},
		{
			Email:   validEmail,
			Account: validAcc,
		},
		{
			Email: validEmail,
		},
		{
			Account: validAcc,
		},
	}

	results := []bool{
		false,
		true,
		false,
		false,
	}

	for i, v := range tests {
		if ok := (v.Validate() != nil); !ok {
			t.Errorf("case num %#v valid result is %v ", v, results[i])
		}
	}
}
