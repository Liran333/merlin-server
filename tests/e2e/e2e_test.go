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

var (
	Auth  context.Context
	Auth2 context.Context
	Api   *swagger.APIClient
)

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

func getToken() []string {
	t, err := os.Open("token")
	if err != nil {
		logrus.Fatal(err)
	}

	res := make([]string, 0)

	reader := bufio.NewScanner(t)
	for reader.Scan() {
		res = append(res, reader.Text())
	}

	return res
}

func TestMain(m *testing.M) {
	cfg := swagger.NewConfiguration()
	if err := LoadFromYaml("./cfg.yaml", cfg); err != nil {
		logrus.Fatal(err)
	}

	token := getToken()

	Api = swagger.NewAPIClient(cfg)

	Auth = newAuthCtx(token[0])
	Auth2 = newAuthCtx(token[1])

	m.Run()
}
