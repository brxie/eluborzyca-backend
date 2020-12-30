package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brxie/eluborzyca-backend/config"
	"github.com/brxie/eluborzyca-backend/controller/session"
	"github.com/brxie/eluborzyca-backend/db/model"
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

func GetUrlParamValue(r *http.Request, paramKey string) (string, error) {
	params := r.Context().Value(config.PARAMS).(map[string]string)
	if param, ok := params[paramKey]; ok {
		return param, nil
	}
	return "", fmt.Errorf("Parameter '%s' doesn't exist. ", paramKey)
}
