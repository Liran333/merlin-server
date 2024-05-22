/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"errors"
	"strings"
)

// SDK is an interface that defines the method to get the SDK name.
type SDK interface {
	SDK() string
}

// NewSDK creates a new SDK instance based on the given version string.
func NewSDK(v string) (SDK, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if _, ok := sdkObjects[v]; v == "" || !ok {
		return nil, errors.New("unsupported sdk")
	}

	return sdk(v), nil
}

// CreateSDK creates a new SDK instance based on the given version string.
func CreateSDK(v string) SDK {
	return sdk(v)
}

type sdk string

// SDK returns the string representation of the sdk.
func (r sdk) SDK() string {
	return string(r)
}

const (
	static = "static"
	gradio = "gradio"
)

var (
	// StaticSdk represents static sdk.
	StaticSdk = sdk(static)
	// GradioSdk represents gradio sdk.
	GradioSdk = sdk(gradio)
)