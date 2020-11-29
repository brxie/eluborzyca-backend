package server

import (
	"net/http"

	"github.com/brxie/ebazarek-backend/server/handler"
)

var Handlers = map[string]func(http.ResponseWriter, *http.Request){
	"GetSession":     handler.GetSession,
	"NewSession":     handler.NewSession,
	"DestroySession": handler.DestroySession,
}
