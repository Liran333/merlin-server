package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
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

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
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
	backoff := 10 * time.Millisecond

	for retries := 1; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = hc.client.Do(req); err == nil {
			break
		}
	}
	return
}
