package users

import "github.com/metaclips/LetsTalk/api/pkg/rooms"

type User struct {
	Email            string              `bson:"_id" json:"userID"`
	Name             string              `bson:"name" json:"name"`
	DOB              string              `bson:"age" json:"age"`
	Class            string              `bson:"class" json:"class"`
	Password         []byte              `bson:"password" json:"-"`
	PasswordInString string              `bson:"-" json:"password"` // Password when sent by client for registration purposes.
	RoomsJoined      []rooms.RoomsJoined `bson:"roomsJoined" json:"roomsJoined"`
	JoinRequest      []rooms.JoinRequest `bson:"joinRequest" json:"joinRequest"`
	UUID             string              `bson:"loginUUID" json:"uuid"`
}
