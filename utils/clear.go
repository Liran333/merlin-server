/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package utils provides utility functions for various purposes.
package utils

import "unsafe"

// ClearByteArrayMemory clears the memory of byte arrays by setting each byte to 0.
func ClearByteArrayMemory(bs ...[]byte) {
	f := func(b []byte) {
		for i := range b {
			b[i] = 0
		}
	}

	for i := range bs {
		f(bs[i])
	}
}

// ClearStringMemory clears the memory of strings by calling clearStringMemory function.
func ClearStringMemory(s ...string) {
	for i := range s {
		clearStringMemory(s[i])
	}
}

// clearStringMemory clears the memory of a string by setting each byte to 0.
func clearStringMemory(s string) {
	// for strings represented by a single character, Go's runtime implements string sharing of character data,
	// which resides in a unified staticbytes. Therefore, it is forbidden to clear it,
	// as doing so would result in all strings parsed from a single character having a value of '\0'
	// during later runtime execution of the program. This is to avoid setting sensitive information such as passwords
	// and credentials with a length less than or equal to 1.
	if len(s) <= 1 {
		return
	}

	bs := *(*[]byte)(unsafe.Pointer(&s)) // #nosec G103
	for i := 0; i < len(bs); i++ {
		bs[i] = 0
	}
}
