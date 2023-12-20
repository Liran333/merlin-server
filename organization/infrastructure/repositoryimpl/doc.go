package repositoryimpl

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	fieldName           = "name"
	fieldOwner          = "owner"
	fieldCount          = "count"
	fieldAccount        = "account"
	fieldBio            = "bio"
	fieldAvatarId       = "avatar_id"
	fieldVersion        = "version"
	fieldFollower       = "follower"
	fieldFollowing      = "following"
	fieldFollowerCount  = "follower_count"
	fieldFollowingCount = "following_count"
	fieldIsFollower     = "is_follower"
	fieldUser           = "user_name"
	fieldOrg            = "org_name"
	fieldRole           = "role"
)

// TODO: senstive data should be store in vault
type Organization struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name        string    `bson:"name"       json:"name"`
	AvatarId    string    `bson:"avatar_id"  json:"avatar_id"`
	Website     string    `bson:"website"    json:"website"`
	PlatformId  string    `bson:"platform_id"        json:"platform_id"`
	Description string    `bson:"description"        json:"description"`
	FullName    string    `bson:"full_name"        json:"full_name"`
	Owner       string    `bson:"owner"        json:"owner"`
	Approves    []Approve `bson:"approves"        json:"approves"`

	OwnerTeamId int64 `bson:"owner_team_id"        json:"owner_team_id"`
	ReadTeamId  int64 `bson:"read_team_id"        json:"read_team_id"`
	WriteTeamId int64 `bson:"write_team_id"        json:"write_team_id"`
	AdminTeamId int64 `bson:"admin_team_id"        json:"admin_team_id"`
	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type Member struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Username string `bson:"user_name"       json:"user_name"`
	Orgname  string `bson:"org_name"       json:"org_name"`
	Role     string `bson:"role"        json:"role"`
	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}

type Approve struct {
	Username string `bson:"user_name"       json:"user_name"`
	Orgname  string `bson:"org_name"       json:"org_name"`
	Role     string `bson:"role"        json:"role"`
	Expire   int64  `bson:"expire"        json:"expire"`
}
