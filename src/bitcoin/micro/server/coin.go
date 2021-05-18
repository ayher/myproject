package server

import (
	"context"
	"myproject/public/micro/proto"
	"myproject/src/bitcoin/public"
)

type Btc struct {

}

func (self *Btc) NewAddress(ctx context.Context, req *proto.NewAddressRequest, rsp *proto.NewAddressResponse) error {
	pr,pu:= public.NewKey()
	ad:= public.PubkeyToAddress(pu)

	rsp.Privateket=pr
	rsp.PublicKey=pu
	rsp.Address=ad
	return nil
}

func (self *Btc)GetBrowser(ctx context.Context, req *proto.GetBrowserRequest, rsp *proto.GetBrowserResponse) error  {
	var url string
	if req.T=="test"{
		url="https://chain.so/btc"
	}else{
		url="https://btc.com"
	}
	rsp.BrowserUrl=url
	return nil
}