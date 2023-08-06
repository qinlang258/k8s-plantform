package main

import (
	"github.com/gin-gonic/gin"
	"k8s-plantform/config"
	"k8s-plantform/controller"
	"k8s-plantform/db"
	"k8s-plantform/middle"
	"k8s-plantform/service"
)

func main() {
	// 初始化gin
	db.Init()
	defer db.Close()
	service.K8s.Init()
	r := gin.Default()
	r.Use(middle.JWTAuth())
	r.Use(middle.Cors())
	controller.Router.InitApiRouter(r)
	// gin 程序启动
	//r.Run(config.ListenAdd)
	r.Run(config.ListenAddr)

}
