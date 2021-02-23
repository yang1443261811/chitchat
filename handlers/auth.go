package handlers

import (
	"fmt"
    "chitchat/models"
    "net/http"
)

// GET /login
// 登录页面
func Login(writer http.ResponseWriter, request *http.Request) {
    t := parseTemplateFiles("auth.layout", "navbar", "login")
    t.Execute(writer, nil)
}

// GET /signup
// 注册页面
func Signup(writer http.ResponseWriter, request *http.Request) {
    generateHTML(writer, nil, "auth.layout", "navbar", "signup")
}

func SignupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
    }
	user := models.User{
		Name: request.PostFormValue("name"),
		Email: request.PostFormValue("email"),
		Password: request.PostFormValue("password"),

	}
	ok := user.Create()
	if ok != nil {
		danger(ok, "Cannot create user")
	} else {
		http.Redirect(writer, request, "/login", 302)
	}
}

func Authenticate(writer http.ResponseWriter, request *http.Request)  {
	request.ParseForm()
	user, err := models.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "cannot find user")
	}
	if user.Password == models.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
            fmt.Println("Cannot create session")
        }
		cookie := http.Cookie{
			Name: "_cookie",
			Value: session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
        http.Redirect(writer, request, "/", 302)
	} else {
		http.Redirect(writer, request, "/login", 302)
	}
}

// GET /logout
// 用户退出
func Logout(writer http.ResponseWriter, request *http.Request) {
    cookie, err := request.Cookie("_cookie")
    if err != http.ErrNoCookie {
        fmt.Println("Failed to get cookie")
        session := models.Session{Uuid: cookie.Value}
        session.DeleteByUUID()
    }
    http.Redirect(writer, request, "/", 302)
}