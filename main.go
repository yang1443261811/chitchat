package main

import (
	"log"
	"net/http"
	. "chitchat/routes"
	. "chitchat/config"
)

func main()  {
	startWebServer("8080")
}

func startWebServer(port string)  {
	// 在入口位置初始化全局配置
	config := LoadConfig()
	r := NewRouter()
	// 处理静态资源文件
	assets := http.FileServer(http.Dir("public"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", assets))

	http.Handle("/", r)

	log.Println("Starting HTTP service at " + port)
    err := http.ListenAndServe(config.App.Address, nil) // 启动协程监听请求

    if err != nil {
        log.Println("An error occured starting HTTP listener at port " + config.App.Address)
        log.Println("Error: " + err.Error())
    }
}