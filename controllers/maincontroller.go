package controllers

import (
	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

func (this *baseController) Prepare() {

}

type MainController struct {
	baseController
}

func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "index.tpl"
}
