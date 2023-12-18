package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoCmdSet  = "$set"
	mongoCmdPush = "$push"
)

type mongodbClient interface {
	IsDocExists(error) bool
	IsDocNotExists(error) bool

	Collection() *mongo.Collection
	ObjectIdFilter(s string) (bson.M, error)
	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)
	UpdateDoc(ctx context.Context, filterOfDoc, update bson.M, op string, version int) error
	DeleteOne(ctx context.Context, filterOfDoc bson.M) error
	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error
	GetDocs(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error
	AddToSimpleArray(ctx context.Context, array string, filterOfDoc, value interface{}) error
	RemoveFromSimpleArray(ctx context.Context, array string, filterOfDoc, value interface{}) error
}
