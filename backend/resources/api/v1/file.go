package v1resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileDetail struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id" validate:"required"`
	UserID       primitive.ObjectID `json:"userID" bson:"userID" validate:"required"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	Hash         string             `json:"hash" bson:"hash" validate:"required"`
	CurrentIndex int                `json:"currentIndex" bson:"currentIndex"`
	TotalIndex   int                `json:"totalIndex" bson:"totalIndex" validate:"required"`
	CreatedAt    time.Time          `bson:"createdAt" json:"-"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}
