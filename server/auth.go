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

func extractSession(c context.Context, input *openapi3filter.AuthenticationInput) error {
	var sessionToken = ""
	req := input.RequestValidationInput.Request
	cookieKey := input.SecurityScheme.Name

	for _, cookie := range req.Cookies() {
		if cookie.Name == cookieKey {
			sessionToken = cookie.Value
			break
		}
	}
	if sessionToken == "" {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	_, err := model.GetSession(&model.Session{Token: sessionToken})
	if err != nil {
		return &SessionError{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized}
	}

	return nil
}
