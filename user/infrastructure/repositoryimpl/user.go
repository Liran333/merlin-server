package repositoryimpl

import (
	"context"
	"errors"
	"fmt"
	"sort"

	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
	"github.com/openmerlin/merlin-server/infrastructure/repositories"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
)

func NewUserRepo(m mongodbClient) repository.User {
	return &userRepoImpl{m}
}

type userRepoImpl struct {
	cli mongodbClient
}

func (impl *userRepoImpl) Save(u *domain.User) (r domain.User, err error) {
	if u.Id != "" {
		if err = impl.update(u); err != nil {
			err = fmt.Errorf("failed to update user: %w", err)
		} else {
			r = *u
			r.Version += 1
		}

		return
	}

	v, err := impl.insert(u)
	if err != nil {
		err = fmt.Errorf("failed to add user info %w", err)
	} else {
		r = *u
		r.Id = v
	}

	return
}

func (impl *userRepoImpl) Delete(u *domain.User) (err error) {
	filter, err := mongo.ObjectIdFilter(u.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.DeleteOne(
			ctx, filter,
		)
	}

	if err = primitive.WithContext(f); err != nil {
		err = fmt.Errorf("Delete failed: %w", err)
	}

	return
}

func (impl *userRepoImpl) GetByAccount(account domain.Account) (r domain.User, err error) {
	if r, _, err = impl.GetByFollower(account, nil); err != nil {
		err = repositories.ConvertError(err)

		return
	}

	return
}

func (impl *userRepoImpl) update(u *domain.User) (err error) {
	var user DUser
	err = toUserDoc(*u, &user)
	if err != nil {
		return
	}
	doc, err := mongo.GenDoc(user)
	if err != nil {
		return
	}

	filter, err := mongo.ObjectIdFilter(u.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter, doc, mongoCmdSet, u.Version,
		)
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocNotExists(err) {
		err = fmt.Errorf("concurrent updating: %w", err)
	}

	return
}

func (impl *userRepoImpl) insert(u *domain.User) (id string, err error) {
	var user DUser
	err = toUserDoc(*u, &user)
	if err != nil {
		return
	}

	doc, err := mongo.GenDoc(user)
	if err != nil {
		return
	}

	doc[fieldVersion] = 0
	doc[fieldFollower] = bson.A{}
	doc[fieldFollowing] = bson.A{}

	f := func(ctx context.Context) error {
		v, err := impl.cli.NewDocIfNotExist(
			ctx, mongo.UserDocFilterByAccount(u.Account.Account()), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (impl *userRepoImpl) GetByFollower(owner, follower domain.Account) (
	u domain.User, isFollower bool, err error,
) {
	var v []struct {
		DUser `bson:",inline"`

		IsFollower     bool `bson:"is_follower"`
		FollowerCount  int  `bson:"follower_count"`
		FollowingCount int  `bson:"following_count"`
	}

	f := func(ctx context.Context) error {
		fields := bson.M{
			fieldFollowerCount:  bson.M{"$size": "$" + fieldFollower},
			fieldFollowingCount: bson.M{"$size": "$" + fieldFollowing},
		}

		if follower != nil {
			fields[fieldIsFollower] = bson.M{
				"$in": bson.A{follower.Account(), "$" + fieldFollower},
			}
		}

		pipeline := bson.A{
			bson.M{"$match": mongo.UserDocFilterByAccount(owner.Account())},
			bson.M{"$addFields": fields},
			bson.M{"$project": bson.M{
				fieldFollowing: 0,
				fieldFollower:  0,
			}},
		}

		cursor, err := impl.cli.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = primitive.WithContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = repositories.NewErrorDataNotExists(errors.New("no user"))

		return
	}

	item := &v[0]
	if err = toUser(item.DUser, &u); err != nil {
		return
	}

	if follower != nil {
		isFollower = item.IsFollower
	}

	return
}

func (impl *userRepoImpl) FindUsersInfo(accounts []domain.Account) (r []domain.UserInfo, err error) {
	var v []DUser

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldName: bson.M{
				"$in": accounts,
			},
		}

		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{
				fieldName:     1,
				fieldAvatarId: 1,
			}, &v,
		)
	}

	if err := primitive.WithContext(f); err != nil {
		err = fmt.Errorf("failed to find user info: %w", err)

		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	r = make([]domain.UserInfo, len(v))
	for i := range v {
		if err := toUserInfo(v[i], &r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl *userRepoImpl) GetUserAvatarId(account domain.Account) (id domain.AvatarId, err error) {

	var v DUser

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, bson.M{fieldName: account},
			bson.M{fieldAvatarId: 1}, &v,
		)
	}

	if err := primitive.WithContext(f); err != nil {
		err = fmt.Errorf("failed to get user avatar id: %w", err)

		return nil, err
	}

	if id, err = domain.NewAvatarId(v.AvatarId); err != nil {
		return
	}

	return
}

func (impl *userRepoImpl) Search(opt *repository.UserSearchOption) (
	r repository.UserSearchResult, err error,
) {
	key := "$" + fieldFollower
	fields := bson.M{fieldCount: bson.M{
		"$cond": bson.M{
			"if":   bson.M{"$isArray": key},
			"then": bson.M{"$size": key},
			"else": 0,
		},
	}}

	pipeline := bson.A{
		bson.M{mongo.MongoCmdMatch: bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$" + fieldName,
					"regex":   opt.Name,
					"options": "i",
				},
			},
		}},
		bson.M{"$addFields": fields},
		bson.M{"$project": bson.M{
			fieldName:  1,
			fieldCount: 1,
		}},
	}

	var v []userInfo

	f := func(ctx context.Context) error {
		col := impl.cli.Collection()
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = primitive.WithContext(f); err != nil || len(v) == 0 {
		return
	}

	total := len(v)

	items := make([]*userInfo, total)
	for i := range v {
		items[i] = &v[i]
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count >= items[j].Count
	})

	n := opt.TopNum
	if total < n {
		n = total
	}

	top := make([]domain.Account, n)
	for i := range top {
		if top[i], err = primitive.NewAccount(items[i].Name); err != nil {
			return
		}
	}

	r.Top = top
	r.Total = n

	return
}

func (impl *userRepoImpl) getPageItems(items []string, option *repository.FollowFindOption) []string {
	if option.CountPerPage <= 0 {
		return items
	}

	total := len(items)

	if option.PageNum <= 1 {
		if total > option.CountPerPage {
			return items[:option.CountPerPage]
		}

		return items
	}

	skip := option.CountPerPage * (option.PageNum - 1)
	if skip >= total {
		return nil
	}

	if n := total - skip; n > option.CountPerPage {
		return items[skip : skip+option.CountPerPage]
	}

	return items[skip:]
}
