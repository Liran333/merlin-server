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
	fieldInvitee        = "user_name"
	fieldInviter        = "inviter"
	fieldStatus         = "status"
)

// TODO: senstive data should be store in vault
type Organization struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Name              string `bson:"name"       json:"name"`
	AvatarId          string `bson:"avatar_id"  json:"avatar_id"`
	Website           string `bson:"website"    json:"website"`
	PlatformId        string `bson:"platform_id"        json:"platform_id"`
	Description       string `bson:"description"        json:"description"`
	FullName          string `bson:"fullname"        json:"fullname"`
	Owner             string `bson:"owner"        json:"owner"`
	DefaultRole       string `bson:"default_role"        json:"default_role"`
	Type              int    `bson:"type"        json:"type"`
	CreatedAt         int64  `bson:"created_at"        json:"created_at"`
	AllowRequest      bool   `bson:"allow_request"        json:"allow_request"`
	OwnerTeamId       int64  `bson:"owner_team_id"        json:"owner_team_id"`
	ReadTeamId        int64  `bson:"read_team_id"        json:"read_team_id"`
	WriteTeamId       int64  `bson:"write_team_id"        json:"write_team_id"`
	ContributorTeamId int64  `bson:"contributor_team_id"        json:"contributor_team_id"`
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
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Username  string `bson:"user_name"       json:"user_name"`
	Orgname   string `bson:"org_name"       json:"org_name"`
	Role      string `bson:"role"        json:"role"`
	Expire    int64  `bson:"expire"        json:"expire"`
	Inviter   string `bson:"inviter"        json:"inviter"`
	Status    string `bson:"status"        json:"status"`
	By        string `bson:"by"        json:"by"`
	Msg       string `bson:"msg"        json:"msg"`
	CreatedAt int64  `bson:"created_at"        json:"created_at"`
	UpdatedAt int64  `bson:"updated_at"        json:"updated_at"`
	Version   int    `bson:"version"    json:"-"`
}

type MemberRequest struct {
	Id primitive.ObjectID `bson:"_id"       json:"-"`

	Username  string `bson:"user_name"       json:"user_name"`
	Orgname   string `bson:"org_name"       json:"org_name"`
	Role      string `bson:"role"        json:"role"`
	Status    string `bson:"status"        json:"status"`
	By        string `bson:"by"        json:"by"`
	CreatedAt int64  `bson:"created_at"        json:"created_at"`
	UpdatedAt int64  `bson:"updated_at"        json:"updated_at"`
	Msg       string `bson:"msg"        json:"msg"`
	// Version will be increased by 1 automatically.
	// So, don't marshal it to avoid setting it occasionally.
	Version int `bson:"version"    json:"-"`
}
