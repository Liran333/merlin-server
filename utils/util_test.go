/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package utils

import "testing"

// TestAnonymizeEmail test email anonymize
func TestAnonymizeEmail(t *testing.T) {
	tests := []struct {
		email          string
		expectedResult string
	}{
		{"example@example.com", "exam***@example.com"},
		{"test@test.com", "t***@test.com"},
		{"user@domain.com", "u***@domain.com"},
		{"abcd@efg.com", "a***@efg.com"},
		{"abc@def.com", "***@def.com"},
		{"co@outlook.com", "***@outlook.com"},
		{"2@gmail.com", "***@gmail.com"},
		{"@gmail.com", "***@gmail.com"},
		{"abc123", "abc123"},
		{"", ""},
	}

	for _, test := range tests {
		result := AnonymizeEmail(test.email)
		if result != test.expectedResult {
			t.Errorf("AnonymizeEmail(%s) returned %s, expected %s", test.email, result, test.expectedResult)
		}
	}
}
