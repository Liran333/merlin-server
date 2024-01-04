package domain

import (
	"testing"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

const (
	token    string = "3c1c575c980bf6daa1bd3eae3a5fbe6edc3aa0c4"
	salt     string = "VOGILiF8vCIODil7eR5r+31/YD3powAiatTo2yPbcVc"
	encToken string = "WUdbCLY50CZvp6hNznE/06XWJ3kfyn7Vz2pJW3fjVRo"
)

func TestPermCheck(t *testing.T) {
	pt := PlatformToken{
		Account:    primitive.CreateAccount("test"),
		Permission: primitive.NewReadPerm(),
		Expire:     0,
		Salt:       salt,
		Token:      encToken,
	}

	type test struct {
		t      string
		expect primitive.TokenPerm
		has    primitive.TokenPerm
		expire int64
		result string
	}

	tests := []test{
		{
			token,
			primitive.NewReadPerm(),
			primitive.NewWritePerm(),
			0,
			"",
		},
		{
			token,
			primitive.NewReadPerm(),
			primitive.NewReadPerm(),
			0,
			"",
		},
		{
			token,
			primitive.NewWritePerm(),
			primitive.NewWritePerm(),
			0,
			"",
		},
		{
			token,
			primitive.NewWritePerm(),
			primitive.NewReadPerm(),
			0,
			tokenPermDenied,
		},
		{
			"1234",
			primitive.NewReadPerm(),
			primitive.NewReadPerm(),
			0,
			tokenInvalid,
		},
		{
			token,
			primitive.NewReadPerm(),
			primitive.NewReadPerm(),
			1234,
			tokenExpired,
		},
		{
			token,
			primitive.NewReadPerm(),
			primitive.NewReadPerm(),
			9704334556,
			"",
		},
	}

	for _, tc := range tests {
		pt.Expire = tc.expire
		pt.Permission = tc.has
		if result := pt.Check(tc.t, tc.expect); (result == nil && tc.result != "") || (result != nil && result.Error() != tc.result) {
			t.Errorf("user  perm check failed, expect %s, result is %s", tc.result, result)
		}
	}
}
