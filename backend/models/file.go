package models

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
)

type File struct {
	db database.Database
}

func NewFile(db database.Database) File {
	return File{
		db,
	}
}

// Create new Room instance. Note should be called right after storing
// success message.
func (b File) CreateNewFile(ctx context.Context, file *v1resources.FileDetail) error {
	if err := b.db.Save(ctx, file, nil); err != nil {
		return err
	}

	return nil
}

func (b File) IncFileIndex(ctx context.Context, fileHash string, count int) error {
	return b.db.IncCount(ctx, bson.M{"hash": fileHash}, count)
}
