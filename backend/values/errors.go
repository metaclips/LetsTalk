package values

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrIncorrectUUID              = errors.New("incorrect UUID")
	ErrInvalidExpiryTime          = errors.New("invalid expiry time")
	ErrCookieExpired              = errors.New("generated cookie has expired")
	ErrInvalidUser                = errors.New("invalid user")
	ErrInvalidDetails             = errors.New("invalid signin details")
	ErrRetrieveUUID               = errors.New("could not retrieve UUID")
	ErrMarshal                    = errors.New("could not marshal content")
	ErrWrite                      = errors.New("error while sending content")
	ErrAuthentication             = errors.New("Authentication error")
	ErrIllicitJoinRequest         = errors.New("User was not originally requested to join")
	ErrUserExistInRoom            = errors.New("User already exist in room")
	ErrUserAlreadyRequested       = errors.New("User already requested to join room")
	ErrUserNotRegisteredToRoom    = errors.New("User was not registered to room")
	ErrFileUpload                 = errors.New("error while uploading file to server")
	ErrPeerConnectionNotFound     = errors.New("PeerConnection not found")
	ErrFileUploadLink             = errors.New("could not generate file upload link")
	ErrTokenNotSpecified          = errors.New("Token was not specified")
	ErrInvalidToken               = errors.New("Token specified is invalid")
	ErrPassCompare                = errors.New("error validating users password")
	ErrInvalidUserJoinRoomRequest = errors.New("an invalid users made a join")
	ErrRetrieveUser               = errors.New("user does not exist")
	ErrSetCookie                  = errors.New("error setting cookies")
)

func IsMongoDuplicateError(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
