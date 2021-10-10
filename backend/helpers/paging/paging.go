package paging

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Paging struct {
	Limit  int    `json:"limit" bson:"limit" validate:"required,min=0,max=30"`
	Skip   int    `json:"skip" bson:"skip" validate:"required,min=0"`
	Search string `json:"search" bson:"search"`
}

func (p *Paging) PagedRequest(ctx context.Context, fields bson.M, store database.Database, data interface{}) (int, error) {

	dd := []v1resources.PrunedUser{}

	err := store.FindMany(ctx, bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "/150/"}}}, nil, nil, int64(p.Limit), int64(p.Skip), &dd)
	if err != nil {
		return 0, err
	}

	count, err := store.Count(ctx, bson.M{"email": bson.M{"$regex": primitive.Regex{Pattern: "/150/"}}})
	if err != nil {
		return 0, err
	}

	return count, nil
}
