package controllers

import (
	
)

type ApiController struct {
	baseController
}

func (this *ApiController) CreateImage() {
	jsoninfo := this.GetString("json")
	if jsoninfo == "" {
		this.Ctx.WriteString("json is empty")
		return
	} else {
		this.Ctx.WriteString(jsoninfo)
		return
	}
}

func (this *ApiController) ExistsImage() {
	//TODO check image exists
}
	
func (this *ApiController) OnlineAll() {
	//TODO online deploy
}