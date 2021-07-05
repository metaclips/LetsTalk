package paging

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"go.mongodb.org/mongo-driver/bson"
)

type Paging struct {
	Limit  int64  `json:"limit" bson:"limit" validate:"required,min=0,max=30"`
	Skip   int64  `json:"skip" bson:"skip" validate:"required,min=0"`
	Search string `json:"search" bson:"search"`
}

func (p *Paging) PagedRequest(ctx context.Context, fields bson.M, store database.Database, data interface{}) (int, error) {
	queries := make(bson.M)

	if p.Search != "" {
		queries["$text"] = bson.M{
			"$search": p.Search,
		}
	}

	for key, value := range fields {
		queries[key] = value
	}

	err := store.FindMany(ctx, queries, nil, bson.M{"_id": -1}, p.Limit, p.Skip, &data)
	if err != nil {
		return 0, err
	}

	count, err := store.Count(ctx, queries)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
