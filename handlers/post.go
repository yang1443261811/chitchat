package handlers

import (
    "fmt"
    "chitchat/models"
    "net/http"
)

func PostThread(writer http.ResponseWriter, request *http.Request)  {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot find user")
		}
		valid := request.ParseForm()
		if valid != nil {
			danger(valid, "Missing required parameter")
		}
		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		thread, err := models.ThreadByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read thread")
		}
		_, ok := user.CreatePost(thread, body)
		if ok != nil {
			danger(ok, "Cannot create post")
		}
		url := fmt.Sprint("/thread/read?uuid=", uuid)
        http.Redirect(writer, request, url, 302)
	}
}