package routers

import (
	"context"
	"github.com/gin-gonic/gin"
	"myproject/public/fmt"
	"myproject/public/micro/proto"
	"github.com/micro/go-micro"
)

func getBrowser(c *gin.Context) {
	body:=GetPostBody(c)

	fmt.Println(body)
	var data interface{}
	if body!=nil{
		if coin,ok:=body["coin"].(string);ok{

			clientName:="go.micro."+coin
			service := micro.NewService(
				micro.Name(clientName),
			)
			// Create new greeter client
			greeter := proto.NewHelloService(clientName, service.Client())

			rsp1, err := greeter.GetBrowser(context.TODO(), &proto.GetBrowserRequest{T:"test"})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(rsp1)
		}


	}
	//返回JSON
	c.JSON(200,
		gin.H{
			"data": data, //自定义函数，生成prod列表
		})
}