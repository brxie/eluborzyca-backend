package server

import (
	"net/http"

	"github.com/brxie/ebazarek-backend/server/handler"
)

var Handlers = map[string]func(http.ResponseWriter, *http.Request){
	// session
	"GetSession":     handler.GetSession,
	"NewSession":     handler.NewSession,
	"DestroySession": handler.DestroySession,

	// user
	"GetUser":    handler.GetUser,
	"CreateUser": handler.CreateUser,
	"UpdateUser": handler.UpdateUser,

	// item
	"CreateItem": handler.CreateItem,
	"GetItem":    handler.GetItem,

	// items
	"GetItems": handler.GetItems,

	// images
	"GetImages": handler.GetImages,

	// units
	"GetUnits": handler.GetUnits,

	// categories
	"GetCategories": handler.GetCategories,

	// villages
	"GetVillages": handler.GetVillages,
}
