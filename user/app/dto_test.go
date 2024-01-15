package app

import (
	"testing"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

func TestCreateCmdValidate(t *testing.T) {
	validAcc := primitive.CreateAccount("xxx")
	validEmail, _ := primitive.NewEmail("yyy")
	tests := []domain.UserCreateCmd{
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
