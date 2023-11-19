package session

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

var errDocNotExists = errors.New("doc doesn't exist")

type mongodbClient interface {
	IsDocExists(error) bool
	IsDocNotExists(error) bool

	ReplaceDoc(ctx context.Context, filterOfDoc, project bson.M) (string, error)
	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error
}

func isDocNotExists(err error) bool {
	return errors.Is(err, errDocNotExists)
}
