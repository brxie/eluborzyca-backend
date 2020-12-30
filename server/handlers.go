package server

import (
	"net/http"

	"github.com/brxie/eluborzyca-backend/server/handler"
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
	"VerifyUser": handler.VerifyUser,

	// item
	"GetItem":      handler.GetItem,
	"CreateItem":   handler.CreateItem,
	"ActivateItem": handler.ActivateItem,
	"UpdateItem":   handler.UpdateItem,
	"DeleteItem":   handler.DeleteItem,

	// items
	"GetItems": handler.GetItems,

	// image
	"GetImage":    handler.GetImage,
	"UploadImage": handler.UploadImage,

	// units
	"GetUnits": handler.GetUnits,

	// categories
	"GetCategories": handler.GetCategories,

	// villages
	"GetVillages": handler.GetVillages,
}
