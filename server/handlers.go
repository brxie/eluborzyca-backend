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
	"GetItem":      handler.GetItem,
	"CreateItem":   handler.CreateItem,
	"ActivateItem": handler.ActivateItem,
	"UpdateItem":   handler.UpdateItem,
	"DeleteItem":   handler.DeleteItem,

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
