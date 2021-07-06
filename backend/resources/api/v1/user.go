package v1resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegistration struct {
	Email        string    `bson:"email" json:"email" validate:"required,email"`
	FirstName    string    `bson:"firstName" json:"firstName" validate:"required"`
	LastName     string    `bson:"lastName" json:"lastName" validate:"required"`
	MatricNumber string    `bson:"matricNumber" json:"matricNumber"`
	Department   string    `bson:"department" json:"department" validate:"required"`
	Faculty      string    `bson:"faculty" json:"faculty"`
	Password     string    `bson:"-" json:"password" validate:"required"`
	Year         int       `bson:"year" json:"year"`
	IsStaff      bool      `bson:"isStaff" json:"isStaff"`
	DOB          time.Time `bson:"dob" json:"dob" validate:"required"`
}

type UserLogin struct {
	Email    string `bson:"_id" json:"email" validate:"required,email"`
	Password string `bson:"-" json:"password" validate:"required"`
}

type PrunedUser struct {
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}

type PrunedUserList struct {
	Data  []PrunedUser
	Count int
}

type UserStatus struct {
	UserID   primitive.ObjectID `json:"email"`
	IsOnline bool               `json:"isOnline"`
}

type UserData struct {
	ID           primitive.ObjectID `bson:"_id" json:"-"`
	Email        string             `bson:"email" json:"email"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	MatricNumber string             `bson:"matricNumber" json:"matricNumber"`
	Department   string             `bson:"department" json:"department" validate:"required"`
	Faculty      string             `bson:"faculty" json:"faculty"`
	UUID         string             `bson:"uuid" json:"uuid"`
	IsStaff      bool               `bson:"isStaff" json:"isStaff"`
	Year         int                `bson:"year" json:"year"`
	Password     []byte             `bson:"password" json:"-"`
	DOB          time.Time          `bson:"dob" json:"dob"`
	CreatedAt    time.Time          `bson:"createdAt" json:"-"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}
