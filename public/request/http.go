package request

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"myproject/public/fmt"
	"os/exec"
	"strings"
	"time"
	reatHttp "net/http"
)

type http struct {
}

var (
	Http http
	httpSetting httplib.BeegoHTTPSettings
)

func init() {
	httpSetting = httplib.BeegoHTTPSettings{
		UserAgent:        "beegoServer",
		ConnectTimeout:   15 * time.Second,
		ReadWriteTimeout: 15 * time.Second,
		Gzip:             true,
		DumpBody:         true,
	}
}

//重写beego原有http请求的Post方法
func (h *http) Post(url string) *httplib.BeegoHTTPRequest {
	req := httplib.Post(url)
	req.Setting(httpSetting)
	return req
}

//发送post json请求
func (h *http) PostJson(url string, body map[string]interface{}) (map[string]interface{}, error) {
	//logs.Info("send post json, ", url, body)
	req := h.Post(url)
	req.Header("content-type", "application/json")
	req.JSONBody(body)

	var data map[string]interface{}
	err := req.ToJSON(&data)
	return data, err
}

//发送post json请求
func (h *http) PostString(url string,s string) (map[string]interface{}, error) {
	//logs.Info("send post json, ", url, body)
	response ,err:= reatHttp.Post(url,"",strings.NewReader(s))

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	json.Unmarshal(body,data)
	fmt.Println(string(body))
	return data, err
}

func (h *http)Curl(arg ...string) (map[string]interface{},error){
	cmd := exec.Command("curl",arg...)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:can not obtain stdout pipe for command:%s\n", err)
		return nil,err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return nil,err
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return nil,err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return nil,err
	}

	var data map[string]interface{}
	json.Unmarshal(bytes,&data)
	return data,nil
}
