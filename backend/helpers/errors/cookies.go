package errors

import "errors"

var (
	ErrTokenNotSpecified     = errors.New("token was not specified")
	ErrInvalidToken          = errors.New("invalid token")
	ErrCookieExpired         = errors.New("cookie expired")
	ErrIncorrectUUID         = errors.New("incorrect uuid")
	ErrInvalidExpiryTime     = errors.New("invalid expiry time")
	ErrValidatingCredentials = errors.New("error validating credentials")
)
