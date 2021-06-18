package cookies

import (
	"context"
	"net/http"
)

// Cookie interface for cookie generation.
type Cookies interface {
	// Generates cookie provided the responsewriter.
	GenerateCookie(ctx context.Context, w http.ResponseWriter, email, cookieName, path string) error
	// Validates cookie provided the request handle.
	CheckCookie(ctx context.Context, r *http.Request, cookieName string) (string, error)
	// Extend cookie validity.
	ExtendCookie(ctx context.Context, email string) error
}
