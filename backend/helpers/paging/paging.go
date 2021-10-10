package paging

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"go.mongodb.org/mongo-driver/bson"
)

type Paging struct {
	Limit  int    `json:"limit" bson:"limit" validate:"required,min=0,max=30"`
	Skip   int    `json:"skip" bson:"skip" validate:"required,min=0"`
	Search string `json:"search" bson:"search"`
}

func (p *Paging) PagedRequest(ctx context.Context, fields bson.M, store database.Database, data interface{}) (int, error) {
	err := store.FindMany(ctx, fields, nil, bson.M{"_id": -1}, int64(p.Limit), int64(p.Skip), data)
	if err != nil {
		return 0, err
	}

	count, err := store.Count(ctx, fields)
	if err != nil {
		return 0, err
	}

	return count, nil
}
