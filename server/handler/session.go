package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brxie/eluborzyca-backend/config"
	"github.com/brxie/eluborzyca-backend/controller/session"
	"github.com/brxie/eluborzyca-backend/controller/user"
	"github.com/brxie/eluborzyca-backend/db/model"
	"github.com/brxie/eluborzyca-backend/utils"
	"github.com/brxie/eluborzyca-backend/utils/ilog"
	fb "github.com/huandu/facebook/v2"
)

type SessionRequest struct {
	Email    string
	Password string
}

type FacebookSessionRequest struct {
	Name                     string `json:"name"`
	Email                    string `json:"email"`
	ID                       string `json:"id"`
	AccessToken              string `json:"accessToken"`
	UserID                   string `json:"userID"`
	ExpiresIn                int    `json:"expiresIn"`
	SignedRequest            string `json:"signedRequest"`
	GraphDomain              string `json:"graphDomain"`
	DataAccessExpirationTime int    `json:"data_access_expiration_time"`
}

const sessionCookieKey = "SESSION_ID"

func GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	session.Token = ""
	json.NewEncoder(w).Encode(session)
}

func NewSession(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+err.Error())
		return
	}
	var sessionRequest SessionRequest
	if err := json.Unmarshal(body, &sessionRequest); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized)+": "+err.Error())
		return
	}

	if err := user.CheckPassword(sessionRequest.Email, sessionRequest.Password); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized))
		return
	}

	verified, err := user.IsVerified(sessionRequest.Email)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+err.Error())
		return
	}

	if !verified {
		utils.WriteMessageResponse(&w, http.StatusForbidden, "User not verified")
		return
	}

	sessionToken, err := session.NewSession(sessionRequest.Email)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}

	ttl := config.Viper.GetInt64("SESSION_TOKEN_TTL")
	expire := time.Now().Add(time.Duration(ttl * int64(time.Second)))

	setCookie(&w, sessionCookieKey, sessionToken, expire)
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}

func NewFacebookSession(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+err.Error())
		return
	}
	var facebookSessionRequest FacebookSessionRequest
	if err := json.Unmarshal(body, &facebookSessionRequest); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized)+": "+err.Error())
		return
	}

	fbSession := &fb.Session{}
	fbSession.SetAccessToken(facebookSessionRequest.AccessToken)
	if err := fbSession.Validate(); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized)+": "+err.Error())
		return
	}

	if err := user.CreateIfNotExist(&model.User{
		Email:      facebookSessionRequest.Email,
		Username:   facebookSessionRequest.Name,
		FacebookID: facebookSessionRequest.UserID,
		Verified:   true,
		Created:    time.Now(),
	}); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}

	sessionToken, err := session.NewFacebookSession(facebookSessionRequest.Email,
		facebookSessionRequest.AccessToken, facebookSessionRequest.UserID)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}

	expireTime := time.Unix(int64(facebookSessionRequest.DataAccessExpirationTime), 0)
	setCookie(&w, sessionCookieKey, sessionToken, expireTime)
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}

func DestroySession(w http.ResponseWriter, r *http.Request) {
	var (
		s   *model.Session
		err error
	)

	if s, err = extractSession(r); err != nil {
		utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
		return
	}

	if err := session.DestroySession(s.Token); err != nil {
		ilog.Error(err)
	}

	setCookie(&w, sessionCookieKey, "", time.Now())
	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
