package handler

import (
	"net/http"

	"github.com/brxie/ebazarek-backend/utils"
)

func GetUnits(w http.ResponseWriter, r *http.Request) {

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
