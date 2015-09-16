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
		    beego.NSRouter("/createImage", &controllers.CreateImage{}),
			// /api/onlineAll
			beego.NSRouter("/onlineAll", &controllers.OnlineAll{}),
		)
		
	beego.AddNamespace(ns)
}