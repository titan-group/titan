package controllers

import (
	"strings"
	"os/exec"
	"os"
	"encoding/json"
	"github.com/astaxie/beego"
)

type ApiController struct {
	baseController
}

var (
    titanws = beego.AppConfig.String("titanws")
	dockerregistory =beego.AppConfig.String("dockerregistory")
)

type CreateImageInfo struct {
	AgileId string `json:"agileId"`
	Module string `json:"module"`
	SvnUrl string `json:"svnUrl"`
	ProdCmd string `json:"prodCmd"`
}

func createImage(info CreateImageInfo) {
	// create workspace /home/titan/images/create/100000
	workspace := titanws+"/images/create/"+info.AgileId
	err := os.MkdirAll(workspace, 0755)
	if err != nil {
		beego.Error("[CreateImage]Uable to create dir:", workspace)
		return
	}
	// cd workspace
	err = os.Chdir(workspace)
	if err != nil {
		return
	}
	//TODO get Dockerfile from git/svn : later
	//Now cp Dockerfile from /home/titan/dockerfiles/module/Dockerfile
	dockerfilePath := titanws+"/dockerfiles/"+info.Module+"/Dockerfile"
	cmd := exec.Command("cp", dockerfilePath, ".")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to get dockerfile:", dockerfilePath)
		return
	}
	//modify dockerfile to set prodCmd
	prodCmd := "RUN "+info.ProdCmd+" ./tmp"
	prodCmd = strings.Replace(prodCmd, "/", "\\/",-1)
	prodCmd = "s/#heresetprodCmd/"+prodCmd+"/g"
	cmd = exec.Command("sed", "-i", "-e", prodCmd, "./Dockerfile")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to modify dockerfile:", prodCmd)
		return
	}
	//create image by dockerfile
	image := dockerregistory+info.Module+":"+info.AgileId
	cmd = exec.Command("docker", "build", "-t", image, ".")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to create image by dockerfile:", err)
		return
	}
	//push image
	cmd = exec.Command("docker", "push", image)
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to push image:", err)
		return
	}
	beego.Info("[CreateImage]SUCCESS:", image)
	return
}

func (this *ApiController) CreateImage() {
	jsonInfo := this.GetString("json")
    beego.Info("[CreateImage]Get create info:", jsonInfo)
	var createImageInfo CreateImageInfo
	if jsonInfo == "" {
		this.Ctx.WriteString("ERROR")
		return
	} else {
		if err := json.Unmarshal([]byte(jsonInfo), &createImageInfo); err == nil {
			//TODO createe
			go createImage(createImageInfo)
		}
		this.Ctx.WriteString("SUCCESS")
		return
	}
}

func (this *ApiController) ExistsImage() {
	//TODO
	return
}
	
func (this *ApiController) OnlineAll() {
	//TODO online deploy
	return
}