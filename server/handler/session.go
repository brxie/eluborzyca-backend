package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brxie/ebazarek-backend/config"
	"github.com/brxie/ebazarek-backend/controller/session"
	"github.com/brxie/ebazarek-backend/controller/user"
	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
)

type SessionRequest struct {
	Email    string
	Password string
}

const sessionCookieKey = "SESSION_ID"

func GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := extractSession(r)
	if err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
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
	json.Unmarshal(body, &sessionRequest)

	if err := user.CheckPassword(sessionRequest.Email, sessionRequest.Password); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized)+": "+err.Error())
		return
	}

	token, err := session.NewSession(sessionRequest.Email)
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}

	ttl, err := config.SessionTTL()
	if err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}
	expire := time.Now().Add(time.Duration(int64(ttl) * int64(time.Second)))

	setCookie(&w, sessionCookieKey, token, expire)

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}

func DestroySession(w http.ResponseWriter, r *http.Request) {
	var (
		s   *model.Session
		err error
	)

	if s, err = extractSession(r); err != nil {
		utils.WriteMessageResponse(&w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	if err := session.DestroySession(s.Token); err != nil {
		ilog.Error(err)
	}

	setCookie(&w, sessionCookieKey, "", time.Now())

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))
}
