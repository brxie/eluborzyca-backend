package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brxie/ebazarek-backend/config"
	"github.com/brxie/ebazarek-backend/controller/session"
	"github.com/brxie/ebazarek-backend/controller/user"
	"github.com/brxie/ebazarek-backend/utils"
	"github.com/brxie/ebazarek-backend/utils/ilog"
)

type SessionRequest struct {
	Email    string
	Password string
}

const sessionCookieKey = "SESSION_ID"

func TestSession(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
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

	if err := addCookie(&w, sessionCookieKey, token); err != nil {
		ilog.Error(err)
		utils.WriteMessageResponse(&w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError)+" "+err.Error())
		return
	}

	utils.WriteMessageResponse(&w, http.StatusOK, http.StatusText(http.StatusOK))

}

func addCookie(w *http.ResponseWriter, name, value string) error {
	ttl, err := config.SessionTTL()
	if err != nil {
		return err
	}

	expire := time.Now().Add(time.Duration(int64(ttl) * int64(time.Second)))
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(*w, &cookie)
	return nil
}
