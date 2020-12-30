package session

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/brxie/eluborzyca-backend/db/model"
)

const sessionTokenLen = 256

func NewSession(email string) (string, error) {
	token := randSeq(sessionTokenLen)
	session := &model.Session{
		Email:   email,
		Created: time.Now(),
		Token:   token,
	}
	err := model.InsertSession(session)

	sessionString, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sessionString), err
}

func GetSession(token string) (*model.Session, error) {
	session := &model.Session{
		Token: token,
	}
	return model.GetSession(session)
}

func DestroySession(token string) error {
	session := &model.Session{
		Token: token,
	}
	return model.DestroySession(session)
}

func DecodeSession(token string) (*model.Session, error) {
	jsonStr, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var session model.Session
	if err := json.Unmarshal(jsonStr, &session); err != nil {
		return nil, err
	}
	return &session, nil

}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
