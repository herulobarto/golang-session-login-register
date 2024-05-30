package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/herulobarto/go-auth/config"
	"github.com/herulobarto/go-auth/entities"
	"github.com/herulobarto/go-auth/models"
	"golang.org/x/crypto/bcrypt"

	en_translations "github.com/go-playground/validator/v10/translations/en"
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

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		// melakukan proses registrasi

		// mengambil inputan form
		r.ParseForm()

		user := entities.User{
			NamaLengkap: r.Form.Get("nama_lengkap"),
			Email:       r.Form.Get("email"),
			Username:    r.Form.Get("username"),
			Password:    r.Form.Get("password"),
			Cpassword:   r.Form.Get("cpassword"),
		}

		// errorMessages := make(map[string]interface{})

		// if user.NamaLengkap == "" {
		// 	errorMessages["NamaLengkap"] = "Nama Lengkap harus diisi"

		// }
		// if user.Email == "" {
		// 	errorMessages["Email"] = "Email harus diisi"

		// }
		// if user.Username == "" {
		// 	errorMessages["Username"] = "Username harus diisi"

		// }
		// if user.Password == "" {
		// 	errorMessages["Password"] = "Password harus diisi"

		// }
		// if user.Cpassword == "" {
		// 	errorMessages["Cpassword"] = "Konfirmasi Password harus diisi"

		// } else {
		// 	if user.Cpassword != user.Password {
		// 		errorMessages["Cpassword"] = "Konfirmasi password tidak cocok"
		// 	}
		// }

		// if len(errorMessages) > 0 {
		// 	//  validasi form gagal

		// 	data := map[string]interface{}{
		// 		"validation": errorMessages,
		// 	}

		// 	temp, _ := template.ParseFiles("views/register.html")
		// 	temp.Execute(w, data)
		// } else {
		// 	// hash password menggunakan bcrypt
		// 	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		// 	user.Password = string(hashPassword)

		// 	// insert ke database
		// 	_, err := userModel.Create(user)

		// 	var message string
		// 	if err != nil {
		// 		message = "Proses Registrasi Gagal: " + message
		// 	} else {
		// 		message = "Registrasi Berhasil, Silahkan Login"
		// 	}

		// 	data := map[string]interface{}{
		// 		"pesan": message,
		// 	}

		// 	temp, _ := template.ParseFiles("views/register.html")
		// 	temp.Execute(w, data)
		// }

		// memanggil paket translator
		translator := en.New()
		uni := ut.New(translator, translator)

		trans, _ := uni.GetTranslator("en")

		validate := validator.New()

		// register default translation (en)
		en_translations.RegisterDefaultTranslations(validate, trans)

		// mengubah label default
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			labelName := field.Tag.Get("label")
			return labelName
		})

		// memakai pesan bahasa indonesia
		validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} tidak boleh kosong", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})

		validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
			params := fl.Param()
			split_params := strings.Split(params, "-")

			tableName := split_params[0]
			fieldName := split_params[1]
			fieldValue := fl.Field().String()

			return checkIsUnique(tableName, fieldName, fieldValue)
		})

		validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
			return ut.Add("isunique", "{0} sudah digunakan", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isunique", fe.Field())
			return t
		})

		// melakukan proses validasi
		vErrors := validate.Struct(user)

		errorMessages := make(map[string]interface{})
		if vErrors != nil {
			fmt.Println(vErrors)

			for _, e := range vErrors.(validator.ValidationErrors) {
				errorMessages[e.StructField()] = e.Translate(trans)
			}

			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		}

	}

}

func checkIsUnique(tableName, fieldName, fieldValue string) bool {

	conn, err := config.DBConn()
	if err != nil {
		panic(err)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", fieldName, tableName, fieldName)
	row, err2 := conn.Query(query, fieldValue)
	if err2 != nil {
		panic(err)
	}
	defer row.Close()

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	// jika ketemu email yang sama akan dibandingkan dengan filedvalue yang diisi
	return result != fieldValue
}
