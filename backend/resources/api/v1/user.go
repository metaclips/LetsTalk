package v1resources

import (
	"strings"
	"time"

	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegistration struct {
	ID           primitive.ObjectID `bson:"_id" json:"-"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	FirstName    string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName     string             `bson:"lastName" json:"lastName" validate:"required"`
	MatricNumber string             `bson:"matricNumber" json:"matricNumber"`
	Department   string             `bson:"department" json:"department" validate:"required"`
	Faculty      string             `bson:"faculty" json:"faculty"`
	Password     string             `bson:"-" json:"password" validate:"required"`
	Year         int                `bson:"year" json:"year"`
	IsStaff      bool               `bson:"isStaff" json:"isStaff"`
	DOB          ModifiedTimeFormat `bson:"dob" json:"dob" validate:"required"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	DeletedAt    *time.Time         `json:"deletedAt" bson:"deletedAt"`
}

type ModifiedTimeFormat time.Time

func (t *ModifiedTimeFormat) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), "\"\\")
	format, err := time.Parse(constants.DateOnly, value)
	if err != nil {
		return err
	}

	*t = ModifiedTimeFormat(format)
	return nil
}

type UserLogin struct {
	Email    string `bson:"_id" json:"email" validate:"required,email"`
	Password string `bson:"-" json:"password" validate:"required"`
}

type PrunedUser struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Email     string             `bson:"email" json:"email"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
}

type PrunedUserList struct {
	Data  []PrunedUser
	Count int
}

type UserStatus struct {
	UserID   primitive.ObjectID `json:"userID"`
	IsOnline bool               `json:"isOnline"`
}

type UserData struct {
	ID           primitive.ObjectID `bson:"_id" json:"-"`
	Email        string             `bson:"email" json:"email"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	Department   string             `bson:"department" json:"department" validate:"required"`
	Faculty      string             `bson:"faculty" json:"faculty"`
	UUID         string             `bson:"uuid" json:"uuid"`
	IsStaff      bool               `bson:"isStaff" json:"isStaff"`
	Verified     bool               `json:"verified" bson:"verified"`
	Year         int                `bson:"year" json:"year"`
	Password     []byte             `bson:"password" json:"-"`
	DOB          time.Time          `bson:"dob" json:"dob"`
	IsBanned     bool               `json:"isBanned" bson:"isBanned"`
	MatricNumber *string            `bson:"matricNumber" json:"matricNumber"`
	CreatedAt    time.Time          `bson:"createdAt" json:"-"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

type AssociateData struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
}
