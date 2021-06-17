package cookies

import (
	"net/http"
)

// Cookie interface for cookie generation.
type Cookies interface {
	// Generates cookie provided the responsewriter.
	GenerateCookie(w http.ResponseWriter, email, cookieName, path string) error
	// Validates cookie provided the request handle.
	CheckCookie(r *http.Request, cookieName string) error
	// Extend cookie validity.
	ExtendCookie(email string) error
}
