package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"text/template"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/taufiqkba/go_auth/config"
	"github.com/taufiqkba/go_auth/entities"
	"github.com/taufiqkba/go_auth/models"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string
	Password string
}

var UserModel = models.NewUserModel()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)
	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			data := map[string]interface{}{
				"name": session.Values["name"],
			}

			temp, err := template.ParseFiles("views/index.gohtml")
			err = temp.Execute(w, data)
			if err != nil {
				panic(err)
			}
		}
	}

}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/login.gohtml")
		err = temp.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else if r.Method == http.MethodPost {
		//login process
		err := r.ParseForm()
		if err != nil {
			return
		}
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
		var user entities.User
		err = UserModel.Where(&user, "username", UserInput.Username)
		if err != nil {
			return
		}

		var message error
		if user.Username == "" {
			//	user record not found
			message = errors.New("username or password is invalid")
		} else {
			//	check password
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
			if errPassword != nil {
				message = errors.New("username or password is invalid")
			}
		}

		if message != nil {
			data := map[string]interface{}{
				"error": message,
			}
			temp, err := template.ParseFiles("views/login.gohtml")
			if err != nil {
				panic(err)
			}
			err = temp.Execute(w, data)
			if err != nil {
				panic(err)
			}
		} else {
			//	login using session
			session, err := config.Store.Get(r, config.SESSION_ID)
			session.Values["loggedIn"] = true
			session.Values["email"] = user.Email
			session.Values["username"] = user.Username
			session.Values["name"] = user.Name
			err = session.Save(r, w)
			if err != nil {
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	//	delete session
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		//	direct to register view page
		temp, err := template.ParseFiles("views/register.gohtml")
		err = temp.Execute(w, nil)
		if err != nil {
			return
		}
	} else if r.Method == http.MethodPost {
		//	registration process
		err := r.ParseForm()
		if err != nil {
			return
		}
		user := entities.User{
			Name:      r.Form.Get("name"),
			Email:     r.Form.Get("email"),
			Username:  r.Form.Get("username"),
			Password:  r.Form.Get("password"),
			CPassword: r.Form.Get("c_password"),
		}

		//errorMessage := make(map[string]interface{})
		//if user.Name == "" {
		//	errorMessage["Name"] = "Name can't empty"
		//}
		//if user.Email == "" {
		//	errorMessage["Email"] = "Email can't empty"
		//}
		//if user.Username == "" {
		//	errorMessage["Username"] = "Username can't empty"
		//}
		//if user.Password == "" {
		//	errorMessage["Password"] = "Password can't empty"
		//}
		//if user.CPassword == "" {
		//	errorMessage["CPassword"] = "Confirmation Password can't empty"
		//} else {
		//	if user.CPassword != user.Password {
		//		errorMessage["CPassword"] = "Confirmation Password not match"
		//	}
		//}
		//if len(errorMessage) > 0 {
		//	//	validation form failes
		//	data := map[string]interface{}{
		//		"validation": errorMessage,
		//	}
		//	temp, err := template.ParseFiles("views/register.gohtml")
		//	err = temp.Execute(w, data)
		//	if err != nil {
		//		return
		//	}
		//} else {
		//	//	hash password using bcrypt
		//	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		//	if err != nil {
		//		return
		//	}
		//	user.Password = string(hashPassword) //return back to entities
		//
		//	//	insert to database
		//	_, err = UserModel.Create(user)
		//	var message string
		//	if err != nil {
		//		message = "Register failed: " + message
		//	} else {
		//		message = "Registration success, click login"
		//	}
		//
		//	data := map[string]interface{}{
		//		"message": message,
		//	}
		//
		//	temp, _ := template.ParseFiles("views/register.gohtml")
		//	err = temp.Execute(w, data)
		//	if err != nil {
		//		return
		//	}
		//}

		//translator error
		translator := en.New()
		uni := ut.New(translator, translator)
		trans, _ := uni.GetTranslator("en")

		//	validation using go playground
		validate := validator.New()

		//register default translation english
		err = en2.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return
		}

		//change default label error
		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			labelName := field.Tag.Get("label")
			return labelName
		})

		validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} can't empty", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})

		//validation process
		vErrors := validate.Struct(user)

		errorMessages := make(map[string]interface{})
		if vErrors != nil {
			for _, e := range vErrors.(validator.ValidationErrors) {
				errorMessages[e.StructField()] = e.Translate(trans)
			}
			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}
			temp, _ := template.ParseFiles("views/register.gohtml")
			err := temp.Execute(w, data)
			if err != nil {
				return
			}
		}
	}
}
