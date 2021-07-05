package cookies

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type secureCookie struct {
	handler *securecookie.SecureCookie
	db      database.Database
}

type cookieData struct {
	ID      primitive.ObjectID
	Email   string
	UUID    string
	Expires time.Time
}

const maxTime = time.Hour * 2

func NewSecureCookie(userDB database.Database) Cookies {
	return &secureCookie{
		securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
		userDB,
	}
}

func (b *secureCookie) GenerateCookie(ctx context.Context, email, cookieName string) (string, error) {
	cookie := cookieData{}
	exitTime := time.Now().Add(maxTime)
	cookie.UUID = uuid.New().String()
	cookie.Email = email

	// Update user UUID.
	err := b.db.UpdateByID(ctx, email, bson.M{"$set": bson.M{"uuid": cookie.UUID, "expires": exitTime}})
	if err != nil {
		return "", err
	}

	encoded, err := b.handler.Encode(cookieName, cookie)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func (b *secureCookie) CheckCookie(ctx context.Context, r *http.Request, cookieName string) (*primitive.ObjectID, error) {
	authorization := r.Header.Get(constants.AuthorizationHeader)

	bearer := "Bearer"
	if len(authorization) > len(bearer) && strings.ToUpper(authorization[0:len(bearer)]) == bearer {
		authorization = authorization[len(bearer)+1:]
	}

	if authorization == "" {
		return nil, api_errors.ErrTokenNotSpecified
	}

	data := cookieData{}

	if err := b.handler.Decode(cookieName, authorization, &data); err != nil {
		return nil, api_errors.ErrInvalidToken
	}

	// Validate
	dbCookieData := cookieData{}
	if err := b.db.FindByID(ctx, data.Email, nil, &dbCookieData); err != nil {
		return nil, err
	}

	// Ensure UUID in DB is same with that of user
	// else there's a login from two different browsers
	// which should not be allowed.
	if dbCookieData.UUID != data.UUID {
		return nil, api_errors.ErrIncorrectUUID
	}

	// They must be of same tume.
	if dbCookieData.Expires.Sub(dbCookieData.Expires).Seconds() > 0 {
		return nil, api_errors.ErrInvalidExpiryTime
	}

	return &dbCookieData.ID, b.ExtendCookie(ctx, data.Email)
}

func (b *secureCookie) ExtendCookie(ctx context.Context, email string) error {
	// Update user UUID.
	err := b.db.UpdateByID(ctx, email, bson.M{"$set": bson.M{"expires": time.Now().Add(maxTime)}})
	if err != nil {
		return err
	}

	return nil
}
