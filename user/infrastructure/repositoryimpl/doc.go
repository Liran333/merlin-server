package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	fieldName           = "name"
	fieldCount          = "count"
	fieldPlatformId     = "platform_id"
	fieldAccount        = "account"
	fieldBio            = "bio"
	fieldAvatarId       = "avatar_id"
	fieldVersion        = "version"
	fieldFollower       = "follower"
	fieldFollowing      = "following"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldIsFollower     = "is_follower"
)

// TODO: senstive data should be store in vault
type DUser struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name                    string   `bson:"name"       json:"name"`
	Email                   string   `bson:"email"      json:"email"`
	Bio                     string   `bson:"bio"        json:"bio"`
	AvatarId                string   `bson:"avatar_id"  json:"avatar_id"`
	PlatformTokens          []DToken `bson:"tokens"      json:"tokens"`
	PlatformUserId          string   `bson:"uid"        json:"uid"`
	PlatformUserNamespaceId string   `bson:"nid"        json:"nid"`
	PlatformId              int64    `bson:"platform_id"        json:"platform_id"`
	PlatformPwd             string   `bson:"platform_pwd"        json:"platform_pwd"`

	Follower  []string `bson:"follower"   json:"-"`
	Following []string `bson:"following"  json:"-"`

	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type DToken struct {
	Token      string `bson:"-"   json:"-"`
	Name       string `bson:"name"   json:"name"`
	Account    string `bson:"account"   json:"account"`
	Expire     int64  `bson:"expire"   json:"expire"` // timeout in seconds
	CreatedAt  int64  `bson:"created_at"   json:"created_at"`
	Permission string `bson:"permission"   json:"permission"`
}

type DUserRegInfo struct {
	Account  string            `bson:"account"        json:"account"`
	Name     string            `bson:"name"           json:"name"`
	City     string            `bson:"city"           json:"city"`
	Email    string            `bson:"email"          json:"email"`
	Phone    string            `bson:"phone"          json:"phone"`
	Identity string            `bson:"identity"       json:"identity"`
	Province string            `bson:"province"       json:"province"`
	Detail   map[string]string `bson:"detail"         json:"detail"`
	Version  int               `bson:"version"        json:"-"`
}
