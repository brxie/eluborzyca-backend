package handler

import (
	"net/http"
	"time"

	"github.com/brxie/ebazarek-backend/controller/session"
	"github.com/brxie/ebazarek-backend/db/model"
)

func setCookie(w *http.ResponseWriter, name, value string, expire time.Time) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(*w, &cookie)
}

func extractSession(r *http.Request) (*model.Session, error) {
	var (
		cookie *http.Cookie
		err    error
	)

	if cookie, err = r.Cookie(sessionCookieKey); err != nil {
		return nil, err
	}

	return session.DecodeSession(cookie.Value)
}
