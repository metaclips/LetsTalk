package v1resources

import "go.mongodb.org/mongo-driver/bson/primitive"

type LoginResponse struct {
	UserID   primitive.ObjectID `json:"userID"`
	Token    string             `json:"token"`
	WSTicket string             `json:"wsTicket"`
}
