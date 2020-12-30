package user

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/brxie/eluborzyca-backend/config"
	"github.com/brxie/eluborzyca-backend/db/model"
	"github.com/brxie/eluborzyca-backend/utils/ilog"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type EmailTemplate struct {
	URL, Token string
}

func CheckPassword(email, password string) error {
	user, err := model.GetUser(&model.User{Email: email})
	if err != nil {
		ilog.Debug("Get user failed: ", err)
		return err
	}

	if err := checkPassword(password, user.Password); err != nil {
		ilog.Debug("Check password failed: ", err)
		return fmt.Errorf("Wrong credentials")
	}
	return nil
}

func IsVerified(email string) (bool, error) {
	user, err := model.GetUser(&model.User{Email: email})
	if err != nil {
		ilog.Debug("Get user failed: ", err)
		return false, err
	}

	return user.Verified, nil
}

func Encode(password string) (string, error) {
	cipher, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(cipher), nil
}

func SendVeryfyTokenEmail(token *model.VerifyToken) error {
	host := config.Viper.GetString("SMTP_HOST")
	port := config.Viper.GetInt("SMTP_PORT")
	user := config.Viper.GetString("SMTP_USER")
	pass := config.Viper.GetString("SMTP_PASSWORD")
	sender := config.Viper.GetString("SMTP_SENDER_NAME")
	tokenURL := config.Viper.GetString("FRONTEND_URL") + "/user-verify"
	skipCertCheck := config.Viper.GetBool("SMTP_PASSWORD")

	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: skipCertCheck}

	m := gomail.NewMessage()

	m.SetHeader("From", sender)
	m.SetHeader("To", token.Email)
	m.SetHeader("Subject", "Verify Token")

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	t, err := template.ParseFiles(path.Join(ex, "verifyTokenEmail.html"))
	if err != nil {
		return err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, &EmailTemplate{tokenURL, token.Token})
	if err != nil {
		return err
	}
	m.SetBody("text/html", out.String())

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func checkPassword(password, cipher string) error {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(password))
	return err
}
