/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"testing"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
)

// TestCreateCmdValidate is a test function for validating the CreateCmd struct.
func TestCreateCmdValidate(t *testing.T) {
	validAcc := primitive.CreateAccount("xxx")
	validEmail := primitive.CreateEmail("yyy")
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
