package main

import "github.com/kataras/iris"

type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

func main() {
	app := iris.New()

	app.RegisterView(iris.HTML("./web/views", ".html").Reload(true))

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})

	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})

	app.Done(func(ctx iris.Context) {})

	app.Post("/decode", func(ctx iris.Context) {
		var user User
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
	})

	app.Get("/encode", func(ctx iris.Context) {
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}

		ctx.JSON(doe)
	})

	app.Get("/profile/{username:string}/{pwd:string}", profileByUsername)

	userRoutes := app.Party("/users",logThisMiddleware)
	{
		userRoutes.Get("/{id:int min(1)}",getUserByID)
		//userRoutes.Post("/create",createUser)
	}

	app.Run(iris.Addr(":8080"))
}

func logThisMiddleware(ctx iris.Context)  {
	ctx.Application().Logger().Infof("Path: %s |IP: %s",ctx.Path(),ctx.RemoteAddr())
	ctx.Next()
}

func profileByUsername(ctx iris.Context) {

	username := ctx.Params().Get("username")
	pwd := ctx.Params().Get("pwd")

	ctx.ViewData("Username", username)
	ctx.ViewData("Password", pwd)
	ctx.View("profile.html")
}

func getUserByID(ctx iris.Context)  {
	userID := ctx.Params().Get("id")
	user := User{Username:"username"+userID}
	ctx.XML(user)
}

func createUser(ctx iris.Context)  {

}


