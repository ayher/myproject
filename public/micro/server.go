package micro

import (
	"context"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"myproject/public/fmt"
	proto "myproject/public/micro/proto"
	"os"
)

type Micro struct {
	Coin string
}

func (self *Micro)runClient() {
	clientName:="go.micro."+self.Coin
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

func (self *Micro)Register(Api interface{})  {

	service := micro.NewService(
		micro.Name("go.micro."+self.Coin),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloword",
		}),
		micro.Flags(cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			if c.Bool("run_client") {
				self.runClient()
				os.Exit(0)
			}
		}),
	)
	_ = proto.RegisterHelloHandler(service.Server(), Api.(proto.HelloHandler))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}