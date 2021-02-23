package models

import (
	"time"
	"fmt"
)

type User struct {
	Id int
	Uuid string
	Name string
	Email string
	Password string
	CreatedAt time.Time
}

func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) value (?, ?, ?, ?)"
	stmtin, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmtin.Close()

	uuid := createUUID()
	stmtin.Exec(uuid, user.Email, user.Id, time.Now())

	stmtout, err := Db.Prepare("select id, uuid, email, user_id, created_at from sessions where uuid = ?")
	if err != nil {
		return
	}
	defer stmtout.Close()
	session = Session{}
	err = stmtout.QueryRow(uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("select id, uuid, email, user_id, created_at from sessions where user_id = ?", user.Id).
	Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (user *User) Create() (err error) {
	stmtin, err := Db.Prepare("insert into users (uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmtin.Close()
	uuid := createUUID()
	pwd := Encrypt(user.Password)
	_, err = stmtin.Exec(uuid, user.Name, user.Email, pwd, time.Now())
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	
	stmtout, err := Db.Prepare("select id, uuid, name, email, created_at from users where uuid = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmtout.Close()
	err = stmtout.QueryRow(uuid).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (user *User) Delete() (err error){
	stmt ,err := Db.Prepare("delete from users where id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return
}

func (user *User) Update() (err error) {
	stmtout, err := Db.Prepare("update users set name = ?, email = ? where id = ?")
	if err != nil {
		return
	}
	defer stmtout.Close()
	_, err = stmtout.Exec(user.Name, user.Email, user.Id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
    statement := "delete from users"
    _, err = Db.Exec(statement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
    return
}

func Users() (users []User, err error)  {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
        return
    }
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
    return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
    user = User{}
    err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
        Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
    return
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
    user = User{}
    err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = ?", uuid).
        Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
    return
}

// Create a new thread
func (user *User) CreateThread(topic string) (conv Thread, err error) {
    statement := "insert into threads (uuid, topic, user_id, created_at) values (?, ?, ?, ?)"
    stmtin, err := Db.Prepare(statement)
    if err != nil {
        return
    }
    defer stmtin.Close()

    uuid := createUUID()
    stmtin.Exec(uuid, topic, user.Id, time.Now())

    stmtout, err := Db.Prepare("select id, uuid, topic, user_id, created_at from threads where uuid = ?")
    if err != nil {
        return
    }
    defer stmtout.Close()

    // use QueryRow to return a row and scan the returned id into the Session struct
    err = stmtout.QueryRow(uuid).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
    return
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
    statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values (?, ?, ?, ?, ?)"
    stmtin, err := Db.Prepare(statement)
    if err != nil {
        return
    }
    defer stmtin.Close()

    uuid := createUUID()
    stmtin.Exec(uuid, body, user.Id, conv.Id, time.Now())

    stmtout, err := Db.Prepare("select id, uuid, body, user_id, thread_id, created_at from posts where uuid = ?")
    if err != nil {
        return
    }
    defer stmtout.Close()

    // use QueryRow to return a row and scan the returned id into the Session struct
    err = stmtout.QueryRow(uuid).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
    return
}