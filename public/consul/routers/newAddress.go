package routers

import (
	"github.com/gin-gonic/gin"
	"myproject/public/consul/Helper"
	"myproject/public/consul/ProdService"
)

func newAddress(context *gin.Context) {
	//请求对象
	var pr Helper.ProdsRequest
	//获取请求参数
	err := context.Bind(&pr)
	//默认为2个
	if err != nil || pr.Size <= 0 {
		pr = Helper.ProdsRequest{Size: 3}
	}
	//返回JSON
	context.JSON(200,
		gin.H{
			"data": ProdService.NewProdList(pr.Size), //自定义函数，生成prod列表
		})
}