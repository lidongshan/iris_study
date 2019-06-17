package controllers

import (
	"awesomeProject/datamodels"
	"awesomeProject/services"
	"errors"
	"github.com/kataras/iris"
)

// MovieController is our /movies controller.
type MovieController struct {
	// Our MovieService, it's an interface which
	// is binded from the main application.
	Service services.MovieService
}

// Get 返回 movies 的列表
// 演示：
// curl -i http://localhost:8080/movies
//
// The correct way if you have sensitive data:
// func (c *MovieController) Get() (results []viewmodels.Movie) {
//     data := c.Service.GetAll()
//
//     for _, movie := range data {
//         results = append(results, viewmodels.Movie{movie})
//     }
//     return
// }
// 否则，只返回数据模型。
func (c *MovieController) Get() (results []datamodels.Movie) {
	return c.Service.GetAll()
}

// GetBy 返回一个 movie
// 演示:
// curl -i http://localhost:8080/movies/1
func (c *MovieController) GetBy(id int64) (movie datamodels.Movie, found bool) {
	return c.Service.GetByID(id) // it will throw 404 if not found.
}

// PutBy 更新一个movie
// 演示:
// curl -i -X PUT -F "genre=Thriller" -F "poster=@/Users/kataras/Downloads/out.gif" http://localhost:8080/movies/1
func (c *MovieController) PutBy(ctx iris.Context, id int64) (datamodels.Movie, error) {
	// 获取请求内的 poster 和 genre 字段数据
	file, info, err := ctx.FormFile("poster")
	if err != nil {
		return datamodels.Movie{}, errors.New("failed due form file 'poster' missing")
	}
	// 关闭
	file.Close()

	// 假设这是一个上传文件的 url ...
	poster := info.Filename
	genre := ctx.FormValue("genre")

	return c.Service.UpdatePosterAndGenreByID(id, poster, genre)
}

// DeleteBy 删除一个 movie
// 演示:
// curl -i -X DELETE -u admin:password http://localhost:8080/movies/1
func (c *MovieController) DeleteBy(id int64) interface{} {
	wasDel := c.Service.DeleteByID(id)
	if wasDel {
		// 返回被删除的 movie 的 id
		return iris.Map{"deleted": id}
	}
	//在这里，我们可以看到一个方法函数可以返回两种类型中的任何一种（map 或者 int）,
	// 我们不用指定特定的返回类型。
	return iris.StatusBadRequest
}
