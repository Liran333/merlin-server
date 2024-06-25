/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package postgresql provides functionality for interacting with PostgreSQL databases.
package obs

import (
	"io"
)

type ObsService interface {
	CreateObject(f io.Reader, bucket, path string) error
	GenFileDownloadURL(bucket, p string, downloadExpiry int) (string, error)
	CopyObject(bucket, dst, src string) error
	DeleteObject(bucket string, path string) error
}
