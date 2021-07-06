package v1resources

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoginResponse struct {
	UserID   primitive.ObjectID
	Encoded  string
	WSTicket string
}
