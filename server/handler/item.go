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
	Images        []string
}

type ItemRequestUpdate struct {
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

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

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

	if item.Owner != session.Email {
		utils.WriteMessageResponse(&w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	var itemRequest ItemRequestUpdate
	if err := json.Unmarshal(body, &itemRequest); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if (ItemRequestUpdate{}) == itemRequest {
		utils.WriteMessageResponse(&w, http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest))
		return
	}

	if itemRequest.Unit != "" {
		if _, err := model.GetUnit(&model.Unit{Name: itemRequest.Unit}); err != nil {
			ilog.Warn(err)
			utils.WriteMessageResponse(&w, http.StatusBadRequest, "Unit doesn't exists")
			return
		}
	}

	if itemRequest.Village != "" {
		if _, err := model.GetVillage(&model.Village{Name: itemRequest.Village}); err != nil {
			ilog.Warn(err)
			utils.WriteMessageResponse(&w, http.StatusBadRequest, "Village doesn't exists")
			return
		}
	}

	if itemRequest.Category != "" {
		if _, err := model.GetCategory(&model.Category{Name: itemRequest.Category}); err != nil {
			ilog.Warn(err)
			utils.WriteMessageResponse(&w, http.StatusBadRequest, "Category doesn't exists")
			return
		}
	}

	itemUpdate := &model.ItemUpdate{
		Name:          itemRequest.Name,
		Price:         itemRequest.Price,
		Unit:          itemRequest.Unit,
		Availability:  itemRequest.Availability,
		FirstLastName: itemRequest.FirstLastName,
		Village:       itemRequest.Village,
		HomeNumber:    itemRequest.HomeNumber,
		Phone:         itemRequest.Phone,
		Category:      itemRequest.Category,
		Description:   itemRequest.Description,
	}

	if err := model.UpdateItem(&model.Item{ID: id}, itemUpdate); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}

func ActivateItem(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	var activateItem model.ItemActivate
	if err := json.Unmarshal(body, &activateItem); err != nil {
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

	if item.Owner != session.Email {
		utils.WriteMessageResponse(&w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	if err := model.ActivateItem(&model.Item{ID: id}, &activateItem); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
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
		Images:        images,
		Created:       time.Now(),
		Popular:       true,
	}
	if err := model.InsertItem(item); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	utils.WriteMessageResponse(&w, http.StatusCreated, http.StatusText(http.StatusCreated))
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

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

	if item.Owner != session.Email {
		utils.WriteMessageResponse(&w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	if err := model.DeleteItem(&model.Item{ID: id}); err != nil {
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
