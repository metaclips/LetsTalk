package cookies

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	api_constants "github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/models/database"
	"go.mongodb.org/mongo-driver/bson"
)

type secureCookie struct {
	handler *securecookie.SecureCookie
	db      database.Database
}

type cookieData struct {
	Email   string
	UUID    string
	Expires time.Time
}

var ctx, _ = context.CreateChildCancelContext()

const maxTime = time.Hour * 2

func NewSecureCookie(userDB, adminDB database.Database) Cookies {
	return &secureCookie{
		securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)),
		userDB,
	}
}

// GenerateCookie generates a token and sends to frontend.
// Generated tokens are to be stored by frontend and verified for every request.
func (b *secureCookie) GenerateCookie(w http.ResponseWriter, email, cookieName, path string) error {
	cookie := cookieData{}
	exitTime := time.Now().Add(maxTime)
	cookie.UUID = uuid.New().String()
	cookie.Email = email

	// Update user UUID.
	err := b.db.UpdateByID(ctx, email, bson.M{"$set": bson.M{"uuid": cookie.UUID, "expires": exitTime}})
	if err != nil {
		return err
	}

	encoded, err := b.handler.Encode(cookieName, cookie)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     api_constants.UserCookieName,
		Value:    encoded,
		Path:     path,
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}

func (b *secureCookie) CheckCookie(r *http.Request, cookieName string) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	if cookie.Value == "" {
		return api_errors.ErrTokenNotSpecified
	}

	data := cookieData{}

	if err := b.handler.Decode(cookieName, cookie.Value, &data); err != nil {
		return api_errors.ErrInvalidToken
	}

	// Ensure hasn't expired.
	if cookie.Expires.Before(time.Now().Local()) {
		return api_errors.ErrCookieExpired
	}

	// Validate
	dbCookieData := cookieData{}
	if err := b.db.FindByID(ctx, data.Email, nil, &dbCookieData); err != nil {
		return err
	}

	// Ensure UUID in DB is same with that of user
	// else there's a login from two different browsers
	// which should not be allowed.
	if dbCookieData.UUID != data.UUID {
		return api_errors.ErrIncorrectUUID
	}

	// They must be of same tume.
	if dbCookieData.Expires.Sub(dbCookieData.Expires).Seconds() > 0 {
		return api_errors.ErrInvalidExpiryTime
	}

	return b.ExtendCookie(data.Email)
}

func (b *secureCookie) ExtendCookie(email string) error {
	// Update user UUID.
	err := b.db.UpdateByID(ctx, email, bson.M{"$set": bson.M{"expires": time.Now().Add(maxTime)}})
	if err != nil {
		return err
	}

	return nil
}
