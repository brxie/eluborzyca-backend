package handler

import (
	"net/http"

	"github.com/brxie/ebazarek-backend/utils"
)

func GetVillages(w http.ResponseWriter, r *http.Request) {

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
