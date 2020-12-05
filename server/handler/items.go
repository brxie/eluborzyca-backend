package handler

import (
	"net/http"

	"github.com/brxie/ebazarek-backend/utils"
)

func GetItems(w http.ResponseWriter, r *http.Request) {

	utils.WriteMessageResponse(&w, http.StatusCreated, http.StatusText(http.StatusCreated))
}
