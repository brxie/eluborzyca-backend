package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func (err *SessionError) Error() string {
	return fmt.Sprintf("HTTP %d %s", err.RespCode, err.Message)
}

type SessionError struct {
	Message  string
	RespCode int
}

func extractSession(c context.Context, input *openapi3filter.AuthenticationInput) error {
	var sessionKey = ""
	req := input.RequestValidationInput.Request
	cookieKey := input.SecurityScheme.Name

	for _, cookie := range req.Cookies() {
		if cookie.Name == cookieKey {
			sessionKey = cookie.Value
			break
		}
	}
	if sessionKey == "" {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	return nil
}
