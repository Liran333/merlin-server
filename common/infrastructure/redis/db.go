package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func Init(cfg *Config, remove bool) error {
	var tlsConfig *tls.Config

	if cfg.DBCert != "" {
		ca, err := os.ReadFile(cfg.DBCert)
		if err != nil {
			return err
		}

		if remove {
			defer os.Remove(cfg.DBCert)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return fmt.Errorf("faild to append certs from PEM")
		}

		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true, // #nosec G402
		}
	}

	client = redis.NewClient(&redis.Options{
		Addr:      cfg.Address,
		Password:  cfg.Password,
		DB:        0,
		TLSConfig: tlsConfig,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil

}
