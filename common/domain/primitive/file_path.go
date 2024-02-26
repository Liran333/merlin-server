/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

// FilePath is an interface that represents a file path.
type FilePath interface {
	FilePath() string
}

// NewCodeFilePath creates a new CodeFilePath instance with the given value.
func NewCodeFilePath(v string) (FilePath, error) {
	// todo judge the length of path
	return codeFilePath(v), nil
}

type codeFilePath string

// FilePath returns the file path as a string.
func (r codeFilePath) FilePath() string {
	return string(r)
}
