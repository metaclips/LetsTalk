package model

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/metaclips/LetsTalk/backend/values"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

// GenerateCookie generates a token and sends to frontend.
// Generated tokens are to be stored by frontend and verified for every request.
func (b CookieDetail) GenerateCookie(w http.ResponseWriter) error {
	exitTime := time.Now().Add(time.Hour * 2).Local()
	b.Data.ExitTime = exitTime
	b.Data.UUID = uuid.New().String()

	_, err := db.Collection(b.Collection).UpdateOne(ctx, bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"loginUUID": b.Data.UUID, "expires": exitTime}})
	if err != nil {
		return err
	}

	encoded, err := cookieHandler.Encode(b.CookieName, b.Data)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     values.UserCookieName,
		Value:    encoded,
		Path:     b.Path,
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}

func (b *CookieDetail) CheckCookie(r *http.Request, cookieName string) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	if cookie.Value == "" {
		return values.ErrTokenNotSpecified
	}

	if err := cookieHandler.Decode(b.CookieName, cookie.Value, &b.Data); err != nil {
		return values.ErrInvalidToken
	}

	if b.Data.ExitTime.Before(time.Now().Local()) {
		return values.ErrCookieExpired
	}

	b.Email = b.Data.Email
	result := db.Collection(b.Collection).FindOne(ctx, bson.M{"_id": b.Email})

	cookieUUID := struct {
		LoginUUID string    `json:"loginUUID"`
		Expires   time.Time `json:"expires"`
	}{}

	if err := result.Decode(&cookieUUID); err != nil {
		return err
	}

	if cookieUUID.LoginUUID != b.Data.UUID {
		return values.ErrIncorrectUUID
	}

	if cookieUUID.Expires.Sub(b.Data.ExitTime).Seconds() > 0 {
		return values.ErrInvalidExpiryTime
	}

	return nil
}
