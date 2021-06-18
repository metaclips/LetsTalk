package v1resources

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegistration struct {
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	DOB       string `bson:"age" json:"age"`
	Class     string `bson:"class" json:"class"`
	Password  string `bson:"-" json:"password"`
}

type UserLogin struct {
	Email    string `bson:"_id" json:"userID"`
	Password string `bson:"-" json:"password"`
}

type UserData struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Email     string             `bson:"email" json:"email"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	DOB       string             `bson:"age" json:"age"`
	Class     string             `bson:"class" json:"class"`
	UUID      string             `bson:"uuid" json:"uuid"`
	Password  []byte             `bson:"password" json:"-"`
}
