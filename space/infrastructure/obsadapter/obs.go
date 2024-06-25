/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package postgresql provides functionality for interacting with PostgreSQL databases.
package obsadapter

import (
	"io"
	"net/http"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ObsService interface {
	CreateObject(f io.Reader, bucket, path string) error
	GenFileDownloadURL(bucket, p string, downloadExpiry int) (string, error)
	CopyObject(bucket, dst, src string) error
	DeleteObject(bucket string, path string) error
}

type obsServiceImpl struct {
	cli *obs.ObsClient
}

func NewClient(cli *obs.ObsClient) *obsServiceImpl {
	return &obsServiceImpl{
		cli: cli,
	}
}

func (s *obsServiceImpl) CreateObject(f io.Reader, bucket, path string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = bucket
	input.Key = path
	input.Body = f

	_, err := s.cli.PutObject(input)

	return err
}

func (s *obsServiceImpl) GenFileDownloadURL(bucket, p string, downloadExpiry int) (string, error) {
	input := &obs.CreateSignedUrlInput{}
	input.Method = obs.HttpMethodGet
	input.Bucket = bucket
	input.Key = p
	input.Expires = downloadExpiry

	output, err := s.cli.CreateSignedUrl(input)
	if err != nil {
		return "", err
	}

	return output.SignedUrl, nil
}

func (s *obsServiceImpl) CopyObject(bucket, dst, src string) error {
	input := &obs.CopyObjectInput{}
	input.Bucket = bucket
	input.Key = dst
	input.CopySourceBucket = bucket
	input.CopySourceKey = src
	_, err := s.cli.CopyObject(input)

	return err
}

func (s *obsServiceImpl) DeleteObject(bucket string, path string) error {
	input := &obs.DeleteObjectInput{}
	input.Bucket = bucket
	input.Key = path

	v, err := s.cli.DeleteObject(input)
	if err != nil && v != nil && v.StatusCode == http.StatusNotFound {
		err = nil
	}

	return err
}
