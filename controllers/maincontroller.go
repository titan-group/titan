package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type CreateImage struct {
	beego.Controller
}

type OnlineAll struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplNames = "index.tpl"
}

func (this *CreateImage) Get() {
	jsoninfo := this.GetString("json")
	if jsoninfo == "" {
		this.Ctx.WriteString("json is empty")
		return
	} else {
		this.Ctx.WriteString(jsoninfo)
		return
	}
	
func (this *OnlineAll) Get() {
	//TODO online deploy
}
	
}