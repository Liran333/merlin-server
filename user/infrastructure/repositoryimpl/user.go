package repositoryimpl

import (
	"context"
	"errors"
	"fmt"

	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	mongo "github.com/openmerlin/merlin-server/common/infrastructure/mongo"
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
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}
	}
	return
}

func (impl *userRepoImpl) GetByAccount(account domain.Account) (r domain.User, err error) {
	if r, _, err = impl.GetByFollower(account, nil); err != nil {
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
			ctx, mongo.DocFilterByAccount(u.Account.Account()), doc,
		)

		id = v

		return err
	}

	if err = primitive.WithContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
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
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return
	}

	if len(v) == 0 {
		err = commonrepo.NewErrorResourceNotExists(errors.New("no user"))

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

func (impl *userRepoImpl) GetUserFullname(account domain.Account) (fullname string, err error) {

	var v DUser

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, mongo.DocFilterByAccount(account.Account()),
			bson.M{fieldFullname: 1}, &v,
		)
	}

	if err := primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return "", err
	}

	return v.Fullname, nil
}

func (impl *userRepoImpl) GetUsersAvatarId(names []string) (users []domain.User, err error) {
	var v []DUser

	if len(names) == 0 {
		err = commonrepo.NewErrorResourceNotExists(err)

		return
	}

	filter := bson.M{}
	filter[fieldName] = bson.M{
		"$in": names,
	}

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, filter,
			bson.M{fieldAvatarId: 1, fieldName: 1}, &v,
		)
	}

	if err := primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	users = make([]domain.User, len(v))
	for i := range v {
		users[i] = domain.User{
			Account:  primitive.CreateAccount(v[i].Name),
			AvatarId: domain.CreateAvatarId(v[i].AvatarId),
		}
	}

	return
}

func (impl *userRepoImpl) GetUserAvatarId(account domain.Account) (id domain.AvatarId, err error) {

	var v DUser

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, mongo.DocFilterByAccount(account.Account()),
			bson.M{fieldAvatarId: 1}, &v,
		)
	}

	if err := primitive.WithContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = commonrepo.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	if id, err = domain.NewAvatarId(v.AvatarId); err != nil {
		return
	}

	return
}
