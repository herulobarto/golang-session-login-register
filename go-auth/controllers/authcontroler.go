package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/herulobarto/go-auth/config"
	"github.com/herulobarto/go-auth/entities"
	"github.com/herulobarto/go-auth/models"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string
	Password string
}

var userModel = models.NewUserModel()

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if session.Values["loggedIn"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"nama_lengkap": session.Values["nama_lengkap"],
	}

	temp, _ := template.ParseFiles("views/index.html")
	temp.Execute(w, data)
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
		userModel.Where(&user, "username", UserInput.Username)

		var message error
		if user.Username == "" {
			// tidak ditemukan di database
			message = errors.New("username atau password salah")
		} else {
			// pengecekan password
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
			if errPassword != nil {
				message = errors.New("username atau password salah")
			}
		}

		if message != nil {

			data := map[string]interface{}{
				"error": message,
			}

			temp, _ := template.ParseFiles("views/Login.html")
			temp.Execute(w, data)
		} else {
			// set session

			session, _ := config.Store.Get(r, config.SESSION_ID)

			session.Values["loggedIn"] = true
			session.Values["email"] = user.Email
			session.Values["username"] = user.Username
			session.Values["nama_lengkap"] = user.NamaLengkap

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	// delete SESSION
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
