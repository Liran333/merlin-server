/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package utils provides utility functions for various purposes.
package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	statusCodeUpLimit   = 200
	statusCodeDownLimit = 299
	defaultBackoff      = 10 * time.Millisecond
)

type HttpClient struct {
	client     http.Client
	maxRetries int
}

func newClient(timeout int) http.Client {
	return http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func NewHttpClient(n, timeout int) HttpClient {
	return HttpClient{
		maxRetries: n,
		client:     newClient(timeout),
	}
}

func (hc *HttpClient) SendAndHandle(req *http.Request, handle func(http.Header, io.Reader) error) error {
	resp, err := hc.do(req)
	if err != nil || resp == nil {
		logrus.Errorf("req remote url is err:%s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < statusCodeUpLimit || resp.StatusCode > statusCodeDownLimit {
		rb, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("response has status:%s and body:%q", resp.Status, rb)
	}

	if handle != nil {
		return handle(resp.Header, resp.Body)
	}

	return nil
}

func (hc *HttpClient) do(req *http.Request) (resp *http.Response, err error) {
	if resp, err = hc.client.Do(req); err == nil {
		return
	}

	maxRetries := hc.maxRetries
	backoff := defaultBackoff

	for retries := 1; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = hc.client.Do(req); err == nil {
			break
		}
	}
	return
}
