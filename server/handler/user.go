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

type UserRequest struct {
	Email       string
	NewPassword string
	Password    string
	Username    string
	Village     string
	HomeNumber  string
	Phone       string
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
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	var userRequest UserRequest
	if err := json.Unmarshal(body, &userRequest); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if _, err := model.GetVillage(&model.Village{Name: userRequest.Village}); err != nil {
		ilog.Warn(err)
		utils.WriteMessageResponse(&w, http.StatusBadRequest, "Village doesn't exists")
		return
	}

	passwdCipher, err := user.Encode(userRequest.NewPassword)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = model.InsertUser(&model.User{
		Email:      userRequest.Email,
		Password:   passwdCipher,
		Username:   userRequest.Username,
		Village:    userRequest.Village,
		HomeNumber: userRequest.HomeNumber,
		Phone:      userRequest.Phone,
		Created:    time.Now(),
	})

	if err != nil {
		if e, ok := err.(mongo.WriteException); ok {
			if e.WriteErrors[0].Code == 11000 {
				utils.WriteMessageResponse(&w, http.StatusConflict,
					fmt.Sprintf("Email '%s' already used", userRequest.Email))
				return
			}
		}

		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	utils.WriteMessageResponse(&w, http.StatusCreated, http.StatusText(http.StatusCreated))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	var userRequest UserRequest
	if err := json.Unmarshal(body, &userRequest); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		return
	}

	if userRequest.NewPassword == userRequest.Password {
		utils.WriteMessageResponse(&w, http.StatusBadRequest,
			"New password must differ")
		return
	}

	if userRequest.NewPassword != "" {
		if err := user.CheckPassword(userRequest.Email, userRequest.Password); err != nil {
			utils.WriteMessageResponse(&w, http.StatusBadRequest,
				"Old password is incorrect")
			return
		}

		if userRequest.NewPassword, err = user.Encode(userRequest.NewPassword); err != nil {
			ilog.Error(err)
			utils.WriteMessageResponse(&w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError))
			return
		}

	}

	update := &model.User{
		Username:   userRequest.Username,
		Village:    userRequest.Village,
		HomeNumber: userRequest.HomeNumber,
		Phone:      userRequest.Phone,
		Password:   userRequest.NewPassword,
	}

	if err := model.UpdateUser(&model.User{Email: session.Email}, update); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+": "+err.Error())
		return
	}

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
