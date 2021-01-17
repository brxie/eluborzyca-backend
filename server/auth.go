package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brxie/eluborzyca-backend/controller/session"
	"github.com/brxie/eluborzyca-backend/db/model"
	"github.com/getkin/kin-openapi/openapi3filter"
	fb "github.com/huandu/facebook/v2"
)

func (err *SessionError) Error() string {
	return fmt.Sprintf("HTTP %d %s", err.RespCode, err.Message)
}

type SessionError struct {
	Message  string
	RespCode int
}

const sessionCookieKey = "SESSION_ID"

func extractSession(c context.Context, input *openapi3filter.AuthenticationInput) error {
	var (
		cookie *http.Cookie
		err    error
	)
	r := input.RequestValidationInput.Request

	if cookie, err = r.Cookie(sessionCookieKey); err != nil {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	session, err := session.DecodeSession(cookie.Value)
	if err != nil {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	if session.FacebbokID != "" {
		fbSession := &fb.Session{}
		fbSession.SetAccessToken(session.Token)
		if err := fbSession.Validate(); err != nil {
			return err
		}
		return nil
	}

	_, err = model.GetSession(&model.Session{Token: session.Token, Email: session.Email})
	if err != nil {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	return nil
}
