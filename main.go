package main

import (
	"github.com/kataras/iris"
)

type AppConfig struct {
	AppName    string `json:"app_name"`
	Port       int
	StaticPath string `json:"static_path"`
	Mode       string `json:"mode"`
}


func main() {
	//1.创建app结构体对象
	app := iris.New()

	app.StaticWeb("/manage/static", "./static")

	app.RegisterView(iris.HTML("./static", ".html"))

	//usersRouter := app.Party("admin",UserMiddleware)
	//2.端口监听
	app.Run(iris.Addr(":8000"), iris.WithoutServerError(iris.ErrServerClosed))


}
