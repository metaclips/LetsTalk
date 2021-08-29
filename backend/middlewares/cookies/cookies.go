package cookies

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cookie interface for cookie generation.
type Cookies interface {
	// Generates cookie provided the responsewriter.
	GenerateCookie(ctx context.Context, email, cookieName string) (string, string, error)
	// Validates cookie provided the request handle.
	CheckCookie(ctx context.Context, r *http.Request, cookieName string) (*primitive.ObjectID, error)
	// Extend cookie validity.
	ExtendCookie(ctx context.Context, email string) error
}
