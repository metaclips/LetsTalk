package v1resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileChunks struct {
	FileID         primitive.ObjectID `json:"_id" bson:"_id"`
	UniqueFileHash string             `bson:"_id" json:"fileHash,omitempty"`
	ChunkIndex     int                `bson:"chunkIndex" json:"chunkIndex"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	DeletedAt      *time.Time         `bson:"deletedAt" json:"-"`
}
