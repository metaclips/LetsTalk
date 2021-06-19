package models

import (
	"context"

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

func (u User) GetUserByID(ctx context.Context, ID primitive.ObjectID) (*v1resources.UserData, error) {
	user := v1resources.UserData{}
	if err := u.db.FindByID(ctx, ID, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u User) GetUserByEmail(ctx context.Context, email string) (*v1resources.UserData, error) {
	user := v1resources.UserData{}
	if err := u.db.FindOne(ctx, bson.M{"email": email}, nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
