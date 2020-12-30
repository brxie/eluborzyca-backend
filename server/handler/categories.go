package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brxie/eluborzyca-backend/db/model"
	"github.com/brxie/eluborzyca-backend/utils"
	"github.com/brxie/eluborzyca-backend/utils/ilog"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := model.GetCategories(&model.Category{})
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(categories)
}
