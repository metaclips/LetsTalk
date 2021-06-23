package v1resources

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileChunks struct {
	FileID         primitive.ObjectID `json:"_id" bson:"_id"`
	UniqueFileHash string             `bson:"_id" json:"fileHash,omitempty"`
	ChunkIndex     int                `bson:"chunkIndex" json:"chunkIndex"`
}
