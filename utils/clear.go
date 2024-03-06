/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

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
	if len(s) < 1 {
		return
	}

	bs := *(*[]byte)(unsafe.Pointer(&s)) // #nosec G103
	for i := 0; i < len(bs); i++ {
		bs[i] = 0
	}
}
