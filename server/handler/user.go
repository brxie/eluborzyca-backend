package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brxie/ebazarek-backend/controller/user"
	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewUserRequest struct {
	Email      string
	Password   string
	Username   string
	Village    string
	HomeNumber string
	Phone      string
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	user, err := model.GetUser(&model.User{Email: session.Email})
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+err.Error())
		return
	}

	var newUserRequest NewUserRequest
	if err := json.Unmarshal(body, &newUserRequest); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// var passwdCipher string
	passwdCipher, err := user.Encode(newUserRequest.Password)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = model.InsertUser(&model.User{
		Email:      newUserRequest.Email,
		Password:   passwdCipher,
		Username:   newUserRequest.Username,
		Village:    newUserRequest.Village,
		HomeNumber: newUserRequest.HomeNumber,
		Phone:      newUserRequest.Phone,
		Created:    time.Now(),
	})

	if err != nil {
		if e, ok := err.(mongo.WriteException); ok {
			if e.WriteErrors[0].Code == 11000 {
				utils.WriteMessageResponse(&w, http.StatusConflict,
					fmt.Sprintf("Email '%s' already used", newUserRequest.Email))
				return
			}
		}

		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	utils.WriteMessageResponse(&w, http.StatusCreated, http.StatusText(http.StatusCreated))
}
