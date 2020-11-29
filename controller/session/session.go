package session

import (
	"math/rand"
	"time"

	"github.com/brxie/ebazarek-backend/db/model"
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

	return token, err
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

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
