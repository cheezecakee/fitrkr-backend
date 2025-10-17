package middleware

import (
	"errors"
	"net/http"

	"github.com/cheezecakee/logr"
)

type Cookie string

const (
	Session      Cookie = "session"
	RefreshToken Cookie = "refresh_token"
)

var ErrTokenNotFound = errors.New("cookie not found")

func ExtractToken(r *http.Request, cookie Cookie) (string, error) {
	if c, err := r.Cookie(string(cookie)); err == nil && c.Value != "" {
		logr.Get().Debugf("found token in cookie: %s", cookie)
		return c.Value, nil
	}

	logr.Get().Info("no token found in cookies")
	return "", ErrTokenNotFound
}
