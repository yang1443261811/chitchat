package handlers

import (
    "chitchat/models"
    "net/http"
)

func NewThread(writer http.ResponseWriter, request *http.Request)  {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "login", 302)
	} else {
		generateHTML(writer, nil, "layout", "auth.navbar", "new.thread")
	}
}

func CreateThread(writer http.ResponseWriter, request *http.Request)  {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
        if err != nil {
            danger(err, "Cannot parse form")
        }
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot create user")
		}
		topic := request.PostFormValue("topic")
		_, err = user.CreateThread(topic)
		if err != nil {
			danger(err, "Cannot create thread")
		}
		http.Redirect(writer, request, "/", 302)
	}
}

func ReadThread(writer http.ResponseWriter, request *http.Request)  {
	values := request.URL.Query()
	uuid := values.Get("uuid")
	thread, err := models.ThreadByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		_, err := session(writer, request)
        if err != nil {
            generateHTML(writer, &thread, "layout", "navbar", "thread")
        } else {
            generateHTML(writer, &thread, "layout", "auth.navbar", "auth.thread")
        }
	}
}