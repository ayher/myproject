package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

var GinRouter *gin.Engine
func init()  {
	//gin 路由
	ginRouter := gin.Default()
	v1Group := ginRouter.Group("/v1")
	v1Group.Handle("POST", "/newAddress", newAddress)
	v1Group.Handle("POST", "/getBrowser", getBrowser)
	GinRouter=ginRouter
}

func GetPostBody(context *gin.Context) map[string]interface{}{
	buf := make([]byte, 1024)
	context.Request.Body.Read(buf)

	// 去掉空字节
	index := bytes.IndexByte(buf, 0)
	b:= buf[:index]

	var m map[string]interface{}
	err:=json.Unmarshal(b,&m)

	if err!=nil{
		fmt.Println(err)
		return nil
	}

	return m

}

func getClient(){

}