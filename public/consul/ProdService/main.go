package ProdService

import "strconv"

//产品model
type ProdModel struct {
	ProdID int	//产品id
	ProdName string		//产品名称
}

//实例化对象
func NewProd(id int,pname string) *ProdModel {
	return &ProdModel{
		ProdID: id,
		ProdName: pname,
	}
}

//生成产品列表
func NewProdList(n int) []*ProdModel {
	ret := make([]*ProdModel, 0)
	for i:=0;i<n;i++ {
		ret = append(ret, NewProd(10+i, "prodname"+strconv.Itoa(100+i)))
	}
	return ret
}