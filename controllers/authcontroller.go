package controllers

import (
	"errors"
	"github.com/taufiqkba/go_auth/config"
	"github.com/taufiqkba/go_auth/entities"
	"github.com/taufiqkba/go_auth/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"text/template"
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
		if err != nil {
			panic(err)
		}
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
