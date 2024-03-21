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

	swaggerInternal "e2e/client_internal"
	swaggerRest "e2e/client_rest"
	swaggerWeb "e2e/client_web"
)

const minElements = 2

var (
	AuthRest   context.Context
	AuthRest2  context.Context
	Interal    context.Context
	ApiInteral *swaggerInternal.APIClient
	ApiRest    *swaggerRest.APIClient
	ApiWeb     *swaggerWeb.APIClient
	ComConfig  ComConfiguration
)

type ComConfiguration struct {
	ACCOUNT_NAME_REGEXP      string `yaml:"ACCOUNT_NAME_REGEXP"`
	ACCOUNT_NAME_MIN_LEN     int    `yaml:"ACCOUNT_NAME_MIN_LEN"`
	ACCOUNT_NAME_MAX_LEN     int    `yaml:"ACCOUNT_NAME_MAX_LEN"`
	ACCOUNT_DESC_MAX_LEN     int    `yaml:"ACCOUNT_DESC_MAX_LEN"`
	ACCOUNT_FULLNAME_MAX_LEN int    `yaml:"ACCOUNT_FULLNAME_MAX_LEN"`
	ORG_FULLNAME_MIN_LEN     int    `yaml:"ORG_FULLNAME_MIN_LEN"`
	MSD_NAME_REGEXP          string `yaml:"MSD_NAME_REGEXP"`
	MSD_NAME_MIN_LEN         int    `yaml:"MSD_NAME_MIN_LEN"`
	MSD_NAME_MAX_LEN         int    `yaml:"MSD_NAME_MAX_LEN"`
	MSD_DESC_MAX_LEN         int    `yaml:"MSD_DESC_MAX_LEN"`
	MSD_FULLNAME_MAX_LEN     int    `yaml:"MSD_FULLNAME_MAX_LEN"`
	EMAIL_REGEXP             string `yaml:"EMAIL_REGEXP"`
	EMAIL_MAX_LEN            int    `yaml:"EMAIL_MAX_LEN"`
	PHONE_REGEXP             string `yaml:"PHONE_REGEXP"`
	PHONE_MAX_LEN            int    `yaml:"PHONE_MAX_LEN"`
	WEBSITE_REGEXP           string `yaml:"WEBSITE_REGEXP"`
	WEBSITE_MAX_LEN          int    `yaml:"WEBSITE_MAX_LEN"`
	TOKEN_NAME_REGEXP        string `yaml:"TOKEN_NAME_REGEXP"`
	TOKEN_NAME_MIN_LEN       int    `yaml:"TOKEN_NAME_MIN_LEN"`
	TOKEN_NAME_MAX_LEN       int    `yaml:"TOKEN_NAME_MAX_LEN"`
	BRANCH_REGEXP            string `yaml:"BRANCH_REGEXP"`
	BRANCH_NAME_MIN_LEN      int    `yaml:"BRANCH_NAME_MIN_LEN"`
	BRANCH_NAME_MAX_LEN      int    `yaml:"BRANCH_NAME_MAX_LEN"`
}

// LoadFromYaml used for testing
func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func newAuthRestCtx(token string) context.Context {
	return context.WithValue(context.Background(), swaggerRest.ContextAPIKey, swaggerRest.APIKey{
		Key:    token,
		Prefix: "Bearer", // Omit if not necessary.
	})
}

func newInteralCtx(token string) context.Context {
	return context.WithValue(context.Background(), swaggerInternal.ContextAPIKey, swaggerInternal.APIKey{
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
	apiRest := swaggerRest.NewConfiguration()
	if err := LoadFromYaml("./rest.yaml", apiRest); err != nil {
		logrus.Fatal(err)
	}
	apiWeb := swaggerWeb.NewConfiguration()
	if err := LoadFromYaml("./web.yaml", apiWeb); err != nil {
		logrus.Fatal(err)
	}
	apiInteral := swaggerInternal.NewConfiguration()
	if err := LoadFromYaml("./internal.yaml", apiInteral); err != nil {
		logrus.Fatal(err)
	}

	token := getToken()

	// Check if token slice contains at least 2 elements.
	if len(token) < minElements {
		logrus.Fatal("Insufficient tokens provided. Need at least 2 tokens.")
	}

	ApiRest = swaggerRest.NewAPIClient(apiRest)
	ApiWeb = swaggerWeb.NewAPIClient(apiWeb)
	ApiInteral = swaggerInternal.NewAPIClient(apiInteral)

	AuthRest = newAuthRestCtx(token[0])  // Use the first token.
	AuthRest2 = newAuthRestCtx(token[1]) // Use the second token.
	Interal = newInteralCtx("12345")

	// Load specification config from yaml
	if err := LoadFromYaml("../../common.yaml", &ComConfig); err != nil {
		logrus.Fatal(err)
	}

	m.Run()
}
