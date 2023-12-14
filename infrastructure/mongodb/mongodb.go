package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var cli *client

func Initialize(conn, dbName, dbCert string, remove bool) error {
	uri := conn
	clientOpt := options.Client().ApplyURI(uri)

	if dbCert != "" {
		ca, err := os.ReadFile(dbCert) // #nosec G304
		if err != nil {
			return err
		}

		if remove {
			defer os.Remove(dbCert)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return fmt.Errorf("faild to append certs from PEM")
		}

		tlsConfig := &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true, // #nosec G402
		}

		clientOpt.SetTLSConfig(tlsConfig)
	}

	c, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		return err
	}

	cli = &client{
		c:  c,
		db: c.Database(dbName),
	}

	return nil
}

func Close() error {
	if cli != nil {
		return cli.disconnect()
	}

	return nil
}

type client struct {
	c  *mongo.Client
	db *mongo.Database
}

func (cli *client) disconnect() error {
	return primitive.WithContext(cli.c.Disconnect)
}

func (cli *client) collection(name string) *mongo.Collection {
	return cli.db.Collection(name)
}
