package routers

import (
	"beego_movie/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/beego_movie", &controllers.CwMovieController{}, "*:CwMovie")}
