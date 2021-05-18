package main

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"myproject/public/consul/routers"
)

// https://blog.csdn.net/qq_36453564/article/details/107781480
func main(){
	//consul注册中心
	consulReg := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	//创建web服务器
	server := web.NewService(
		web.Name("prodservice"),	// 服务名
		//web.Address(":8001"),
		web.Handler(routers.GinRouter),	// 路由
		web.Metadata(map[string]string{"protocol": "http"}), // 添加这行代码
		web.Registry(consulReg),	// 注册服务
	)
	//初始化服务器
	server.Init()
	//运行
	err := server.Run()
	if err != nil{
		log.Panic(err)
	}
}