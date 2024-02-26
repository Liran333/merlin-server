/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package e2e provides end-to-end testing functionality for the application.
package e2e

import (
	"bufio"
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	swagger "e2e/client"
)

const minElements = 2

var (
	Auth       context.Context
	Auth2      context.Context
	Interal    context.Context
	Api        *swagger.APIClient
	InteralApi *swagger.APIClient
)

// LoadFromYaml used for testing
func LoadFromYaml(path string, cfg *swagger.Configuration) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func newAuthCtx(token string) context.Context {
	return context.WithValue(context.Background(), swagger.ContextAPIKey, swagger.APIKey{
		Key:    token,
		Prefix: "Bearer", // Omit if not necessary.
	})
}

func newInteralCtx(token string) context.Context {
	return context.WithValue(context.Background(), swagger.ContextAPIKey, swagger.APIKey{
		Key: token,
	})
}

func getToken() []string {
	t, err := os.Open("token")
	if err != nil {
		logrus.Fatal(err)
	}
	defer t.Close()

	res := make([]string, 0)
	reader := bufio.NewScanner(t)
	for reader.Scan() {
		res = append(res, reader.Text())
	}

	if err := reader.Err(); err != nil {
		logrus.Fatal(err)
	}

	return res
}

// TestMain used for testing
func TestMain(m *testing.M) {
	api := swagger.NewConfiguration()
	if err := LoadFromYaml("./api.yaml", api); err != nil {
		logrus.Fatal(err)
	}

	internal := swagger.NewConfiguration()
	if err := LoadFromYaml("./internal.yaml", internal); err != nil {
		logrus.Fatal(err)
	}

	token := getToken()

	// Check if token slice contains at least 2 elements.
	if len(token) < minElements {
		logrus.Fatal("Insufficient tokens provided. Need at least 2 tokens.")
	}

	Api = swagger.NewAPIClient(api)
	InteralApi = swagger.NewAPIClient(internal)

	Auth = newAuthCtx(token[0])  // Use the first token.
	Auth2 = newAuthCtx(token[1]) // Use the second token.
	Interal = newInteralCtx("12345")

	m.Run()
}
