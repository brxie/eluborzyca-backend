package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user, err := model.GetUser(&model.User{Email: session.Email})
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user.Password = ""
	json.NewEncoder(w).Encode(user)
}
