/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"golang.org/x/xerrors"
)

// Hardware is an interface that defines hardware-related operations.
type Task interface {
	Task() string
}

// NewHardware creates a new Hardware instance decided by sdk based on the given string.
func NewTask(v string) (Task, error) {
	if v == "" {
		return task(v), nil
	}

	if !tasks.Has(v) {
		return nil, xerrors.Errorf("unsupported task, %s", v)
	}

	return task(v), nil
}

// CreateHardware creates a new Hardware instance based on the given string.
func CreateTask(v string) Task {
	return task(v)
}

type task string

// BaseImage returns the base image of the base image.
func (r task) Task() string {
	return string(r)
}
