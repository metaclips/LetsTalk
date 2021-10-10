package models

import (
	"context"
	"time"

	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	db database.Database
}

func NewUser(db database.Database) User {
	return User{
		db,
	}
}

func (u User) CreateUser(ctx context.Context, user *v1resources.UserData) error {
	return u.db.Save(ctx, user, nil)
}

func (u User) UpdateUser(ctx context.Context, user *v1resources.UserData) error {
	return u.db.UpdateByID(ctx, user.ID, user)
}

func (u User) GetManyAssociates(ctx context.Context, associates []primitive.ObjectID) ([]v1resources.AssociateData, error) {
	data := make([]bson.M, len(associates))
	for i, associate := range associates {
		data[i] = bson.M{"_id": associate}
	}

	if len(data) == 0 || data == nil {
		return []v1resources.AssociateData{}, nil
	}

	val := []v1resources.AssociateData{}

	err := u.db.FindMany(ctx, bson.M{"$or": data}, nil, nil, 0, 0, &val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (u User) GetUserByID(ctx context.Context, ID primitive.ObjectID) (*v1resources.UserData, error) {
	user := v1resources.UserData{}
	if err := u.db.FindByID(ctx, ID, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) GetUserByUUID(ctx context.Context, uuid string) (*v1resources.UserData, error) {
	user := v1resources.UserData{}
	if err := u.db.FindOne(ctx, bson.M{"uuid": uuid}, nil, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) GetRegisteredMinorsPruned(ctx context.Context) (*v1resources.PrunedUserList, error) {
	user := []v1resources.PrunedUser{}
	if err := u.db.FindOne(ctx, bson.M{"dob": bson.M{"lt": time.Now().Add(-(time.Hour * 365 * 24) - time.Hour)}}, nil, nil, &user); err != nil {
		return nil, err
	}

	count, err := u.db.Count(ctx, bson.M{"dob": bson.M{"$lt": time.Now().Add(-(time.Hour * 365 * 24) - time.Hour)}})
	if err != nil {
		return nil, err
	}

	return &v1resources.PrunedUserList{
		Data:  user,
		Count: count,
	}, nil
}

func (u User) GetRegisteredMinors(ctx context.Context, paging *paging.Paging) (*v1resources.UserDataList, error) {
	users := []v1resources.UserData{}

	fields := bson.M{
		"dob": bson.M{
			"$lt": time.Now().Add(-(time.Hour * 365 * 18) - time.Hour),
		},
	}

	count, err := paging.PagedRequest(ctx, fields, u.db, &users)
	if err != nil {
		return nil, err
	}

	return &v1resources.UserDataList{
		Data:  users,
		Count: count,
	}, nil
}

func (u User) GetUsersPrunedData(ctx context.Context, paging *paging.Paging) (*v1resources.PrunedUserList, error) {
	users := make([]v1resources.PrunedUser, 0)

	err := u.db.FindMany(ctx, bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: paging.Search}}}, nil, nil, int64(paging.Limit), int64(paging.Skip), &users)
	if err != nil {
		return nil, err
	}

	return &v1resources.PrunedUserList{
		Data:  users,
		Count: len(users),
	}, nil
}

func (u User) GetUsers(ctx context.Context, paging *paging.Paging, fields bson.M) ([]v1resources.UserData, error) {
	bannedUsers := []v1resources.UserData{}
	_, err := paging.PagedRequest(ctx, fields, u.db, &bannedUsers)
	if err != nil {
		return nil, err
	}
	return bannedUsers, nil
}

func (u User) GetUserByEmail(ctx context.Context, email string) (*v1resources.UserData, error) {
	user := v1resources.UserData{}
	if err := u.db.FindOne(ctx, bson.M{"email": email}, nil, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
