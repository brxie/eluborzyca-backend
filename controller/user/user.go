package user

import (
	"fmt"

	"github.com/brxie/ebazarek-backend/db/model"
	"github.com/brxie/ebazarek-backend/utils/ilog"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(email, password string) error {
	user, err := model.GetUser(&model.User{Email: email})
	if err != nil {
		ilog.Debug("Get user failed: ", err)
		return fmt.Errorf("Wrong credentials")
	}

	if err := checkPassword(password, user.Password); err != nil {
		ilog.Debug("Check password failed: ", err)
		return fmt.Errorf("Wrong credentials")
	}
	return nil
}

func encode(password string) (string, error) {
	cipher, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(cipher), nil
}

func checkPassword(password, cipher string) error {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(password))
	return err
}
