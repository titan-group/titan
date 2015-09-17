package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	ns :=
	    //API
	    beego.NewNamespace("/api",
		    // /api/createImage
		    beego.NSRouter("/createImage", &controllers.ApiController{}, "post:CreateImage"),
			// /api/existsImage
			beego.NSRouter("/existsImage", &controllers.ApiController{}, "post:ExistsImage"),
			// /api/onlineAll
			beego.NSRouter("/onlineAll", &controllers.ApiController{}, "get:OnlineAll"),
		)
		
	beego.AddNamespace(ns)
}