package controllers

import (
	"net/http"
	"bytes"
	"time"
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
	titannl = beego.AppConfig.String("titannodes")
	dockremoteport = beego.AppConfig.String("dockremoteport")
	agileCBUrl = beego.AppConfig.String("agilecallback")
)

type CreateImageInfo struct {
	AgileId string `json:"agileId"`
	Module string `json:"module"`
	SvnUrl string `json:"svnUrl"`
	ProdCmd string `json:"prodCmd"`
}

type ContainerInfo struct {
	Id string
	Image string
}

type AgileCallBackInfo struct {
	AgileId string `json:"agileId"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func pushImage(image string) error {
	cmd := exec.Command("docker", "push", image)
	return cmd.Run() 
}

func getContainerByIpAndImage(nodeIp, image string) ContainerInfo {
	// get all running containers
	var target ContainerInfo
	apiUrl := "http://"+nodeIp+":"+dockremoteport+"/containers/json"
	cmd := exec.Command("curl", apiUrl)
	var cmdout bytes.Buffer
	cmd.Stdout = &cmdout
	err := cmd.Run()
	if err != nil {
		beego.Error("[OnlineAll]Uable to get container info by remote api:", err)
		return target
	}
	var containerList []ContainerInfo
	if err = json.Unmarshal([]byte(cmdout.String()), &containerList); err != nil {
		beego.Error("[OnlineAll]Uable to get container info by remote api:", err)
    }
	// get target container
	for i:=0; i < len(containerList); i++ {
		if 0 == strings.Compare(containerList[i].Image, image) {
			target = containerList[i]
			break;
		}
	}
	return target
}

func agileCallBack(url string, agileCB AgileCallBackInfo) {
	if jsonString, err := json.Marshal(agileCB); err == nil {
		postUrl := agileCBUrl+"/"+url
		beego.Info("[AgileCallBack]Callback info:", url+"-"+string(jsonString))
		resp, postErr := http.Post(postUrl, "application/json/raw", strings.NewReader(string(jsonString)))
		defer resp.Body.Close()
		if postErr != nil {
			beego.Error("[AgileCallBack]Callback info:", postErr)
		}
		return
	}
	return
}

func createImage(info CreateImageInfo) {
	var agileCB AgileCallBackInfo
	agileCB.AgileId = info.AgileId
	agileCB.Status = "FALSE"
	// create workspace /home/titan/images/create/100000
	workspace := titanws+"/images/create/"+info.AgileId
	err := os.MkdirAll(workspace, 0755)
	if err != nil {
		beego.Error("[CreateImage]Uable to create dir:", workspace)
		agileCB.Message = "[CreateImage]Uable to create dir:"+workspace
		agileCallBack("createImage", agileCB)
		return
	}
	// cd workspace
	err = os.Chdir(workspace)
	if err != nil {
		beego.Error("[CreateImage]Uable to cd dir:", workspace)
		agileCB.Message = "[CreateImage]Uable to cd dir:"+workspace
		agileCallBack("createImage", agileCB)
		return
	}
	//TODO get Dockerfile from git/svn : later
	//Now cp Dockerfile from /home/titan/dockerfiles/module/Dockerfile
	dockerfilePath := titanws+"/dockerfiles/"+info.Module+"/Dockerfile"
	cmd := exec.Command("cp", dockerfilePath, ".")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to get dockerfile:", dockerfilePath)
		agileCB.Message = "[CreateImage]Uable to get dockerfile:"+dockerfilePath
		agileCallBack("createImage", agileCB)
		return
	}
	//modify dockerfile to set prodCmd
	prodCmd := "RUN "+info.ProdCmd+" -P ./tmp"
	prodCmd = strings.Replace(prodCmd, "/", "\\/",-1)
	prodCmd = "s/#heresetprodCmd/"+prodCmd+"/g"
	cmd = exec.Command("sed", "-i", "-e", prodCmd, "./Dockerfile")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to modify dockerfile:", prodCmd)
		agileCB.Message = "[CreateImage]Uable to modify dockerfile:"+prodCmd
		agileCallBack("createImage", agileCB)
		return
	}
	//create image by dockerfile
	image := dockerregistory+info.Module+":"+info.AgileId
	cmd = exec.Command("docker", "build", "-t", image, ".")
	err = cmd.Run()
	if err != nil {
		beego.Error("[CreateImage]Uable to create image by dockerfile:", err)
		agileCB.Message = "[CreateImage]Uable to create image by dockerfile"
		agileCallBack("createImage", agileCB)
		return
	}
	//push image
	err = pushImage(image)
	if err != nil {
		beego.Error("[CreateImage]Uable to push image:", err)
		agileCB.Message = "[CreateImage]Uable to push image"
		agileCallBack("createImage", agileCB)
		return
	}
	beego.Info("[CreateImage]SUCCESS:", image)
	agileCB.Message = "SUCCESS"
	agileCB.Status = "TRUE"
	agileCallBack("createImage", agileCB)
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
			go createImage(createImageInfo)
		}
		this.Ctx.WriteString("SUCCESS")
		return
	}
}

func (this *ApiController) ExistsImage() {
	this.Ctx.WriteString("SUCCESS")
	return
}

func singleDeploy(nodeIp, image string) {
	containerInfo := getContainerByIpAndImage(nodeIp, image)
	if "" == containerInfo.Id {
		return
	}
	// begin deploy
	// stop old container
	cmdString := "curl -v --raw -X POST http://"+nodeIp+":"+dockremoteport+"/containers/"+containerInfo.Id+"/stop?t=5"
	in := bytes.NewBuffer(nil)
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = in
	in.WriteString(cmdString)
	if err := cmd.Run(); err != nil {
		beego.Error("[OnlineAll]Uable to stop container:", err)
	}
	return
}
	
func (this *ApiController) OnlineAll() {
	agileId := this.GetString("agileId")
	module := this.GetString("module")
	beego.Info("[OnlineAll]Get online info:", agileId+"@"+module)
	// tag latest
	image := dockerregistory+module
	cmd := exec.Command("docker", "tag", "-f", image+":"+agileId, image+":latest")
	err := cmd.Run()
	if err != nil {
		beego.Error("[OnlineAll]Uable to tag latest:", image+":"+agileId)
	}
	// push image
	err = pushImage(image+":latest")
	if err != nil {
		beego.Error("[OnlineAll]Uable to push latest image:", image)
	}
	// online deploy
	nodes := strings.Split(titannl, ",")
	for i := 0; i < len(nodes); i++ {
		go singleDeploy(nodes[i], image+":latest")
		time.Sleep(20 * time.Second)
	}
	var agileCB AgileCallBackInfo
	agileCB.AgileId = agileId
	agileCB.Status = "TRUE"
	agileCB.Message = "SUCCESS"
	agileCallBack("onlineAll", agileCB)
	this.Ctx.WriteString("SUCCESS")
	return
}