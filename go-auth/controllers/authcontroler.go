package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/herulobarto/go-auth/entities"
	"github.com/herulobarto/go-auth/models"
)

type UserInput struct {
	Username string
	Password string
}

var UserModel = models.NewUserModel()

func Index(w http.ResponseWriter, r *http.Request) {

	temp, _ := template.ParseFiles("views/Index.html")
	temp.Execute(w, nil)

}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/Login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		r.ParseForm()
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		var user entities.User
		userModel.where(&user, "username", UserInput.Username)

		var message error
		if user.Username == "" {
			// tidak ditemukan di database
			message = errors.New("Username atau Password salah!")
		} else {
			// pengecekan password
		}

	}

}

// lalu ketik pada terminal go get golang.org.x/crypto/bcrypt
