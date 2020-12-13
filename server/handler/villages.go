package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
)

func GetVillages(w http.ResponseWriter, r *http.Request) {
	villages, err := model.GetVillages(&model.Village{})
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := []string{}
	for _, village := range villages {
		resp = append(resp, village.Name)
	}
	json.NewEncoder(w).Encode(resp)
}
