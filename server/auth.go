package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/getkin/kin-openapi/openapi3filter"
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

	if cookie.Value == "" {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	_, err = model.GetSession(&model.Session{Token: cookie.Value})
	if err != nil {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	return nil
}
