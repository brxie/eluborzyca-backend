package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brxie/eluborzyca-backend/db/model"
	"github.com/brxie/eluborzyca-backend/utils"
	"github.com/brxie/eluborzyca-backend/utils/ilog"
)

func GetUnits(w http.ResponseWriter, r *http.Request) {
	units, err := model.GetUnits(&model.Unit{})
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := []string{}
	for _, unit := range units {
		resp = append(resp, unit.Name)
	}
	json.NewEncoder(w).Encode(resp)
}
