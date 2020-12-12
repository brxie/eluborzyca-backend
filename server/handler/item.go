package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemRequest struct {
	Name          string
	Price         uint64
	Unit          string
	Availability  int
	FirstLastName string
	Village       string
	HomeNumber    string
	Phone         string
	Category      string
	Description   string
	Popular       bool
	Active        bool
	Images        []string
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := GetUrlParamValue(r, "itemID")
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		return
	}

	id, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusNotFound,
			http.StatusText(http.StatusNotFound))
		return
	}
	item, err := model.GetItem(&model.Item{ID: id})
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusNotFound,
			http.StatusText(http.StatusNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	var itemRequest ItemRequest
	if err := json.Unmarshal(body, &itemRequest); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if _, err := model.GetUnit(&model.Unit{Name: itemRequest.Unit}); err != nil {
		ilog.Warn(err)
		utils.WriteMessageResponse(&w, http.StatusBadRequest, "Unit doesn't exists")
		return
	}

	if _, err := model.GetVillage(&model.Village{Name: itemRequest.Village}); err != nil {
		ilog.Warn(err)
		utils.WriteMessageResponse(&w, http.StatusBadRequest, "Village doesn't exists")
		return
	}

	if _, err := model.GetCategory(&model.Category{Name: itemRequest.Category}); err != nil {
		ilog.Warn(err)
		utils.WriteMessageResponse(&w, http.StatusBadRequest, "Category doesn't exists")
		return
	}

	var images []model.Image
	for _, val := range itemRequest.Images {
		id, err := primitive.ObjectIDFromHex(val)
		if err != nil {
			ilog.Warn(err)
			utils.WriteMessageResponse(&w, http.StatusBadRequest, "Wrong image id")
			return
		}
		image, err := model.GetImage(&model.Image{ID: id})
		if err != nil {
			ilog.Warn(err)
			utils.WriteMessageResponse(&w, http.StatusBadRequest, "Image doesn't exists")
			return
		}
		images = append(images, *image)
	}

	item := &model.Item{
		Name:          itemRequest.Name,
		Owner:         session.Email,
		Price:         itemRequest.Price,
		Unit:          itemRequest.Unit,
		Availability:  itemRequest.Availability,
		FirstLastName: itemRequest.FirstLastName,
		Village:       itemRequest.Village,
		HomeNumber:    itemRequest.HomeNumber,
		Phone:         itemRequest.Phone,
		Category:      itemRequest.Category,
		Description:   itemRequest.Description,
		Popular:       itemRequest.Popular,
		Active:        itemRequest.Active,
		Images:        images,
		Created:       time.Now(),
	}
	if err := model.InsertItem(item); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	utils.WriteMessageResponse(&w, http.StatusCreated, http.StatusText(http.StatusCreated))
}
